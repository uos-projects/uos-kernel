package capacities

import (
	"context"
	"fmt"

	"github.com/uos-projects/uos-kernel/actors/capacities"
)

// RaiseLowerCommandMessage 升降命令消息（CIM Control 类型）
type RaiseLowerCommandMessage struct {
	Delta float64
}

// RaiseLowerCommandCapacity 实现 RaiseLowerCommand 控制能力（CIM 特定）
type RaiseLowerCommandCapacity struct {
	*capacities.BaseCapacity
}

// NewRaiseLowerCommandCapacity 创建新的 RaiseLowerCommandCapacity
func NewRaiseLowerCommandCapacity(controlID string) *RaiseLowerCommandCapacity {
	return &RaiseLowerCommandCapacity{
		BaseCapacity: capacities.NewBaseCapacity("RaiseLowerCommandCapacity", controlID),
	}
}

func (c *RaiseLowerCommandCapacity) CanHandle(msg capacities.Message) bool {
	_, ok := msg.(*RaiseLowerCommandMessage)
	return ok
}

func (c *RaiseLowerCommandCapacity) Execute(ctx context.Context, msg capacities.Message) error {
	cmdMsg, ok := msg.(*RaiseLowerCommandMessage)
	if !ok {
		return fmt.Errorf("invalid message type for RaiseLowerCommandCapacity")
	}

	// 执行升降命令逻辑
	// TODO: 实现具体的业务逻辑
	_ = cmdMsg
	return nil
}
