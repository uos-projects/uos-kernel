package capacities

import (
	"context"
	"fmt"

	"github.com/uos-projects/uos-kernel/actors/capacities"
)

// AccumulatorResetMessage 累加器复位消息（CIM Control 类型）
type AccumulatorResetMessage struct {
	Value int
}

// AccumulatorResetCapacity 实现 AccumulatorReset 控制能力（CIM 特定）
type AccumulatorResetCapacity struct {
	*capacities.BaseCapacity
}

// NewAccumulatorResetCapacity 创建新的 AccumulatorResetCapacity
func NewAccumulatorResetCapacity(controlID string) *AccumulatorResetCapacity {
	return &AccumulatorResetCapacity{
		BaseCapacity: capacities.NewBaseCapacity("AccumulatorResetCapacity", controlID),
	}
}

func (c *AccumulatorResetCapacity) CanHandle(msg capacities.Message) bool {
	_, ok := msg.(*AccumulatorResetMessage)
	return ok
}

func (c *AccumulatorResetCapacity) Execute(ctx context.Context, msg capacities.Message) error {
	resetMsg, ok := msg.(*AccumulatorResetMessage)
	if !ok {
		return fmt.Errorf("invalid message type for AccumulatorResetCapacity")
	}

	// 执行累加器复位逻辑
	// TODO: 实现具体的业务逻辑
	_ = resetMsg
	return nil
}
