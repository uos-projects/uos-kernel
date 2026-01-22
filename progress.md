# UOS Kernel 开发进展

## 2026-01-22 - 消息系统简化与事件机制重构

### 主要改动

#### 1. 消息系统简化

**改动内容：**
- 删除了 `InternalMessage` 接口和所有内部消息类型（`SetPropertyMessage`、`TransitionStateMessage`、`SuspendMessage`、`ResumeMessage` 等）
- 重命名 `CapabilityCommand` → `Command`
- 重命名 `CoordinationEvent` → `Event`
- 消息类型从 3 种简化为 2 种：`Command` 和 `Event`

**影响文件：**
- `actor/message.go` - 消息定义简化
- `actor/base_resource_actor.go` - 删除内部消息处理逻辑
- `examples/substation_maintenance/commands.go` - 更新消息类型
- `examples/substation_maintenance/events.go` - 更新消息类型

#### 2. 事件发射器重新设计

**改动内容：**
- 删除了 `Event` 结构体（事件发射器用的那个）
- 删除了 `EventType` 常量
- `EventEmitter` 现在直接发射 `Event`（Message）类型
- 简化了 `BaseEventEmitter`，只保留 `EmitEvent(event Event)` 方法

**影响文件：**
- `actor/event.go` - 事件发射器重新设计
- `actor/base_resource_actor.go` - 更新事件发射调用
- `examples/substation_maintenance/*.go` - 更新所有事件发射代码

#### 3. System 事件订阅机制

**改动内容：**
- 在 `System` 中添加了 `SubscribeEvent(actorID string)` 方法
- 添加了 `UnsubscribeEvent(actorID string)` 方法
- `PublishEvent()` 现在直接通过 `system.Send()` 发送给所有订阅的 Actor
- 删除了 `EventListener` 接口，不再需要适配器类

**优势：**
- 更简单：无需定义适配器类
- 更直接：事件直接作为消息发送给 Actor
- 更符合 Actor 模型：事件就是消息

**影响文件：**
- `actor/system.go` - 添加事件订阅机制
- `examples/substation_maintenance/main.go` - 删除 `DispatcherEventListener`，直接使用 `system.SubscribeEvent()`

#### 4. 状态管理改为直接方法调用

**改动内容：**
- 添加了 `SetProperty(name string, value interface{})` 方法
- 添加了 `GetProperty(name string) (interface{}, bool)` 方法
- 添加了 `GetAllProperties() map[string]interface{}` 方法
- 添加了 `Suspend(reason string)` 和 `Resume(reason string)` 方法
- 删除了所有 `SetPropertyMessage` 的使用

**影响文件：**
- `actor/base_resource_actor.go` - 添加属性管理方法
- `examples/substation_maintenance/breaker_actor.go` - 所有 `SetPropertyMessage` 改为直接调用 `SetProperty()`

#### 5. Capacity Effects 简化

**改动内容：**
- `Effect.Events` 现在直接是 `Event`（Message）列表
- 删除了 `CapacityEvent` 类型
- `applyEffects()` 中的事件发射改为直接发送 Event（Message）

**影响文件：**
- `actor/capacity.go` - 简化 Effects 定义
- `actor/base_resource_actor.go` - 更新 `applyEffects()` 实现

#### 6. 修复 BaseActor.run() 调用子类 Receive 的问题

**问题：**
- `BaseActor.run()` 中 `a` 是 `*BaseActor` 类型，调用 `a.Receive()` 时调用的是 `BaseActor.Receive`（空实现），而不是子类的 `Receive`

**解决方案：**
- 在 `BaseActor` 中添加 `self Actor` 字段
- 在 `System.Register()` 时设置 `self` 引用
- `BaseActor.run()` 通过 `self.Receive()` 调用，确保调用子类的实现

**影响文件：**
- `actor/actor.go` - 添加 `self` 字段和 `SetSelf()` 方法
- `actor/system.go` - 在 `Register()` 时设置 `self` 引用
- `actor/base_resource_actor.go` - 添加 `GetBaseActor()` 方法

### 验证结果

运行 `examples/substation_maintenance` 示例，确认完整事件流程：

1. **Breaker → Dispatcher（事件）**
   - Breaker 发射 `MaintenanceRequiredEvent` 或 `DeviceAbnormalEvent`
   - 通过 `emitter.EmitEvent()` → `System.PublishEvent()` → `System.Send("DISPATCHER", event)`
   - Dispatcher 接收并处理事件，创建检修任务

2. **Dispatcher → Operator（命令）**
   - Dispatcher 发送 `StartMaintenanceCommand` 给 Operator
   - 通过 `system.Send("OP-001", cmd)` 直接发送
   - Operator 接收命令并执行检修操作

3. **Operator → Dispatcher（事件）**
   - Operator 发送 `MaintenanceCompletedEvent` 给 Dispatcher
   - 通过 `system.Send("DISPATCHER", event)` 直接发送
   - Dispatcher 接收并更新任务状态

### 技术要点

1. **事件订阅机制**：只有 `Dispatcher` 订阅了全局事件，所有 Breaker 发射的事件都会发送给它
2. **两种通信方式**：
   - 事件（Event）：通过 `EmitEvent()` → `System.PublishEvent()` → 发送给所有订阅者
   - 命令（Command）：通过 `system.Send()` 直接发送给指定 Actor
3. **消息路由**：`BaseActor.run()` 通过 `self` 引用调用子类的 `Receive` 方法，确保多态性

### 代码统计

- 修改文件数：约 15 个
- 删除代码行数：约 200+ 行（内部消息、Event 结构体等）
- 新增代码行数：约 100 行（属性管理方法、事件订阅机制等）
- 简化程度：消息类型从 3 种减少到 2 种，事件机制从适配器模式简化为直接发送

### 下一步计划

1. 考虑是否需要事件过滤机制（只订阅特定类型的事件）
2. 优化事件描述符的使用（目前注册但未充分利用）
3. 考虑是否需要事件持久化机制

---

## Actor 系统设计实现（基于调度示例）

### 系统架构概览

Actor 系统采用分层架构，核心组件包括：

```
System (系统层)
  ├── Actor Registry (Actor 注册表)
  ├── Event Subscribers (事件订阅者列表)
  └── Message Routing (消息路由)
      │
      └── Actor (Actor 层)
          ├── BaseActor (基础 Actor)
          │   ├── Mailbox (消息邮箱)
          │   ├── Context (执行上下文)
          │   └── Lifecycle (生命周期管理)
          │
          └── BaseResourceActor (资源 Actor)
              ├── BaseActor (嵌入)
              ├── Capabilities (能力管理)
              ├── Bindings (外部绑定)
              ├── Events (事件管理)
              ├── Properties (属性存储)
              └── EventEmitter (事件发射器)
```

### 核心组件设计

#### 1. Actor 基础层 (`BaseActor`)

**职责：**
- 消息邮箱管理（基于 Go channel）
- 消息顺序处理（单线程模型）
- 生命周期管理（启动/停止）
- 消息路由到子类 `Receive` 方法

**关键实现：**
```go
type BaseActor struct {
    id      string
    mailbox chan Message
    ctx     context.Context
    cancel  context.CancelFunc
    self    Actor  // 指向完整 Actor 实现的引用
    system  *System
}

// run() 方法通过 self.Receive() 调用子类实现
func (a *BaseActor) run() {
    for {
        select {
        case msg := <-a.mailbox:
            if err := a.self.Receive(a.ctx, msg); err != nil {
                // 错误处理
            }
        case <-a.ctx.Done():
            return
        }
    }
}
```

**设计要点：**
- 使用 `self Actor` 字段确保多态性，避免调用基类的空实现
- 消息邮箱是缓冲 channel，支持异步非阻塞发送
- 单线程消息处理保证状态一致性

#### 2. 资源 Actor 层 (`BaseResourceActor`)

**职责：**
- 能力（Capacity）管理：动态添加/移除能力
- 外部绑定（Binding）管理：与真实世界交互的适配层
- 事件管理：注册和发射业务事件
- 属性管理：键值对状态存储
- 生命周期状态机：`offline` → `online` → `suspended` → `busy`

**关键实现：**
```go
type BaseResourceActor struct {
    *BaseActor
    resourceID   string
    resourceType string
    capabilities map[string]Capacity
    bindings     map[BindingType]Binding
    events       map[string]*EventDescriptor
    properties   map[string]interface{}
    eventEmitter *BaseEventEmitter
}

// Receive 方法路由消息到 Capacity 或子类处理
func (a *BaseResourceActor) Receive(ctx context.Context, msg Message) error {
    switch msg.MessageType() {
    case MessageCategoryCommand:
        return a.handleCommand(ctx, msg)
    case MessageCategoryEvent:
        return a.handleEvent(ctx, msg)
    }
    return nil
}
```

**设计要点：**
- 嵌入 `BaseActor` 获得基础消息处理能力
- 通过 `AddCapacity()` 动态扩展功能
- 通过 `AddBinding()` 连接外部系统（设备协议/人机交互）
- 属性管理使用 `sync.RWMutex` 保证并发安全

#### 3. Capacity（能力）模式

**设计理念：**
- Capacity 封装特定的业务能力（如 `BreakerSwitchingCapacity`）
- 每个 Capacity 定义可接受的命令类型、前置条件、执行语义和副作用

**接口定义：**
```go
type Capacity interface {
    // CommandShape: 定义可接受的命令类型
    AcceptableCommands() []CommandType
    
    // Guards: 定义前置条件（如设备必须在线）
    Requires() []Requirement
    CheckGuards(ctx context.Context, actor interface{}) (bool, []string, error)
    
    // ExecutionSemantics: 定义执行模式（同步/异步/超时等）
    ExecutionInfo() ExecutionInfo
    
    // Effects: 定义副作用（状态变更、事件发射）
    Execute(ctx context.Context, actor interface{}, command Command) (*Effect, error)
}
```

**执行流程：**
1. `BaseResourceActor` 接收 `Command` 消息
2. 遍历所有 `Capacity`，找到能处理该命令的 Capacity
3. 检查前置条件（Guards）
4. 执行命令（Execute）
5. 应用副作用（Effects）：更新状态、发射事件

**示例：`BreakerSwitchingCapacity`**
```go
type BreakerSwitchingCapacity struct {
    breaker *BreakerActor
}

func (c *BreakerSwitchingCapacity) Execute(ctx context.Context, actor interface{}, cmd Command) (*Effect, error) {
    switch cmd := cmd.(type) {
    case *OpenBreakerCommand:
        // 1. 通过 Binding 执行真实操作
        binding := c.breaker.GetBinding(BindingTypeDevice)
        if err := binding.ExecuteExternal(ctx, cmd); err != nil {
            return nil, err
        }
        
        // 2. 更新状态
        c.breaker.SetProperty("isOpen", true)
        
        // 3. 返回副作用（发射事件）
        return &Effect{
            StateChanges: map[string]interface{}{"isOpen": true},
            Events: []Event{&BreakerOpenedEvent{...}},
        }, nil
    }
    return nil, nil
}
```

#### 4. Binding（绑定）模式

**设计理念：**
- Binding 是 Actor 与真实世界交互的适配层
- 将 Actor 的命令转换为外部协议（设备协议/人机交互/服务 API）
- 将外部事件转换为 Actor 消息

**接口定义：**
```go
type Binding interface {
    Type() BindingType  // device / human / service / mq
    ResourceID() string
    Start(ctx context.Context) error
    Stop() error
    
    // 外部事件 → Actor 消息
    OnExternalEvent(ctx context.Context, event interface{}) error
    
    // Actor 命令 → 外部操作
    ExecuteExternal(ctx context.Context, command interface{}) error
}
```

**示例：`SimulatedBreakerBinding`**
```go
// 模拟断路器设备协议
type SimulatedBreakerBinding struct {
    *BaseBinding
    breaker *BreakerActor
}

func (b *SimulatedBreakerBinding) ExecuteExternal(ctx context.Context, cmd interface{}) error {
    switch cmd := cmd.(type) {
    case *OpenBreakerCommand:
        // 模拟设备协议：发送控制信号
        fmt.Printf("[设备协议] 发送打开命令到设备 %s\n", b.ResourceID())
        
        // 模拟设备响应：发送状态更新事件
        event := &BreakerStateChangedEvent{
            DeviceID: b.ResourceID(),
            NewState: "open",
        }
        return b.SendToActor(event)  // 转换为消息发送给 Actor
    }
    return nil
}
```

**示例：`SimulatedOperatorBinding`**
```go
// 模拟操作员行为
type SimulatedOperatorBinding struct {
    *BaseBinding
    operator *DispatcherOperatorActor
}

func (b *SimulatedOperatorBinding) ExecuteExternal(ctx context.Context, cmd interface{}) error {
    switch cmd := cmd.(type) {
    case *StartMaintenanceCommand:
        // 模拟操作员执行检修操作
        fmt.Printf("[操作员] 开始执行检修任务：%s\n", cmd.TaskID)
        
        // 模拟操作步骤
        time.Sleep(2 * time.Second)
        
        // 发送完成事件
        event := &MaintenanceCompletedEvent{
            TaskID: cmd.TaskID,
            Result: "success",
        }
        return b.SendToActor(event)
    }
    return nil
}
```

#### 5. 事件系统

**事件类型：**
- **业务事件（Event）**：实现 `Message` 接口，通过 `EmitEvent()` 发射
- **系统事件**：生命周期变化等（暂未实现）

**事件流程：**
```
Actor 发射事件
  ↓
emitter.EmitEvent(event)
  ↓
System.PublishEvent(ctx, event)
  ↓
遍历所有订阅者
  ↓
System.Send(actorID, event)  // 直接作为消息发送
  ↓
订阅者 Actor 的 Receive() 方法处理
```

**事件订阅：**
```go
// 在 main.go 中
system.SubscribeEvent("DISPATCHER")  // Dispatcher 订阅所有事件

// 在 System.PublishEvent() 中
func (s *System) PublishEvent(ctx context.Context, event Event) {
    for _, actorID := range s.eventSubscribers {
        s.Send(actorID, event)  // 直接发送给订阅者
    }
}
```

**事件描述符：**
- 用于声明 Actor 可以发射的事件类型
- 目前主要用于文档和类型检查，未来可用于事件过滤

### 调度示例中的 Actor 实现

#### 1. BreakerActor（设备 Actor）

**职责：**
- 代表真实的断路器设备
- 持续监测状态（温度、运行时间）
- 检测异常并发射事件
- 响应开关命令

**关键特性：**
```go
type BreakerActor struct {
    *BaseResourceActor
    
    // 设备状态
    isOpen      bool
    temperature float64
    operationHours int64
    
    // 监测协程
    monitorCtx    context.Context
    monitorCancel context.CancelFunc
}

// 监测协程持续运行
func (b *BreakerActor) monitorStatus() {
    ticker := time.NewTicker(1 * time.Second)
    for {
        select {
        case <-ticker.C:
            // 检查温度异常
            if b.temperature > b.maxTemperature {
                b.emitAbnormalEvent("temperature_high", ...)
            }
            // 检查运行时间
            if b.operationHours > b.maxOperationHours {
                b.emitMaintenanceRequiredEvent(...)
            }
        case <-b.monitorCtx.Done():
            return
        }
    }
}
```

**事件发射：**
```go
func (b *BreakerActor) emitMaintenanceRequiredEvent(...) {
    event := &MaintenanceRequiredEvent{
        DeviceID: b.ResourceID(),
        Reason:   "scheduled",
        ...
    }
    emitter := b.GetEventEmitter()
    if emitter != nil {
        emitter.EmitEvent(event)  // → System.PublishEvent() → Dispatcher
    }
}
```

#### 2. DispatcherActor（调度中心 Actor）

**职责：**
- 接收设备事件（异常/检修需求）
- 创建检修任务
- 分配任务给操作员
- 跟踪任务状态

**关键实现：**
```go
type DispatcherActor struct {
    *BaseResourceActor
    
    maintenancePlans []MaintenancePlan
    pendingTasks     []MaintenanceTask
    operators        []string
    system           *System
}

// Receive 方法处理事件
func (d *DispatcherActor) Receive(ctx context.Context, msg Message) error {
    switch event := msg.(type) {
    case *MaintenanceRequiredEvent:
        return d.handleMaintenanceRequiredEvent(ctx, event)
    case *DeviceAbnormalEvent:
        return d.handleDeviceAbnormalEvent(ctx, event)
    case *MaintenanceCompletedEvent:
        return d.handleMaintenanceCompletedEvent(ctx, event)
    }
    return d.BaseResourceActor.Receive(ctx, msg)
}

// 创建任务并分配给操作员
func (d *DispatcherActor) handleMaintenanceRequiredEvent(ctx context.Context, event *MaintenanceRequiredEvent) error {
    // 1. 创建任务
    task := &MaintenanceTask{
        TaskID: generateTaskID(),
        Type:   "scheduled",
        Devices: []string{event.DeviceID},
        ...
    }
    
    // 2. 分配给操作员
    operatorID := d.selectOperator()
    task.AssignedTo = operatorID
    
    // 3. 发送命令给操作员
    cmd := &StartMaintenanceCommand{
        TaskID: task.TaskID,
        Devices: task.Devices,
        ...
    }
    d.system.Send(operatorID, cmd)  // 直接发送命令
    
    // 4. 发射任务分配事件
    d.emitTaskAssignedEvent(task)
    
    return nil
}
```

#### 3. DispatcherOperatorActor（操作员 Actor）

**职责：**
- 接收检修任务命令
- 通过 Binding 执行实际检修操作
- 发射检修完成事件

**关键实现：**
```go
type DispatcherOperatorActor struct {
    *BaseResourceActor
    
    operatorID   string
    operatorName string
    currentTask  *MaintenanceTask
    system       *System
}

// 接收开始检修命令
func (o *DispatcherOperatorActor) handleStartMaintenanceCommand(ctx context.Context, cmd *StartMaintenanceCommand) error {
    // 1. 更新状态
    o.currentTask = &MaintenanceTask{...}
    o.SetProperty("currentTask", o.currentTask)
    
    // 2. 通过 Binding 执行实际操作
    binding := o.GetBinding(BindingTypeHuman)
    if binding != nil {
        return binding.ExecuteExternal(ctx, cmd)  // → SimulatedOperatorBinding
    }
    
    return nil
}
```

**与 Binding 的交互：**
- Actor 只负责状态管理
- Binding 负责实际行为执行
- Binding 执行完成后，通过 `SendToActor()` 发送事件给 Actor
- Actor 接收事件后更新状态并通知调度中心

### 设计模式总结

#### 1. Actor 模型
- **消息传递**：Actor 之间通过消息通信，不共享状态
- **单线程处理**：每个 Actor 顺序处理消息，保证状态一致性
- **位置透明**：Actor 通过 ID 引用，不关心物理位置

#### 2. Capacity 模式
- **能力封装**：将功能封装为可插拔的 Capacity
- **命令路由**：消息自动路由到对应的 Capacity
- **前置条件**：通过 Guards 确保执行条件满足

#### 3. Binding 模式
- **适配器模式**：将 Actor 命令适配为外部协议
- **双向转换**：外部事件 → Actor 消息，Actor 命令 → 外部操作
- **可插拔**：支持多种绑定类型（设备/人/服务/消息队列）

#### 4. 事件驱动
- **发布-订阅**：Actor 发射事件，订阅者接收
- **解耦**：事件发射者不关心接收者
- **统一消息**：事件就是消息，通过 `System.Send()` 分发

#### 5. 两阶段设计
- **命令阶段**：Actor 接收命令，通过 Binding 执行外部操作
- **反馈阶段**：Binding 执行完成后，发送事件反馈给 Actor
- **状态同步**：Actor 根据反馈更新状态

### 最佳实践

1. **Actor 职责单一**：每个 Actor 代表一个资源或角色
2. **状态私有化**：状态只能通过消息修改，不直接访问
3. **事件驱动**：使用事件通知状态变化，而不是轮询
4. **Binding 分离**：将外部交互逻辑放在 Binding 中，Actor 只管理状态
5. **Capacity 封装**：将业务能力封装为 Capacity，便于复用和测试
6. **消息类型明确**：Command 用于操作，Event 用于通知

### 扩展性考虑

1. **事件过滤**：未来可以支持按事件类型过滤订阅
2. **分布式支持**：Actor ID 可以包含网络地址，支持跨节点通信
3. **持久化**：可以添加状态快照和事件日志，支持 Actor 恢复
4. **监控**：可以添加 Actor 指标收集和健康检查
5. **监督机制**：可以添加 Actor 监督树，处理 Actor 故障
