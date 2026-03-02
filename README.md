# UOS Kernel — 电力系统资源操作系统内核

## 项目概述

UOS Kernel 将电力系统中的设备（断路器、变压器、线路等）抽象为 **Actor**，在 Actor 之上构建 POSIX 风格的资源访问内核，再通过 **App 框架** 和 **gRPC 服务** 支撑上层业务流程的编排。

```
Python / Java 业务流程
    ↓  gRPC
┌─────────────────────────────────┐
│  gRPC Server (server/)          │  ← 远程接入层
└────────────┬────────────────────┘
             ↓
┌─────────────────────────────────┐
│  App 框架 (app/)                │  ← 业务编排层：Process + helpers
└────────────┬────────────────────┘
             ↓
┌─────────────────────────────────┐
│  Kernel (kernel/)               │  ← 系统调用层：Open/Close/Read/Write/Ioctl/Watch
└────────────┬────────────────────┘
             ↓
┌─────────────────────────────────┐
│  Meta (meta/)                   │  ← 类型系统：TypeRegistry / TypeDescriptor
└────────────┬────────────────────┘
             ↓
┌─────────────────────────────────┐
│  Actor System (actor/)          │  ← 设备层：Actor + Capacity + Event
└─────────────────────────────────┘
```

核心理念：**业务流程不直接接触 Actor，只通过 Kernel API 操控资源**。

## 项目结构

```
uos-kernel/
├── actor/                         # 设备层：Actor 系统
│   ├── actor.go                   #   Actor 接口 + BaseActor
│   ├── system.go                  #   System：注册、事件分发、Watch
│   ├── resource_actor.go          #   ResourceActor 接口
│   ├── base_resource_actor.go     #   BaseResourceActor 实现
│   ├── capacity.go                #   Capacity 接口
│   ├── binding.go                 #   Binding 机制
│   ├── message.go                 #   Message 接口
│   ├── event.go                   #   Event 定义
│   └── event_descriptor.go        #   EventDescriptor 元数据
│
├── meta/                          # 类型系统
│   ├── types.go                   #   TypeDescriptor / AttributeDescriptor
│   ├── registry.go                #   TypeRegistry
│   └── loader.go                  #   YAML 加载器
│
├── kernel/                        # 系统调用层
│   ├── kernel.go                  #   Kernel：Open / Close / Read / Write / Stat / Ioctl / Watch
│   ├── manager.go                 #   Manager：资源描述符管理
│   ├── read_write.go              #   Read / Write 实现
│   └── rctl.go                    #   RCtl：控制命令派发
│
├── app/                           # 业务编排层
│   ├── app.go                     #   App + Process 接口
│   └── helpers.go                 #   WaitForProperty / WaitForEvent
│
├── server/                        # gRPC 远程接入层
│   ├── server.go                  #   KernelServer：7 个 RPC 方法
│   ├── registry.go                #   MessageRegistry：命令工厂 + 事件序列化
│   ├── codec.go                   #   protobuf Struct ↔ map 转换
│   └── gen/                       #   生成的 Go protobuf 代码
│
├── proto/uos/kernel/v1/
│   └── kernel.proto               # gRPC 服务定义
│
├── sdk/python/                    # Python SDK
│   ├── uos_kernel/
│   │   ├── client.py              #   KernelClient（gRPC 客户端）
│   │   ├── process.py             #   Process 基类
│   │   └── helpers.py             #   wait_for_property / execute_capacity
│   ├── examples/
│   │   └── substation_maintenance.py
│   └── gen/                       #   生成的 Python protobuf 代码
│
├── examples/
│   ├── substation_maintenance/        # Go 示例：本地 App 框架
│   └── substation_maintenance_grpc/   # Go 示例：gRPC 服务端
│
├── Makefile                       # proto 代码生成
├── cmd/uos-kernel/                # 主程序入口（含 typesystem.yaml）
├── docs/                          # 设计文档
└── infra/                         # 基础设施配置
```

## 核心概念

### Actor — 设备资源

每个电力设备是一个 Actor，拥有：

- **属性**（Properties）：电压、电流、温度、开关状态等
- **能力**（Capacity）：定义设备可执行的操作，通过 Message 路由
- **事件**（Event）：设备异常、状态变更等，通过 EventEmitter 发射

```go
// 自定义 Actor（以断路器为例）
type BreakerActor struct {
    *actor.BaseResourceActor
    // 设备状态字段...
}

// 注册能力
breaker.AddCapacity(&BreakerSwitchingCapacity{...})

// 发射事件
breaker.GetEventEmitter().EmitEvent(&DeviceAbnormalEvent{...})
```

### Kernel — POSIX 风格 API

```go
fd, _ := k.Open("Breaker", "BREAKER-001", 0)   // 打开资源 → 文件描述符
defer k.Close(fd)

state, _ := k.Read(ctx, fd)                      // 读取属性
stat, _  := k.Stat(ctx, fd)                      // 查询能力和事件定义

k.Ioctl(ctx, fd, 0x1003, map[string]interface{}{  // 执行能力
    "message": &OpenBreakerCommand{Reason: "检修"},
})

ch, cancel, _ := k.Watch(fd, nil)                // 监听事件（nil = 所有类型）
defer cancel()
for event := range ch { ... }
```

**标准 Ioctl 命令**：

| 命令                    | 值     | 说明               |
| ----------------------- | ------ | ------------------ |
| CMD_GET_RESOURCE_INFO   | 0x1000 | 获取资源信息       |
| CMD_LIST_CAPABILITIES   | 0x1001 | 列出能力           |
| CMD_LIST_EVENTS         | 0x1002 | 列出事件           |
| CMD_EXECUTE_CAPACITY    | 0x1003 | 执行能力（最常用） |

**事件类型**：

| 类型              | 说明                       |
| ----------------- | -------------------------- |
| state_change      | 生命周期变更               |
| attribute_change  | 属性变更（name + value）   |
| capability_change | 能力变更                   |
| custom            | 业务自定义事件             |

### App — 业务流程编排

```go
type Process interface {
    Name() string
    Run(ctx context.Context, k *kernel.Kernel) error
}

// 启动应用
application := app.New(k)
application.Run(ctx, &MaintenanceProcess{...}, &MonitorProcess{...})
```

辅助函数：

```go
// 等待属性达到预期值（先 Read 检查，再 Watch 等待）
app.WaitForProperty(ctx, k, fd, "isOpen", true)

// 等待匹配的自定义事件
app.WaitForEvent(ctx, k, fd, func(e kernel.Event) bool { ... })
```

### gRPC — 跨语言接入

gRPC 服务将 Kernel API 1:1 映射为远程 RPC，动态数据使用 `google.protobuf.Struct`，通过 `_type` 字段区分命令/事件的具体类型。

```protobuf
service KernelService {
  rpc Open(OpenRequest)   returns (OpenResponse);
  rpc Close(CloseRequest) returns (CloseResponse);
  rpc Read(ReadRequest)   returns (ReadResponse);
  rpc Write(WriteRequest) returns (WriteResponse);
  rpc Stat(StatRequest)   returns (StatResponse);
  rpc Ioctl(IoctlRequest) returns (IoctlResponse);
  rpc Watch(WatchRequest) returns (stream WatchEvent);  // server streaming
}
```

服务端需要注册 **MessageRegistry**，将命令类型名映射到 Go 工厂函数，将事件类型映射到序列化器。

### Python SDK

```python
from uos_kernel import KernelClient
from uos_kernel.process import Process
from uos_kernel.helpers import wait_for_property, execute_capacity

class MaintenanceProcess(Process):
    def name(self) -> str:
        return "变电站停电检修"

    def run(self, client: KernelClient) -> None:
        fd = client.open("Breaker", "BREAKER-001", 0)
        try:
            # 监听事件
            for event in client.watch(fd):
                if event["data"].get("_type") == "DeviceAbnormalEvent":
                    # 执行能力
                    execute_capacity(client, fd, "OpenBreakerCommand", {
                        "reason": "停电检修",
                    })
                    # 等待确认
                    wait_for_property(client, fd, "isOpen", True)
                    break
        finally:
            client.close(fd)

with KernelClient("localhost:50051") as client:
    MaintenanceProcess().run(client)
```

## 快速开始

### Go 本地示例

```bash
cd examples/substation_maintenance
go run .
```

Actor 创建设备 → Kernel 提供 API → App 运行检修流程 → Watch 接收事件 → Ioctl 执行命令。

### gRPC + Python 跨语言示例

```bash
# 1. 生成 protobuf 代码
make proto-all

# 2. 启动 Go gRPC 服务端
cd examples/substation_maintenance_grpc && go run .

# 3. 运行 Python 客户端（另一个终端）
cd sdk/python
python -m venv .venv && source .venv/bin/activate
pip install -e .
python examples/substation_maintenance.py
```

## 设计要点

- **分层隔离**：业务流程只通过 Kernel API 操控设备，不直接接触 Actor
- **命令-确认模式**：Ioctl 发送命令 → Watch/WaitForProperty 等待属性变更确认
- **事件驱动**：Actor 发射事件 → Watch 推送给所有观察者 → 业务流程响应
- **跨语言**：gRPC + protobuf Struct 解决 Go `interface{}` 的跨语言序列化问题
- **Go 多模块**：actor、meta、kernel、app、server 各自独立模块，依赖关系清晰

## 技术栈

- **Go** — 内核、Actor 系统、gRPC 服务端
- **Python** — SDK 客户端
- **gRPC / Protocol Buffers** — 跨语言通信
- **Actor 模型** — 设备资源并发抽象

## 相关文档

- [Actor 系统设计讨论](docs/actors_system_discussion_20251209.md)
- [POSIX 风格资源服务层设计](docs/posix_resource_service_layer_20251210.md)
- [开发总结](docs/development_summary_20251216.md)
