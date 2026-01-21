package main

import (
	"context"
	"time"

	"github.com/uos-projects/uos-kernel/actors"
)

// BreakerDeviceEvent 表示来自设备侧的状态反馈事件
// 这是 Binding 从真实世界（或模拟设备）感知到的结果
type BreakerDeviceEvent struct {
	Action   string // "opened" 或 "closed"
	Operator string
	Reason   string
}

// SimulatedBreakerBinding 模拟断路器设备绑定
// 不直接修改 Actor 状态，而是在 ExecuteExternal 时：
// 1. 接收到命令
// 2. 延迟一段时间
// 3. 通过 ExternalEventMessage 把“设备状态变化”反馈给 Actor
type SimulatedBreakerBinding struct {
	*actors.BaseBinding
}

// NewSimulatedBreakerBinding 创建模拟绑定
func NewSimulatedBreakerBinding(resourceID string) *SimulatedBreakerBinding {
	base := actors.NewBaseBinding(actors.BindingTypeDevice, resourceID)
	return &SimulatedBreakerBinding{BaseBinding: base}
}

// Start 启动绑定（这里不需要额外逻辑）
func (b *SimulatedBreakerBinding) Start(ctx context.Context) error {
	return nil
}

// Stop 停止绑定
func (b *SimulatedBreakerBinding) Stop() error {
	return b.BaseBinding.Stop()
}

// OnExternalEvent 外部事件入口（本示例未使用，保留以符合接口）
func (b *SimulatedBreakerBinding) OnExternalEvent(ctx context.Context, event interface{}) error {
	// 直接转发为 ExternalEventMessage
	msg := &actors.ExternalEventMessage{
		BindingType: actors.BindingTypeDevice,
		Event:       event,
	}
	return b.SendToActor(msg)
}

// ExecuteExternal 执行外部操作（发送命令到设备）
// 在这个模拟中，我们不直接改 Actor 状态，而是：
// - 启动一个 goroutine，模拟设备执行后再通过 ExternalEventMessage 回馈
func (b *SimulatedBreakerBinding) ExecuteExternal(ctx context.Context, command interface{}) error {
	go func() {
		// 模拟设备执行耗时
		time.Sleep(100 * time.Millisecond)

		var event *BreakerDeviceEvent

		switch cmd := command.(type) {
		case *OpenBreakerCommand:
			event = &BreakerDeviceEvent{
				Action:   "opened",
				Operator: cmd.Operator,
				Reason:   cmd.Reason,
			}
		case *CloseBreakerCommand:
			event = &BreakerDeviceEvent{
				Action:   "closed",
				Operator: cmd.Operator,
				Reason:   cmd.Reason,
			}
		default:
			return
		}

		_ = b.OnExternalEvent(context.Background(), event)
	}()

	return nil
}
