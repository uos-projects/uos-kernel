package main

import (
	"fmt"
	"time"

	"github.com/uos-projects/uos-kernel/actor"
)

// ============================================================================
// 命令定义（Commands）— 设备级命令
// ============================================================================

// OpenBreakerCommand 打开断路器命令
type OpenBreakerCommand struct {
	commandID string
	Reason    string
	Operator  string
}

func (c *OpenBreakerCommand) MessageType() actor.MessageCategory {
	return actor.MessageCategoryCommand
}

func (c *OpenBreakerCommand) CommandID() string {
	if c.commandID == "" {
		c.commandID = fmt.Sprintf("open_breaker_%d", time.Now().UnixNano())
	}
	return c.commandID
}

// CloseBreakerCommand 关闭断路器命令
type CloseBreakerCommand struct {
	commandID string
	Reason    string
	Operator  string
}

func (c *CloseBreakerCommand) MessageType() actor.MessageCategory {
	return actor.MessageCategoryCommand
}

func (c *CloseBreakerCommand) CommandID() string {
	if c.commandID == "" {
		c.commandID = fmt.Sprintf("close_breaker_%d", time.Now().UnixNano())
	}
	return c.commandID
}

// CompleteMaintenanceCommand 完成检修命令
type CompleteMaintenanceCommand struct {
	commandID string
	Operator  string
	Result    string
}

func (c *CompleteMaintenanceCommand) MessageType() actor.MessageCategory {
	return actor.MessageCategoryCommand
}

func (c *CompleteMaintenanceCommand) CommandID() string {
	if c.commandID == "" {
		c.commandID = fmt.Sprintf("complete_maintenance_%d", time.Now().UnixNano())
	}
	return c.commandID
}
