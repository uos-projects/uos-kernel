package main

import (
	"context"

	"github.com/uos-projects/uos-kernel/actor"
)

// BreakerSwitchingCapacity 断路器开关控制能力
// 将断路器的开/合操作封装为一个 Capacity，由 BaseResourceActor 进行路由
type BreakerSwitchingCapacity struct {
	*actor.BaseCapacity

	breaker *BreakerActor
}

// NewBreakerSwitchingCapacity 创建断路器开关控制能力
func NewBreakerSwitchingCapacity(breaker *BreakerActor) *BreakerSwitchingCapacity {
	return &BreakerSwitchingCapacity{
		BaseCapacity: actor.NewBaseCapacity("BreakerSwitchingCapacity", breaker.ResourceID()),
		breaker:      breaker,
	}
}

// CanHandle 检查是否能处理消息
// 这里只处理 OpenBreakerCommand 和 CloseBreakerCommand
func (c *BreakerSwitchingCapacity) CanHandle(msg actor.Message) bool {
	switch msg.(type) {
	case *OpenBreakerCommand, *CloseBreakerCommand:
		return true
	default:
		return false
	}
}

// Execute 执行断路器开/合操作
// 注意：这里只负责"发起控制"，不直接修改 Actor 状态
// 真正的状态变化由设备反馈（通过 Binding）驱动
func (c *BreakerSwitchingCapacity) Execute(ctx context.Context, msg actor.Message) error {
	binding, ok := c.breaker.GetBinding(actor.BindingTypeDevice)
	if !ok {
		return nil
	}

	switch cmd := msg.(type) {
	case *OpenBreakerCommand, *CloseBreakerCommand:
		return binding.ExecuteExternal(ctx, cmd)
	default:
		return nil
	}
}
