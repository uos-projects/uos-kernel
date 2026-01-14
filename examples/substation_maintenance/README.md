# 变电站停电检修操作场景示例

## 场景描述

变电站中有多个设备（主变压器、进线断路器、出线断路器）。系统会根据以下两种方式触发检修操作：

1. **定期检修计划**：设备运行时间达到检修间隔时自动触发
2. **异常检测**：设备状态异常（如温度过高）时触发紧急检修

调度中心接收检修需求，制定检修计划并分配给调度操作员。调度操作员执行停电检修操作，按顺序操作多个断路器，完成检修后恢复供电。

## 设计理念体现

### 1. Actor 是长期存在的运行时实体

- **设备 Actor**：代表真实的断路器设备，在系统中持续运行
- **调度中心 Actor**：持续监听设备事件，管理检修计划
- **操作员 Actor**：随时准备接收检修任务

```go
// 设备 Actor 一旦创建并启动，就会持续运行
breaker1 := NewBreakerActor("BREAKER-001", "主变压器进线断路器")
system.Register(breaker1)
breaker1.Start(ctx)  // Actor 开始运行，持续存在
```

### 2. Actor 维护资源状态

- **设备状态**：断路器打开/关闭、电压、电流、温度
- **运行时间**：记录运行小时数，用于判断是否需要检修
- **检修历史**：记录上次检修时间

```go
type BreakerActor struct {
    isOpen bool
    voltage float64
    current float64
    temperature float64
    lastMaintenanceTime time.Time
    operationHours int64
}
```

### 3. 事件驱动

- **异常检测**：设备监测到异常时发射 `DeviceAbnormalEvent`
- **定期检修**：运行时间达到阈值时发射 `MaintenanceRequiredEvent`
- **任务分配**：调度中心分配任务时发射 `MaintenanceTaskAssignedEvent`
- **检修完成**：操作员完成检修后发射 `MaintenanceCompletedEvent`

```go
// 设备监测到异常，发射事件
func (b *BreakerActor) emitAbnormalEvent(eventType string, details map[string]interface{}) {
    event := &DeviceAbnormalEvent{
        DeviceID:  b.ResourceID(),
        EventType: eventType,
        Severity:  "warning",
        Details:   details,
    }
    // 通过事件发射器发射事件
    emitter.Emit(event)
}
```

### 4. Actor 之间协同工作

- **设备 → 调度中心**：设备发射异常/检修需求事件
- **调度中心 → 操作员**：调度中心分配检修任务
- **操作员 → 设备**：操作员发送操作命令（打开/关闭断路器）
- **操作员 → 调度中心**：操作员完成检修后通知调度中心

```go
// 操作员协调多个设备 Actor 完成检修操作
func (o *DispatcherOperatorActor) executeMaintenanceOperation(ctx context.Context, task *MaintenanceTask) error {
    // 按顺序操作多个断路器
    for _, deviceID := range task.Devices {
        cmd := &OpenBreakerCommand{...}
        o.system.Send(deviceID, cmd)  // 发送命令到设备 Actor
    }
    // ...
}
```

### 5. 私有状态和消息驱动

- 设备状态只能通过消息修改，不能直接访问
- 所有状态变化都通过消息驱动

```go
// 通过消息更新状态
b.Send(&actors.SetPropertyMessage{Name: "isOpen", Value: true})
```

## 文件结构

```
examples/substation_maintenance/
├── main.go              # 主程序，演示完整场景
├── breaker_actor.go     # 断路器 Actor 实现（设备 Actor）
├── dispatcher_actor.go  # 调度中心 Actor 实现
├── operator_actor.go    # 调度操作员 Actor 实现
├── events.go            # 事件定义（Coordination Events）
├── commands.go          # 命令定义（Capability Commands）
└── README.md           # 本文件
```

## 运行示例

```bash
cd examples/substation_maintenance
go run .
```

## 运行流程

1. **初始化**：创建设备 Actor、调度中心 Actor、操作员 Actor
2. **启动**：所有 Actor 启动并开始运行
3. **监测**：设备 Actor 持续监测状态（温度、运行时间）
4. **触发**：
   - 定期检修：运行时间达到阈值 → 发射 `MaintenanceRequiredEvent`
   - 异常检测：温度异常 → 发射 `DeviceAbnormalEvent`
5. **响应**：调度中心接收事件，创建检修任务，分配给操作员
6. **执行**：操作员执行停电检修操作
   - 步骤1：停电（按顺序打开所有断路器）
   - 步骤2：执行检修操作
   - 步骤3：恢复供电（按顺序关闭所有断路器）
7. **完成**：操作员完成检修，通知调度中心，更新设备检修时间

## 关键特性

- ✅ **长期存在**：Actor 持续运行，代表真实资源
- ✅ **状态维护**：Actor 维护设备状态和运行历史
- ✅ **事件驱动**：通过事件触发和响应操作
- ✅ **协同工作**：多个 Actor 协同完成复杂业务流程
- ✅ **消息驱动**：所有状态变化通过消息驱动
- ✅ **异常检测**：自动检测异常并触发响应
- ✅ **定期计划**：支持定期检修计划
