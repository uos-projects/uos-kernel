package capacities

import (
	"context"
	"fmt"
)

// CommandMessage 命令消息
type CommandMessage struct {
	Command string
	Value   int
}

// CommandCapacity 实现 Command 控制能力
type CommandCapacity struct {
	BaseCapacity
}

// NewCommandCapacity 创建新的 CommandCapacity
func NewCommandCapacity(controlID string) *CommandCapacity {
	return &CommandCapacity{
		BaseCapacity: BaseCapacity{
			resourceID: controlID,
			name:       "CommandCapacity",
		},
	}
}

func (c *CommandCapacity) CanHandle(msg Message) bool {
	_, ok := msg.(*CommandMessage)
	return ok
}

func (c *CommandCapacity) Execute(ctx context.Context, msg Message) error {
	cmdMsg, ok := msg.(*CommandMessage)
	if !ok {
		return fmt.Errorf("invalid message type for CommandCapacity")
	}

	// 执行命令逻辑
	// TODO: 实现具体的业务逻辑
	_ = cmdMsg
	return nil
}
