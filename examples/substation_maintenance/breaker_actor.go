package main

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/uos-projects/uos-kernel/actors"
)

// BreakerActor 断路器 Actor
// 代表一个真实的断路器设备，长期存在，持续监测状态
type BreakerActor struct {
	*actors.BaseResourceActor

	// 设备属性
	name string // 设备名称

	// 设备状态
	isOpen      bool
	voltage     float64
	current     float64
	temperature float64 // 温度（用于异常检测）

	// 状态监测
	lastMaintenanceTime time.Time // 上次检修时间
	operationHours      int64     // 运行小时数
	startTime           time.Time // 启动时间

	// 异常检测阈值
	maxTemperature    float64
	maxOperationHours int64 // 最大运行小时数（检修间隔）

	// 状态锁
	mu sync.RWMutex

	// 监测协程控制
	monitorCtx    context.Context
	monitorCancel context.CancelFunc
}

// NewBreakerActor 创建断路器 Actor
func NewBreakerActor(id string, name string) *BreakerActor {
	now := time.Now()
	actor := &BreakerActor{
		BaseResourceActor: actors.NewBaseResourceActor(id, "Breaker"),
		name:                name,
		isOpen:              false,                 // 初始状态：关闭
		voltage:             220.0,                 // 初始电压：220kV
		current:             100.0,                 // 初始电流：100A
		temperature:         45.0,                  // 初始温度：45°C
		lastMaintenanceTime: now.AddDate(0, -2, 0), // 2个月前检修过
		operationHours:      0,
		startTime:           now,
		maxTemperature:      80.0, // 最大温度：80°C
		maxOperationHours:   1440, // 最大运行小时数：60天（1440小时）
	}

	// 设置设备属性（通过消息驱动）
	props := map[string]interface{}{
		"name":                name,
		"isOpen":              false,
		"voltage":             220.0,
		"current":             100.0,
		"temperature":         45.0,
		"lastMaintenanceTime": actor.lastMaintenanceTime,
	}
	for k, v := range props {
		msg := &actors.SetPropertyMessage{Name: k, Value: v}
		actor.Send(msg)
	}

	// 注册断路器开关控制能力（BreakerSwitchingCapacity）
	switchingCap := NewBreakerSwitchingCapacity(actor)
	actor.AddCapacity(switchingCap)

	// 注册业务事件（参考 Capacity 管理）
	actor.registerBusinessEvents()

	// 添加模拟设备绑定（通过 Binding 模拟真实世界反馈）
	binding := NewSimulatedBreakerBinding(actor.ResourceID())
	if err := actor.AddBinding(binding); err != nil {
		fmt.Printf("failed to add simulated breaker binding: %v\n", err)
	}

	return actor
}

// registerBusinessEvents 注册业务事件
func (b *BreakerActor) registerBusinessEvents() {
	// 注册设备异常事件
	deviceAbnormalEventDesc := actors.NewEventDescriptor(
		"DeviceAbnormalEvent",
		actors.EventTypeStateChanged,
		reflect.TypeOf((*DeviceAbnormalEvent)(nil)).Elem(),
		"设备异常事件（温度异常、电压异常等）",
		b.ResourceID(),
	)
	b.RegisterEvent(deviceAbnormalEventDesc)

	// 注册需要检修事件
	maintenanceRequiredEventDesc := actors.NewEventDescriptor(
		"MaintenanceRequiredEvent",
		actors.EventTypeStateChanged,
		reflect.TypeOf((*MaintenanceRequiredEvent)(nil)).Elem(),
		"需要检修事件（定期检修、异常检修等）",
		b.ResourceID(),
	)
	b.RegisterEvent(maintenanceRequiredEventDesc)
}

// Start 启动 Actor（重写，启动状态监测）
func (b *BreakerActor) Start(ctx context.Context) error {
	// 调用基类 Start
	if err := b.BaseResourceActor.Start(ctx); err != nil {
		return err
	}

	// 启动状态监测协程
	b.monitorCtx, b.monitorCancel = context.WithCancel(ctx)
	go b.monitorStatus(b.monitorCtx)

	return nil
}

// Stop 停止 Actor
func (b *BreakerActor) Stop() error {
	if b.monitorCancel != nil {
		b.monitorCancel()
	}
	return b.BaseResourceActor.Stop()
}

// monitorStatus 持续监测设备状态（长期存在的监测任务）
func (b *BreakerActor) monitorStatus(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second) // 每2秒检查一次（演示用，实际可以更长）
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.mu.Lock()

			// 模拟温度变化（实际应该从传感器读取）
			if !b.isOpen {
				// 运行时温度会波动
				b.temperature += (float64(time.Now().Unix()%10) - 5) * 0.5
				if b.temperature < 40.0 {
					b.temperature = 40.0
				}
				if b.temperature > 90.0 {
					b.temperature = 90.0
				}
				// 通过消息更新属性
				msg := &actors.SetPropertyMessage{Name: "temperature", Value: b.temperature}
				b.Send(msg)

				// 更新运行小时数
				b.operationHours = int64(time.Since(b.startTime).Hours())
			}

			// 检查温度异常
			if b.temperature > b.maxTemperature {
				b.mu.Unlock()
				b.emitAbnormalEvent("temperature_high", map[string]interface{}{
					"temperature": b.temperature,
					"threshold":   b.maxTemperature,
				})
				b.mu.Lock()
			}

			// 检查运行时间（需要检修）
			hoursSinceMaintenance := int64(time.Since(b.lastMaintenanceTime).Hours())
			if hoursSinceMaintenance > b.maxOperationHours {
				b.mu.Unlock()
				b.emitMaintenanceRequiredEvent("scheduled", hoursSinceMaintenance)
				b.mu.Lock()
			}

			b.mu.Unlock()
		}
	}
}

// emitAbnormalEvent 发射异常事件
func (b *BreakerActor) emitAbnormalEvent(eventType string, details map[string]interface{}) {
	event := &DeviceAbnormalEvent{
		DeviceID:  b.ResourceID(),
		EventType: eventType,
		Severity:  "warning",
		Details:   details,
		Timestamp: time.Now(),
	}

	// 通过事件发射器发射事件
	if emitter := b.GetEventEmitter(); emitter != nil {
		_ = emitter.Emit(actors.Event{
			Type:    actors.EventTypeStateChanged,
			Payload: event,
		})
	}

	fmt.Printf("[%s] ⚠️  异常检测：%s, 详情：%v\n", b.ResourceID(), eventType, details)
}

// emitMaintenanceRequiredEvent 发射需要检修事件
func (b *BreakerActor) emitMaintenanceRequiredEvent(reason string, operationHours int64) {
	event := &MaintenanceRequiredEvent{
		DeviceID:            b.ResourceID(),
		Reason:              reason,
		LastMaintenanceTime: b.lastMaintenanceTime,
		OperationHours:      operationHours,
		Details: map[string]interface{}{
			"operationHours":    operationHours,
			"maxOperationHours": b.maxOperationHours,
		},
		Timestamp: time.Now(),
	}

	// 通过事件发射器发射事件
	if emitter := b.GetEventEmitter(); emitter != nil {
		_ = emitter.Emit(actors.Event{
			Type:    actors.EventTypeStateChanged,
			Payload: event,
		})
	}

	fmt.Printf("[%s] 📅 需要检修：%s, 运行小时数：%d\n", b.ResourceID(), reason, operationHours)
}

// Receive 重写消息处理逻辑
func (b *BreakerActor) Receive(ctx context.Context, msg actors.Message) error {
	// 先处理来自 Binding 的外部事件（设备反馈）
	if ext, ok := msg.(*actors.ExternalEventMessage); ok && ext.BindingType == actors.BindingTypeDevice {
		if ev, ok := ext.Event.(*BreakerDeviceEvent); ok {
			switch ev.Action {
			case "opened":
				return b.doOpen(ctx, ev.Reason, ev.Operator)
			case "closed":
				return b.doClose(ctx, ev.Reason, ev.Operator)
			}
		}
	}

	// 处理完成检修命令
	if cmd, ok := msg.(*CompleteMaintenanceCommand); ok {
		return b.handleCompleteMaintenanceCommand(ctx, cmd)
	}

	// 其他消息交给基类处理，由 BaseResourceActor 根据 Capacity 进行路由
	return b.BaseResourceActor.Receive(ctx, msg)
}

// handleCompleteMaintenanceCommand 处理完成检修命令
func (b *BreakerActor) handleCompleteMaintenanceCommand(ctx context.Context, cmd *CompleteMaintenanceCommand) error {
	b.CompleteMaintenance()
	return nil
}

// doOpen 执行断路器打开操作（Actor 内部领域方法）
// 注意：不直接暴露给外部，只由 Capacity 调用
func (b *BreakerActor) doOpen(ctx context.Context, reason, operator string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	fmt.Printf("[%s] 🔌 执行打开断路器操作，操作员：%s，原因：%s\n",
		b.ResourceID(), operator, reason)

	// 检查前置条件：断路器必须处于关闭状态
	if b.isOpen {
		return fmt.Errorf("断路器 %s 已经处于打开状态", b.ResourceID())
	}

	// 执行打开操作（模拟）
	time.Sleep(100 * time.Millisecond) // 模拟操作时间

	// 更新状态（通过消息驱动）
	b.isOpen = true
	b.current = 0.0 // 打开后电流为0
	b.Send(&actors.SetPropertyMessage{Name: "isOpen", Value: true})
	b.Send(&actors.SetPropertyMessage{Name: "current", Value: 0.0})

	fmt.Printf("[%s] ✓ 断路器已打开，电流：%.2f A\n", b.ResourceID(), b.current)

	return nil
}

// doClose 执行断路器关闭操作（Actor 内部领域方法）
// 注意：不直接暴露给外部，只由 Capacity 调用
func (b *BreakerActor) doClose(ctx context.Context, reason, operator string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	fmt.Printf("[%s] 🔌 执行关闭断路器操作，操作员：%s，原因：%s\n",
		b.ResourceID(), operator, reason)

	// 检查前置条件：断路器必须处于打开状态
	if !b.isOpen {
		return fmt.Errorf("断路器 %s 已经处于关闭状态", b.ResourceID())
	}

	// 执行关闭操作（模拟）
	time.Sleep(100 * time.Millisecond) // 模拟操作时间

	// 更新状态（通过消息驱动）
	b.isOpen = false
	b.current = 100.0 // 关闭后恢复电流
	b.Send(&actors.SetPropertyMessage{Name: "isOpen", Value: false})
	b.Send(&actors.SetPropertyMessage{Name: "current", Value: 100.0})

	fmt.Printf("[%s] ✓ 断路器已关闭，电流：%.2f A\n", b.ResourceID(), b.current)

	return nil
}

// CompleteMaintenance 完成检修（更新检修时间）
func (b *BreakerActor) CompleteMaintenance() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.lastMaintenanceTime = time.Now()
	b.startTime = time.Now() // 重置启动时间
	b.operationHours = 0
	b.temperature = 45.0 // 重置温度
	b.Send(&actors.SetPropertyMessage{Name: "lastMaintenanceTime", Value: b.lastMaintenanceTime})
	b.Send(&actors.SetPropertyMessage{Name: "temperature", Value: b.temperature})

	fmt.Printf("[%s] ✅ 检修完成，检修时间已更新\n", b.ResourceID())
}

// GetStatus 获取断路器状态
func (b *BreakerActor) GetStatus() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return map[string]interface{}{
		"id":                  b.ResourceID(),
		"name":                b.name,
		"isOpen":              b.isOpen,
		"voltage":             b.voltage,
		"current":             b.current,
		"temperature":         b.temperature,
		"lastMaintenanceTime": b.lastMaintenanceTime,
		"operationHours":      b.operationHours,
	}
}
