package capacities

import (
	"context"
	"fmt"
)

// AccumulatorResetMessage 累加器复位消息
type AccumulatorResetMessage struct {
	Value int
}

// AccumulatorResetCapacity 实现 AccumulatorReset 控制能力
type AccumulatorResetCapacity struct {
	BaseCapacity
}

// NewAccumulatorResetCapacity 创建新的 AccumulatorResetCapacity
func NewAccumulatorResetCapacity(controlID string) *AccumulatorResetCapacity {
	return &AccumulatorResetCapacity{
		BaseCapacity: BaseCapacity{
			resourceID: controlID,
			name:       "AccumulatorResetCapacity",
		},
	}
}

func (c *AccumulatorResetCapacity) CanHandle(msg Message) bool {
	_, ok := msg.(*AccumulatorResetMessage)
	return ok
}

func (c *AccumulatorResetCapacity) Execute(ctx context.Context, msg Message) error {
	resetMsg, ok := msg.(*AccumulatorResetMessage)
	if !ok {
		return fmt.Errorf("invalid message type for AccumulatorResetCapacity")
	}

	// 执行累加器复位逻辑
	// TODO: 实现具体的业务逻辑
	_ = resetMsg
	return nil
}
