package capacities

import (
	"context"
	"fmt"

	"github.com/uos-projects/uos-kernel/actors/capacities"
)

// CommandMessage 命令消息（CIM Control 类型）
type CommandMessage struct {
	Command string
	Value   int
}

// CommandCapacity 实现 Command 控制能力（CIM 特定）
type CommandCapacity struct {
	*capacities.BaseCapacity
}

// NewCommandCapacity 创建新的 CommandCapacity
func NewCommandCapacity(controlID string) *CommandCapacity {
	return &CommandCapacity{
		BaseCapacity: capacities.NewBaseCapacity("CommandCapacity", controlID),
	}
}

func (c *CommandCapacity) CanHandle(msg capacities.Message) bool {
	_, ok := msg.(*CommandMessage)
	return ok
}

func (c *CommandCapacity) Execute(ctx context.Context, msg capacities.Message) error {
	cmdMsg, ok := msg.(*CommandMessage)
	if !ok {
		return fmt.Errorf("invalid message type for CommandCapacity")
	}

	// 执行命令逻辑
	// TODO: 实现具体的业务逻辑
	_ = cmdMsg
	return nil
}
