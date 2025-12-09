package capacities

import (
	"context"
	"fmt"
)

// SetPointMessage 设定点消息
type SetPointMessage struct {
	Value float64
}

// SetPointCapacity 实现 SetPoint 控制能力
type SetPointCapacity struct {
	BaseCapacity
}

// NewSetPointCapacity 创建新的 SetPointCapacity
func NewSetPointCapacity(controlID string) *SetPointCapacity {
	return &SetPointCapacity{
		BaseCapacity: BaseCapacity{
			controlID: controlID,
			name:      "SetPointCapacity",
		},
	}
}

func (c *SetPointCapacity) CanHandle(msg Message) bool {
	_, ok := msg.(*SetPointMessage)
	return ok
}

func (c *SetPointCapacity) Execute(ctx context.Context, msg Message) error {
	setPointMsg, ok := msg.(*SetPointMessage)
	if !ok {
		return fmt.Errorf("invalid message type for SetPointCapacity")
	}

	// 执行设定点逻辑
	// TODO: 实现具体的业务逻辑
	_ = setPointMsg
	return nil
}
