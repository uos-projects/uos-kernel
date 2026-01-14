package capacities

import (
	"context"
	"fmt"

	"github.com/uos-projects/uos-kernel/actors/capacities"
)

// SetPointMessage 设定点消息（CIM Control 类型）
type SetPointMessage struct {
	Value float64
}

// SetPointCapacity 实现 SetPoint 控制能力（CIM 特定）
type SetPointCapacity struct {
	*capacities.BaseCapacity
}

// NewSetPointCapacity 创建新的 SetPointCapacity
func NewSetPointCapacity(controlID string) *SetPointCapacity {
	return &SetPointCapacity{
		BaseCapacity: capacities.NewBaseCapacity("SetPointCapacity", controlID),
	}
}

func (c *SetPointCapacity) CanHandle(msg capacities.Message) bool {
	_, ok := msg.(*SetPointMessage)
	return ok
}

func (c *SetPointCapacity) Execute(ctx context.Context, msg capacities.Message) error {
	setPointMsg, ok := msg.(*SetPointMessage)
	if !ok {
		return fmt.Errorf("invalid message type for SetPointCapacity")
	}

	// 执行设定点逻辑
	// TODO: 实现具体的业务逻辑
	_ = setPointMsg
	return nil
}
