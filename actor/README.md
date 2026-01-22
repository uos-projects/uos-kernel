# 轻量级 Go-Actor 系统

一个简单、轻量级的 Actor 模型实现，适用于 Go 语言的并发编程。特别设计用于映射 CIM (Common Information Model) 中的 PowerSystemResource 和 Control 关系。

## 核心概念

### Actor
Actor 是一个独立的并发实体，它：
- 拥有自己的邮箱（mailbox）用于接收消息
- 顺序处理消息（一次处理一条）
- 通过消息传递与其他 Actor 通信
- 封装状态和行为

### PowerSystemResourceActor（资源 Actor）
代表 CIM 模型中的 PowerSystemResource，是 Actor 的特殊实现：
- 每个 PowerSystemResource 对应一个 Actor
- 可以动态添加和移除能力（Capacity）
- 消息自动路由到对应的能力处理

### Capacity（能力）
代表 CIM 模型中的 Control，是 Actor 的能力抽象：
- 每个 Control 类型对应一个 Capacity 实现
- 通过 interface 定义，支持类型安全的能力扩展
- 能力与 PowerSystemResource 的关联关系决定 Actor 具有哪些能力

### 消息传递
- 异步、非阻塞的消息传递
- 基于 Go channel 实现
- 支持任意类型的消息
- 消息自动路由到对应的 Capacity

## 基本用法

### 基础 Actor

```go
package main

import (
    "context"
    "github.com/uos-projects/uos-kernel/actor"
)

// 创建资源 Actor（所有 Actor 都是 ResourceActor）
func main() {
    ctx := context.Background()
    system := actors.NewSystem(ctx)
    defer system.Shutdown()
    
    // 创建资源 Actor
    actor := actors.NewBaseResourceActor("my-actor", "MyResourceType")
    
    // 注册到系统
    if err := system.Register(actor); err != nil {
        panic(err)
    }
    
    // 发送消息
    system.Send("my-actor", "Hello, Actor!")
}
```

### PowerSystemResourceActor（资源 Actor）

```go
package main

import (
    "context"
    "github.com/uos-projects/uos-kernel/actor"
)

func main() {
    ctx := context.Background()
    system := actors.NewSystem(ctx)
    defer system.Shutdown()
    
    // 创建资源 Actor（对应 PowerSystemResource）
    actor := actors.NewPowerSystemResourceActor(
        "BE-Line_2",
        "ACLineSegment",
        nil,
    )
    
    // 创建能力工厂
    factory := actors.NewCapacityFactory()
    
    // 根据 CIM 关联关系添加能力
    // 如果 Control 关联到 PowerSystemResource，则添加对应能力
    capacity, err := factory.CreateCapacity("AccumulatorReset", "ACC_RESET_1")
    if err != nil {
        panic(err)
    }
    actor.AddCapacity(capacity)
    
    // 注册到系统
    system.Register(actor)
    
    // 发送消息（自动路由到 AccumulatorResetCapacity）
    system.Send("BE-Line_2", &actors.AccumulatorResetMessage{Value: 100})
}
```

## 自定义 Actor 行为

### 实现自定义 Capacity

```go
// 1. 定义消息类型
type MyControlMessage struct {
    Value int
}

// 2. 实现 Capacity 接口
type MyControlCapacity struct {
    actors.BaseCapacity
}

func NewMyControlCapacity(controlID string) *MyControlCapacity {
    return &MyControlCapacity{
        BaseCapacity: actors.BaseCapacity{
            controlID: controlID,
            name:      "MyControlCapacity",
        },
    }
}

func (c *MyControlCapacity) CanHandle(msg actors.Message) bool {
    _, ok := msg.(*MyControlMessage)
    return ok
}

func (c *MyControlCapacity) Execute(ctx context.Context, msg actors.Message) error {
    myMsg, ok := msg.(*MyControlMessage)
    if !ok {
        return fmt.Errorf("invalid message type")
    }
    
    // 实现具体的业务逻辑
    fmt.Printf("Executing MyControl with value: %d\n", myMsg.Value)
    return nil
}

// 3. 注册到工厂（可选）
func init() {
    actors.RegisterCapacityType("MyControl", func(controlID string) capacities.Capacity {
        return NewMyControlCapacity(controlID)
    })
}
```

## 特性

- ✅ 轻量级：最小化依赖，核心代码简洁
- ✅ 类型安全：利用 Go 的类型系统和 interface
- ✅ 并发安全：基于 Go channel 和 sync 包
- ✅ 易于扩展：清晰的接口设计
- ✅ 生命周期管理：支持启动和停止
- ✅ 能力模型：支持动态添加和移除能力（Capacity）
- ✅ 消息路由：自动将消息路由到对应的 Capacity
- ✅ CIM 映射：完美映射 CIM 模型中的 PowerSystemResource 和 Control 关系

## 架构

```
System
  ├── Actor Registry (map[string]Actor)
  └── Actor
      ├── BaseActor
      │   ├── Mailbox (chan Message)
      │   └── Context (context.Context)
      └── PowerSystemResourceActor
          ├── BaseActor (继承)
          ├── ResourceID
          ├── ResourceType
          └── Capabilities (map[string]Capacity)
              ├── AccumulatorResetCapacity
              ├── CommandCapacity
              ├── RaiseLowerCommandCapacity
              └── SetPointCapacity
```

## CIM 模型映射

### 映射关系

- **PowerSystemResource** → `PowerSystemResourceActor`
- **Control** → `Capacity` (interface)
- **Control 子类** → 具体的 `Capacity` 实现
  - `AccumulatorReset` → `AccumulatorResetCapacity`
  - `Command` → `CommandCapacity`
  - `RaiseLowerCommand` → `RaiseLowerCommandCapacity`
  - `SetPoint` → `SetPointCapacity`

### 关联关系

- 如果 `Control` 关联到 `PowerSystemResource`（通过 `Control.PowerSystemResource` 属性），则对应的 Actor 具有该 Control 对应的 Capacity
- 能力可以通过 `AddCapacity()` 动态添加
- 消息根据类型自动路由到对应的 Capacity 处理

## 已实现的功能

- ✅ 基础 Actor 系统
- ✅ PowerSystemResourceActor（资源 Actor）
- ✅ Capacity 接口和实现
- ✅ CapacityFactory（能力工厂）
- ✅ 消息路由到 Capacity
- ✅ 动态能力管理（添加/移除/查询）

## 未来扩展

- [ ] 从 CIM 数据自动创建 Actor 和 Capacity
- [ ] Actor 监控和指标
- [ ] 消息超时和重试
- [ ] Actor 监督（Supervision）
- [ ] 分布式 Actor 支持
- [ ] 消息优先级
- [ ] Actor 池（Pool）
- [ ] Capacity 的持久化

