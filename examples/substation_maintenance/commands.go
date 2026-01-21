package main

import (
	"fmt"
	"time"

	"github.com/uos-projects/uos-kernel/actors"
)

// ============================================================================
// 命令定义（Capability Commands）
// ============================================================================

// OpenBreakerCommand 打开断路器命令
type OpenBreakerCommand struct {
	commandID string // 私有字段
	Reason    string
	Operator  string
	TaskID    string // 关联的检修任务 ID
}

func (c *OpenBreakerCommand) MessageType() actors.MessageCategory {
	return actors.MessageCategoryCapabilityCommand
}

func (c *OpenBreakerCommand) CommandID() string {
	if c.commandID == "" {
		c.commandID = fmt.Sprintf("open_breaker_%d", time.Now().UnixNano())
	}
	return c.commandID
}

// CloseBreakerCommand 关闭断路器命令
type CloseBreakerCommand struct {
	commandID string // 私有字段
	Reason    string
	Operator  string
	TaskID    string
}

func (c *CloseBreakerCommand) MessageType() actors.MessageCategory {
	return actors.MessageCategoryCapabilityCommand
}

func (c *CloseBreakerCommand) CommandID() string {
	if c.commandID == "" {
		c.commandID = fmt.Sprintf("close_breaker_%d", time.Now().UnixNano())
	}
	return c.commandID
}

// StartMaintenanceCommand 开始检修命令
// 由调度中心发送给操作员，启动检修任务
type StartMaintenanceCommand struct {
	commandID   string // 私有字段
	TaskID      string
	Type        string // "scheduled", "emergency"
	Devices     []string
	Description string
	Reason      string
	OperatorID  string // 操作员ID
}

func (c *StartMaintenanceCommand) MessageType() actors.MessageCategory {
	return actors.MessageCategoryCapabilityCommand
}

func (c *StartMaintenanceCommand) CommandID() string {
	if c.commandID == "" {
		c.commandID = fmt.Sprintf("start_maintenance_%d", time.Now().UnixNano())
	}
	return c.commandID
}

// CompleteMaintenanceCommand 完成检修命令
type CompleteMaintenanceCommand struct {
	commandID string // 私有字段
	TaskID    string
	Operator  string
	Result    string
}

func (c *CompleteMaintenanceCommand) MessageType() actors.MessageCategory {
	return actors.MessageCategoryCapabilityCommand
}

func (c *CompleteMaintenanceCommand) CommandID() string {
	if c.commandID == "" {
		c.commandID = fmt.Sprintf("complete_maintenance_%d", time.Now().UnixNano())
	}
	return c.commandID
}
