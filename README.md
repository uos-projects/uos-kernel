# UOS Kernel - 电力系统资源操作系统内核

## 项目概述

UOS Kernel 是一个基于 CIM (Common Information Model) 标准的电力系统资源操作系统内核。它将电力系统中的资源抽象为 Actor，实现了概念上的驱动层，并在此基础上提供了 POSIX 风格的资源访问接口（ROSIX - Resource Operating System Interface）。

## 核心设计思想

### 1. CIM 模型到 Actor 系统的映射

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

### 2. 驱动层概念

Actor + Capability 组合模型类似于操作系统中的设备驱动：

```
┌─────────────────────────────────────────┐
│          ROSIX 服务层                    │
│  (POSIX风格资源访问接口)                 │
│  - Open/Close/Read/Write/RCtl           │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│          驱动层 (Driver Layer)           │
│                                          │
│  ┌──────────────────────────────────┐   │
│  │  PowerSystemResourceActor        │   │
│  │  (资源驱动)                       │   │
│  │  ├── ResourceID                  │   │
│  │  ├── ResourceType                │   │
│  │  └── Capabilities                │   │
│  │      ├── Control Capacity        │   │
│  │      │   ├── SetPoint            │   │
│  │      │   ├── Command             │   │
│  │      │   └── AccumulatorReset   │   │
│  │      └── Measurement Capacity   │   │
│  │          ├── Analog              │   │
│  │          ├── Discrete            │   │
│  │          └── Accumulator         │   │
│  └──────────────────────────────────┘   │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         CIM 模型层                        │
│  (Common Information Model)              │
│  - PowerSystemResource                  │
│  - Control                              │
│  - Measurement                           │
└─────────────────────────────────────────┘
```

**类比关系**：
- **Actor** ↔ **设备驱动**：Actor 代表一个资源实体，设备驱动代表一个硬件设备
- **Capability** ↔ **设备功能**：Capability 代表资源的能力（测量、控制），设备功能代表硬件的功能（读取、写入）
- **消息传递** ↔ **系统调用**：Actor 通过消息传递进行通信，设备驱动通过系统调用与内核交互
- **资源管理** ↔ **文件描述符**：通过统一的接口访问和管理资源

### 3. ROSIX - POSIX 风格的资源访问接口

在驱动层之上，提供了 POSIX 风格的资源访问接口（ROSIX），将资源抽象为可通过描述符访问的对象：

#### 核心接口

- **`Open(resourceID)` / `OpenWithExclusive(resourceID, exclusive)`**
  - 打开资源，返回资源描述符（ResourceDescriptor）
  - 每次 Open 都分配新的描述符，符合 POSIX 语义
  - 支持排他性和可复用性两种资源类型

- **`Close(fd)`**
  - 关闭资源描述符
  - 减少资源的引用计数
  - 当引用计数为 0 时，排他性资源可以再次打开

- **`Read(fd)`**
  - 读取 Actor 的状态信息（ResourceID、ResourceType、Capabilities）
  - 返回资源的元数据，而非测量值

- **`Write(fd, request)`**
  - 改变 Actor 的状态（待完善）
  - 用于更新资源属性、配置等

- **`RCtl(fd, cmd, arg)`**
  - 资源控制接口（类似 POSIX `ioctl`）
  - 将控制命令转换为 Capacity 消息并发送给 Actor
  - 支持各种控制操作（SetPoint、Command、AccumulatorReset 等）

#### 资源描述符模型

```
ResourceDescriptor (int32)
    │
    ├── ResourceHandle
    │   └── Resource (资源本身)
    │       ├── resourceID (string)
    │       ├── actor (*PowerSystemResourceActor)
    │       ├── refCount (int32)      # 引用计数
    │       └── exclusive (bool)       # 排他性标志
    │
    └── 多个描述符可以指向同一个资源
```

## 架构实现

### 目录结构

```
uos-kernel/
├── actors/                    # Actor 系统实现
│   ├── actor.go              # 基础 Actor
│   ├── system.go             # Actor 系统管理器
│   ├── resource_actor.go     # PowerSystemResourceActor
│   ├── capacity_factory.go   # Capacity 工厂
│   └── capacities/           # Capacity 实现
│       ├── base.go           # 基础 Capacity
│       ├── measurement_base.go  # Measurement 基类
│       ├── analog_measurement.go
│       ├── discrete_measurement.go
│       ├── accumulator_measurement.go
│       ├── set_point.go
│       ├── command.go
│       └── accumulator_reset.go
│
├── resource/                  # ROSIX 服务层
│   ├── manager.go            # 资源描述符管理（Open/Close）
│   ├── read_write.go         # 资源状态访问（Read/Write）
│   ├── rctl.go               # 资源控制（RCtl）
│   └── example_test.go       # 使用示例
│
├── ontology/                 # CIM 语义模型
│   └── cim_scope.yaml        # CIM 范围定义
│
└── docs/                     # 设计文档
    └── posix_resource_service_layer_20251210.md
```

### 核心组件

#### 1. Actor 系统 (`actors/`)

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

#### 2. ROSIX 服务层 (`resource/`)

**ResourceManager**
- 管理资源描述符的分配和回收
- 维护资源引用计数和排他性控制
- 提供 Open/Close/Read/Write/RCtl 接口

**资源描述符管理**
- 每次 Open 分配新的 ResourceDescriptor
- Resource 维护引用计数，支持排他性和可复用性
- Close 时减少引用计数，保留 Resource 以便快速重新打开

**资源控制**
- RCtl 将控制命令转换为 Capacity 消息
- 通过 Actor.Send() 发送消息到 Actor 的 mailbox
- Actor 自动路由消息到对应的 Capacity

## 使用示例

### 基本使用

```go
// 1. 创建 Actor 系统
ctx := context.Background()
system := actors.NewSystem(ctx)
defer system.Shutdown()

// 2. 创建资源 Actor
actor := actors.NewPowerSystemResourceActor("BE-G4", "SynchronousMachine", nil)

// 3. 添加 Capabilities
factory := actors.NewCapacityFactory()
setPointCap, _ := factory.CreateCapacity("SetPoint", "SET_PNT_1")
actor.AddCapacity(setPointCap)

system.Register(actor)

// 4. 创建资源管理器（ROSIX）
rm := resource.NewResourceManager(system)

// 5. 打开资源
fd, err := rm.Open("BE-G4")
defer rm.Close(fd)

// 6. 读取资源状态
state, err := rm.Read(ctx, fd)
fmt.Printf("Resource ID: %s\n", state.ResourceID)
fmt.Printf("Capabilities: %v\n", state.Capabilities)

// 7. 使用 RCtl 发送控制命令
setPointArg := map[string]interface{}{"value": 150.0}
rm.RCtl(ctx, fd, resource.CMD_SET_POINT, setPointArg)
```

### 排他性资源

```go
// 以排他性方式打开资源
fd1, err := rm.OpenWithExclusive("EXCLUSIVE-RESOURCE", true)
defer rm.Close(fd1)

// 再次打开会失败（资源已被排他性打开）
fd2, err := rm.OpenWithExclusive("EXCLUSIVE-RESOURCE", true)
if err != nil {
    fmt.Printf("Failed: %v\n", err)
}

// 关闭后可以再次打开
rm.Close(fd1)
fd3, err := rm.OpenWithExclusive("EXCLUSIVE-RESOURCE", true)
```

## 设计优势

1. **符合领域模型**：完美映射 CIM 标准，保持语义一致性
2. **类型安全**：利用 Go 的类型系统和 interface，确保类型安全
3. **可扩展性**：新增 Control 或 Measurement 类型只需实现 Capacity 接口
4. **统一接口**：ROSIX 提供统一的资源访问接口，隐藏底层 Actor 系统复杂性
5. **并发安全**：基于 Go channel 和 sync 包，保证并发安全
6. **资源管理**：支持排他性和可复用性两种资源类型，灵活的资源生命周期管理

## 技术栈

- **Go**: 主要编程语言
- **CIM**: Common Information Model，电力系统语义建模标准
- **Actor 模型**: 并发编程模型，用于资源抽象

## 相关文档

- [POSIX风格资源服务层设计文档](docs/posix_resource_service_layer_20251210.md)
- [Actor系统设计讨论](docs/actors_system_discussion_20251209.md)
- [CIM模型映射分析](docs/cimpyorm_fullgrid_analysis_20251209.md)

## 开发指南

### 添加新的 Control Capacity

1. 在 `actors/capacities/` 下创建新的 Capacity 实现
2. 实现 `Capacity` 接口
3. 在 `actors/capacity_factory.go` 中注册

### 添加新的 Measurement Capacity

1. 继承 `BaseMeasurementCapacity`
2. 实现具体的测量值处理逻辑
3. 在 `actors/capacity_factory.go` 中注册

### 扩展 ROSIX 接口

1. 在 `resource/rctl.go` 中添加新的控制命令
2. 实现对应的消息构造函数
3. 更新使用示例

## 许可证

[待添加]
