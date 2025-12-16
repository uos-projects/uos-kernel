# UOS Kernel - 电力系统资源操作系统内核

## 项目概述

UOS Kernel 是一个基于 CIM (Common Information Model) 标准的电力系统资源操作系统内核。它将电力系统中的资源抽象为 Actor，实现了概念上的驱动层，并在此基础上提供了 POSIX 风格的资源访问接口。

项目参考操作系统内核的对象类型系统（Object Type System）设计，实现了类型系统内核，提供类型验证、操作路由和统一的资源管理接口。

## 核心设计思想

### 1. 分层架构

```
┌─────────────────────────────────────┐
│   用户应用层                          │
│   (使用 ResourceKernel API)          │
└─────────────────────────────────────┘
           ↓
┌─────────────────────────────────────┐
│   Resource Kernel (resource包)      │  ← 面向用户的高级接口
│   - 类型验证                         │
│   - ioctl命令映射                    │
│   - POSIX风格系统调用                │
└─────────────────────────────────────┘
           ↓
┌─────────────────────────────────────┐
│   Type System (kernel包)             │  ← 类型系统定义层
│   - TypeRegistry                    │
│   - TypeDescriptor                  │
│   - CIM Converter                   │
└─────────────────────────────────────┘
           ↓
┌─────────────────────────────────────┐
│   Resource Manager (resource包)     │  ← 资源管理层
│   - 资源描述符管理                   │
│   - 引用计数                         │
│   - 排他性控制                       │
└─────────────────────────────────────┘
           ↓
┌─────────────────────────────────────┐
│   Actor System (actors包)            │  ← 实现层
│   - PowerSystemResourceActor        │
│   - Capacity Implementations        │
│   - Message Routing                 │
└─────────────────────────────────────┘
```

### 2. CIM 模型到 Actor 系统的映射

基于 CIM 标准，将电力系统中的资源映射为 Actor 系统：

- **PowerSystemResource** → `PowerSystemResourceActor`
  - 每个电力系统资源（如发电机、变压器、断路器、线路等）对应一个 Actor
  - Actor 代表资源实体，封装资源的状态和行为

- **Control** → `Capacity` (能力接口)
  - 每个控制类型（如 SetPoint、Command、AccumulatorReset 等）对应一种 Capacity
  - Capacity 定义了资源可以执行的控制操作

- **Measurement** → `MeasurementCapacity`
  - 每个测量类型（如 AnalogMeasurement、DiscreteMeasurement 等）对应一种 MeasurementCapacity
  - MeasurementCapacity 处理测量数据的订阅和更新

### 3. 类型系统内核

参考操作系统内核的对象类型系统，实现了类型系统内核：

- **类型定义**：通过 YAML DSL 定义资源类型、属性、关系、能力
- **类型注册**：TypeRegistry 管理所有类型描述符
- **类型验证**：在执行操作前验证资源类型和能力
- **CIM转换**：自动将 CIM 模型转换为类型系统定义

### 4. POSIX 风格接口

提供类似操作系统内核的系统调用接口：

- `Open(resourceType, resourceID, flags)` - 打开资源（带类型验证）
- `Close(fd)` - 关闭资源
- `Read(ctx, fd)` - 读取资源状态
- `Write(ctx, fd, req)` - 写入资源状态
- `Stat(ctx, fd)` - 查询资源信息（包含类型描述符）
- `Ioctl(ctx, fd, request, argp)` - 控制操作（ioctl命令映射和类型验证）

## 项目结构

```
uos-kernel/
├── actors/                    # Actor 系统实现
│   ├── actor.go              # 基础 Actor
│   ├── system.go             # Actor 系统管理器
│   ├── resource_actor.go     # PowerSystemResourceActor
│   ├── capacity_factory.go   # Capacity 工厂
│   ├── capacities/           # Capacity 实现
│   │   ├── base.go           # 基础 Capacity
│   │   ├── measurement_base.go  # Measurement 基类
│   │   ├── analog_measurement.go
│   │   ├── discrete_measurement.go
│   │   ├── accumulator_measurement.go
│   │   ├── set_point.go
│   │   ├── command.go
│   │   └── accumulator_reset.go
│   ├── mq/                   # MQ Consumer 接口
│   │   ├── consumer.go
│   │   └── mock_consumer.go
│   └── cmd/                  # 示例程序
│
├── kernel/                    # 类型系统定义层
│   ├── typesystem.yaml       # 类型系统DSL定义
│   ├── types.go              # 类型描述符数据结构
│   ├── registry.go           # 类型注册表
│   ├── loader.go             # YAML加载器
│   ├── cim_converter.go      # CIM到类型系统转换器
│   └── README.md             # 类型系统使用文档
│
├── resource/                  # 资源管理层
│   ├── manager.go            # 资源描述符管理（Open/Close）
│   ├── read_write.go         # 资源状态访问（Read/Write）
│   ├── rctl.go               # 资源控制（RCtl）
│   ├── kernel.go             # ResourceKernel（高级接口）
│   ├── cmd/
│   │   └── kernel_example/   # 使用示例
│   └── example_test.go       # 测试示例
│
├── cim/                       # CIM 数据集和模式
│   ├── datasets/             # CIM 数据集
│   └── schemata/             # CIM 模式定义
│
├── docs/                      # 设计文档
│   ├── actors_system_discussion_20251209.md
│   ├── posix_resource_service_layer_20251210.md
│   └── development_summary_20251216.md
│
└── infra/                     # 基础设施配置
    ├── docker-compose.yml
    └── USAGE.txt
```

## 核心组件

### 1. Actor 系统 (`actors/`)

**BaseActor**
- 基础 Actor 实现，包含邮箱（mailbox）和消息处理循环
- 支持异步消息传递和生命周期管理

**PowerSystemResourceActor**
- 代表 CIM 中的 PowerSystemResource
- 动态管理 Capabilities（Control 和 Measurement）
- 消息自动路由到对应的 Capacity

**Capacity 接口**
```go
type Capacity interface {
    Name() string
    CanHandle(msg Message) bool
    Execute(ctx context.Context, msg Message) error
    ResourceID() string
}
```

### 2. 类型系统 (`kernel/`)

**TypeRegistry**
- 管理所有类型描述符
- 支持类型注册、查询、验证

**TypeDescriptor**
- 定义资源类型的元数据
- 包含属性、关系、能力、生命周期等信息
- 支持类型继承和查询

**CIM Converter**
- 将 CIM 实体转换为类型描述符
- 从 CIM 关系推断能力
- 自动推断关系的基数和反向关系

### 3. 资源管理 (`resource/`)

**ResourceManager**
- 管理资源描述符的分配和回收
- 维护资源引用计数和排他性控制
- 提供底层资源管理接口

**ResourceKernel**
- 面向用户的高级接口
- 整合类型系统和资源管理
- 提供 POSIX 风格的系统调用接口
- 支持类型验证和操作路由

## 快速开始

### 1. 基本使用

```go
package main

import (
    "context"
    "github.com/uos-projects/uos-kernel/actors"
    "github.com/uos-projects/uos-kernel/resource"
)

func main() {
    ctx := context.Background()

    // 1. 创建 Actor 系统
    system := actors.NewSystem(ctx)
    defer system.Shutdown()

    // 2. 创建资源 Actor
    actor := actors.NewPowerSystemResourceActor("BREAKER_001", "Breaker", nil)
    system.Register(actor)

    // 3. 创建资源内核
    k := resource.NewResourceKernel(system)

    // 4. 加载类型系统定义
    k.LoadTypeSystem("kernel/typesystem.yaml")

    // 5. 打开资源（带类型验证）
    fd, err := k.Open("Breaker", "BREAKER_001", 0)
    defer k.Close(fd)

    // 6. 查询资源信息
    stat, err := k.Stat(ctx, fd)

    // 7. 读取资源状态
    state, err := k.Read(ctx, fd)

    // 8. 执行控制操作（ioctl命令映射）
    result, err := k.Ioctl(ctx, fd, 0x1004, map[string]interface{}{"value": 150.0})
}
```

### 2. 使用 ResourceManager（底层接口）

```go
// 创建资源管理器
rm := resource.NewResourceManager(system)

// 打开资源
fd, err := rm.Open("BE-G4")
defer rm.Close(fd)

// 读取资源状态
state, err := rm.Read(ctx, fd)

// 使用 RCtl 发送控制命令
setPointArg := map[string]interface{}{"value": 150.0}
rm.RCtl(ctx, fd, resource.CMD_SET_POINT, setPointArg)
```

## 类型系统定义

在 `kernel/typesystem.yaml` 中定义资源类型：

```yaml
resource_types:
  - name: Breaker
    base_type: PowerSystemResource
    attributes:
      - name: ratedCurrent
        type: float
        required: false
    capabilities:
      - name: SwitchControl
        operations: [open, close, queryState]
      - name: Control
        operations: [execute, query, cancel]
```

## 设计优势

1. **分层清晰**：kernel包是纯类型系统定义，resource包整合类型系统和资源管理
2. **类型安全**：通过类型系统验证，减少运行时错误
3. **符合领域模型**：完美映射 CIM 标准，保持语义一致性
4. **统一接口**：POSIX 风格的系统调用接口，符合用户习惯
5. **可扩展性**：新增资源类型只需在 YAML 中定义
6. **并发安全**：基于 Go channel 和 sync 包，保证并发安全

## 技术栈

- **Go**: 主要编程语言
- **CIM**: Common Information Model，电力系统语义建模标准
- **Actor 模型**: 并发编程模型，用于资源抽象
- **YAML**: 类型系统定义格式

## 相关文档

- [类型系统使用文档](kernel/README.md)
- [Actor系统设计讨论](docs/actors_system_discussion_20251209.md)
- [POSIX风格资源服务层设计](docs/posix_resource_service_layer_20251210.md)
- [开发总结 - 2025年12月16日](docs/development_summary_20251216.md)

## 开发指南

### 添加新的 Control Capacity

1. 在 `actors/capacities/` 下创建新的 Capacity 实现
2. 实现 `Capacity` 接口
3. 在 `actors/capacity_factory.go` 中注册

### 添加新的 Measurement Capacity

1. 继承 `BaseMeasurementCapacity`
2. 实现具体的测量值处理逻辑
3. 在 `actors/capacity_factory.go` 中注册

### 定义新的资源类型

1. 在 `kernel/typesystem.yaml` 中添加类型定义
2. 定义属性、关系、能力
3. 在 `kernel/typesystem.yaml` 中定义 ioctl 命令映射

## 许可证

[待添加]
