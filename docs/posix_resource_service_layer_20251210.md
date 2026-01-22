# POSIX风格资源服务层设计与实现

**日期**: 2025年12月10日

## 概述

在现有的Actor系统之上，设计并实现了一个POSIX风格的外部服务层，提供统一的资源访问接口。该服务层将Actor系统中的资源（PowerSystemResourceActor）抽象为可通过描述符访问的资源，类似于操作系统中的文件描述符机制。

## 背景与动机

### Actor + Capability 模型与设备驱动的类比

Actor + Capability组合模型与操作系统中的设备驱动有很强的可类比性：

- **Actor** ↔ **设备驱动**：Actor代表一个资源实体（如发电机、变压器），设备驱动代表一个硬件设备
- **Capability** ↔ **设备功能**：Capability代表资源的能力（如测量、控制），设备功能代表硬件的功能（如读取、写入）
- **消息传递** ↔ **系统调用**：Actor通过消息传递进行通信，设备驱动通过系统调用与内核交互
- **资源管理** ↔ **文件描述符**：需要统一的接口来访问和管理这些资源

### 现有系统的优化

在实现服务层之前，我们发现并优化了`BaseMeasurementCapacity`中的消息传递冗余：

**问题**：`StartConsumeLoop`中订阅到MQ消息后，通过`actorRef.Send()`发送消息给Actor，但实际上这个Actor就是自己，存在不必要的消息传递开销。

**解决方案**：
- 在`BaseMeasurementCapacity`中添加`capacityRef`和`ctx`字段
- 优先直接调用`Capacity.Execute()`，避免通过Actor mailbox的额外传递
- 保留通过Actor mailbox的回退机制，确保向后兼容

## 设计目标

基于POSIX接口设计，实现以下功能：

1. **资源描述符管理**：为每个Actor对应的资源分配一个资源描述符，可以通过这个描述符定位资源
2. **资源访问控制**：提供`Open`/`Close`接口，获取/释放资源访问
3. **资源状态访问**：提供`Read`/`Write`接口，访问资源状态或改变资源状态
4. **资源控制**：提供`RCtl`接口（仿照POSIX的`ioctl`），实现资源控制（调用Actor的Control Capacity）

## 实现细节

### 1. 资源描述符管理 (`resource/manager.go`)

#### ResourceDescriptor
```go
type ResourceDescriptor int32

const (
    InvalidDescriptor ResourceDescriptor = -1
    MinDescriptor     ResourceDescriptor = 0
    MaxDescriptor     ResourceDescriptor = 0x7FFFFFFF
)
```

#### Resource（资源本身）
每个`resourceID`对应一个`Resource`，管理资源本身的引用计数和排他性：
```go
type Resource struct {
    resourceID string
    actor      *actors.PowerSystemResourceActor
    refCount   int32      // 引用计数（有多少个描述符打开了这个资源）
    exclusive  bool       // 是否排他性资源（true表示只能被一个描述符打开）
    mu         sync.RWMutex
}
```

#### ResourceHandle（描述符句柄）
每个描述符对应一个`ResourceHandle`，指向一个`Resource`：
```go
type ResourceHandle struct {
    descriptor ResourceDescriptor
    resource   *Resource // 指向资源
}
```

#### Open接口
- **设计决策**：每次`Open`都分配新的`ResourceDescriptor`，但`Resource`本身会维护引用计数
- **原因**：符合POSIX语义（每次`open()`返回不同的文件描述符），同时支持资源的排他性和可复用性
- **实现**：
  - 检查资源是否存在
  - 获取或创建`Resource`
  - 检查排他性：如果资源是排他性的且`refCount > 0`，则拒绝打开
  - 分配新的描述符（每次Open都分配新的）
  - 创建`ResourceHandle`并映射到描述符
  - 增加`Resource`的引用计数

#### OpenWithExclusive接口
- 支持指定资源是否为排他性资源
- 如果资源已存在且新请求是排他性的，会更新资源的`exclusive`标志

#### Close接口
- **设计决策**：删除`ResourceHandle`，减少`Resource`的引用计数，但保留`Resource`在`resources` map中
- **原因**：
  - 当引用计数变为0时，排他性资源可以再次打开（`Open`检查`refCount > 0`）
  - 保留`Resource`可以避免重新查找Actor，快速重新打开资源
  - 符合POSIX语义：关闭描述符后，资源可以再次打开
- **实现**：
  - 减少`Resource`的引用计数
  - 删除`ResourceHandle`
  - 检查引用计数异常（不应该小于0）

### 2. 资源状态访问 (`resource/read_write.go`)

#### Read接口
读取Actor的状态信息（不是Measurement值）：
```go
type ActorState struct {
    ResourceID   string
    ResourceType string
    Capabilities []string
}
```

#### Write接口
改变Actor的状态（当前为占位实现，待后续完善）：
```go
type WriteRequest struct {
    Updates map[string]interface{}
}
```

### 3. 资源控制 (`resource/rctl.go`)

#### ControlCommand
定义控制命令类型：
```go
const (
    // 通用命令
    CMD_GET_RESOURCE_INFO ControlCommand = iota + 0x1000
    CMD_LIST_CAPABILITIES

    // Control 命令
    CMD_ACCUMULATOR_RESET
    CMD_COMMAND
    CMD_RAISE_LOWER
    CMD_SET_POINT
)
```

#### RCtl接口
- **设计决策**：将参数转换为Control消息，然后通过`Actor.Send()`发送给对应的Actor
- **实现流程**：
  1. 根据命令类型构造对应的Control消息（如`AccumulatorResetMessage`、`SetPointMessage`等）
  2. 通过`Actor.Send()`发送消息到Actor的mailbox
  3. Actor的`Receive()`方法会自动路由消息到对应的Capacity
  4. Capacity处理消息并执行相应的控制操作

#### 辅助函数
提供了多个辅助函数用于构造不同类型的Control消息：
- `buildAccumulatorResetMessage`
- `buildCommandMessage`
- `buildRaiseLowerMessage`
- `buildSetPointMessage`

## 代码结构

```
resource/
├── manager.go          # 资源描述符管理（Open/Close）
├── read_write.go       # 资源状态访问（Read/Write）
├── rctl.go            # 资源控制（RCtl）
└── example_test.go    # 使用示例
```

## 使用示例

### 基本使用（可复用资源）

```go
// 1. 创建Actor系统和资源Actor
system := actors.NewSystem(ctx)
actor := actors.NewPowerSystemResourceActor("BE-G4", "SynchronousMachine", nil)
system.Register(actor)

// 2. 创建资源管理器
rm := resource.NewResourceManager(system)

// 3. 打开资源（每次Open都返回新的描述符）
fd1, err := rm.Open("BE-G4")
defer rm.Close(fd1)

fd2, err := rm.Open("BE-G4")  // 同一个资源可以多次打开
defer rm.Close(fd2)

fmt.Printf("First descriptor: %d\n", fd1)   // 输出: 0
fmt.Printf("Second descriptor: %d\n", fd2) // 输出: 1

// 4. 读取Actor状态
state, err := rm.Read(ctx, fd1)
fmt.Printf("Resource ID: %s\n", state.ResourceID)
fmt.Printf("Capabilities: %v\n", state.Capabilities)

// 5. 使用rctl获取资源信息
info, _ := rm.RCtl(ctx, fd1, resource.CMD_GET_RESOURCE_INFO, nil)

// 6. 使用rctl发送SetPoint控制命令
setPointArg := map[string]interface{}{"value": 150.0}
rm.RCtl(ctx, fd1, resource.CMD_SET_POINT, setPointArg)
```

### 排他性资源

```go
// 1. 以排他性方式打开资源
fd1, err := rm.OpenWithExclusive("EXCLUSIVE-RESOURCE", true)
defer rm.Close(fd1)

// 2. 尝试再次打开排他性资源，应该失败
fd2, err := rm.OpenWithExclusive("EXCLUSIVE-RESOURCE", true)
if err != nil {
    fmt.Printf("Failed: %v\n", err) // 输出: resource EXCLUSIVE-RESOURCE is exclusive and already opened
}

// 3. 关闭第一个描述符后，可以再次打开
rm.Close(fd1)
fd3, err := rm.OpenWithExclusive("EXCLUSIVE-RESOURCE", true) // 成功
defer rm.Close(fd3)
```

## 设计决策说明

### 1. Open为什么每次都要分配新的描述符？
**回答**：每次`Open`都会分配新的`ResourceDescriptor`，符合POSIX语义（每次`open()`返回不同的文件描述符）。但`Resource`本身会维护引用计数，用于：
- **排他性控制**：排他性资源（`exclusive=true`）在`refCount > 0`时拒绝再次打开
- **可复用性**：非排他性资源可以被多个描述符同时打开
- **资源生命周期管理**：通过引用计数跟踪资源的使用情况

**关键区别**：
- `ResourceDescriptor`（描述符）：每次Open都不同，用于标识一个打开的资源实例
- `ResourceID`（资源ID）：字符串类型，用于标识资源本身
- `Resource`（资源对象）：管理资源本身的引用计数和排他性

### 2. 为什么Resource和ResourceHandle要分离？
**回答**：
- **Resource**：代表资源本身，一个`resourceID`对应一个`Resource`，管理引用计数和排他性
- **ResourceHandle**：代表一个打开的资源实例，一个描述符对应一个`Handle`，指向`Resource`
- **优势**：
  - 多个描述符可以指向同一个资源（支持可复用资源）
  - 资源本身的状态（引用计数、排他性）独立于描述符
  - 符合POSIX的文件描述符模型

### 3. Close为什么不删除Resource？
**回答**：`Close`只删除`ResourceHandle`，但保留`Resource`在`resources` map中，原因：
- **快速重新打开**：避免重新查找Actor，提高性能
- **排他性资源支持**：当`refCount=0`时，排他性资源可以再次打开（`Open`检查`refCount > 0`）
- **符合POSIX语义**：关闭描述符后，资源可以再次打开

### 4. Read读取的是什么？
**回答**：Read读取的是Actor的状态（ResourceID、ResourceType、Capabilities），不是Actor的Measurement值。Write也是改变Actor的状态，不是Measurement值。

### 5. RCtl如何工作？
**回答**：RCtl将参数转换为Control消息，然后通过`Actor.Send()`发送给对应的Actor。Actor的`Receive()`方法会根据消息类型路由到对应的Capacity进行处理。

## 设计演进

### 初始设计（已调整）

最初的设计中，`Open`会复用已存在的描述符并增加引用计数。但经过讨论，发现这个设计混淆了`ResourceID`（string类型）和`ResourceDescriptor`（int32类型）的概念。

### 最终设计

**关键改进**：
1. **每次Open都分配新的描述符**：符合POSIX语义，每次`open()`返回不同的文件描述符
2. **Resource和ResourceHandle分离**：
   - `Resource`：代表资源本身，管理引用计数和排他性
   - `ResourceHandle`：代表一个打开的资源实例，映射描述符到Resource
3. **引用计数在Resource层面**：用于控制排他性和资源生命周期，而不是在Handle层面

**设计优势**：
- 符合POSIX文件描述符模型
- 支持排他性和可复用性两种资源类型
- 清晰的资源生命周期管理
- 高效的资源重新打开（保留Resource对象）

## 待完善功能

1. **Write接口实现**：当前Write接口为占位实现，需要完善具体的状态更新逻辑
2. **Measurement相关命令**：扩展RCtl支持Measurement Capacity的查询命令（如`CMD_GET_MEASUREMENT_VALUE`、`CMD_GET_STATE_HISTORY`）
3. **错误处理**：完善错误处理和日志记录
4. **资源清理策略**：考虑添加资源清理策略（如长时间未使用的Resource自动清理）
5. **并发安全**：进一步优化并发访问性能

## 相关文件

- `resource/manager.go` - 资源描述符管理
- `resource/read_write.go` - 资源状态访问
- `resource/rctl.go` - 资源控制
- `resource/example_test.go` - 使用示例
- `actor/resource_actor.go` - Actor实现
- `actor/capacities/measurement_base.go` - Measurement Capacity基类

## 总结

本次实现为Actor系统提供了一个POSIX风格的外部服务层，使得外部系统可以通过统一的接口访问和管理Actor资源。该设计借鉴了操作系统中的文件描述符机制，提供了直观的资源访问方式，同时保持了Actor系统的消息传递模型和并发安全特性。

