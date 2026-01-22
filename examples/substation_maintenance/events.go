package main

import (
	"time"

	"github.com/uos-projects/uos-kernel/actor"
)

// ============================================================================
// 事件定义（Coordination Events）
// ============================================================================

// DeviceAbnormalEvent 设备异常事件
// 当设备检测到异常状态时发射此事件
type DeviceAbnormalEvent struct {
	DeviceID  string
	EventType string // "temperature_high", "voltage_abnormal", "current_overload"
	Severity  string // "warning", "critical"
	Details   map[string]interface{}
	Timestamp time.Time
}

func (e *DeviceAbnormalEvent) MessageType() actor.MessageCategory {
	return actor.MessageCategoryCoordinationEvent
}

func (e *DeviceAbnormalEvent) SourceActorID() string {
	return e.DeviceID
}

// MaintenanceRequiredEvent 需要检修事件
// 当设备达到检修时间或检测到需要检修时发射此事件
type MaintenanceRequiredEvent struct {
	DeviceID            string
	Reason              string // "scheduled", "overdue", "abnormal"
	LastMaintenanceTime time.Time
	OperationHours      int64
	Details             map[string]interface{}
	Timestamp           time.Time
}

func (e *MaintenanceRequiredEvent) MessageType() actor.MessageCategory {
	return actor.MessageCategoryCoordinationEvent
}

func (e *MaintenanceRequiredEvent) SourceActorID() string {
	return e.DeviceID
}

// MaintenanceTaskCreatedEvent 检修任务创建事件
// 调度中心创建检修任务时发射此事件
type MaintenanceTaskCreatedEvent struct {
	TaskID      string
	Type        string // "scheduled", "emergency"
	Devices     []string
	Description string
	Reason      string
	Timestamp   time.Time
}

func (e *MaintenanceTaskCreatedEvent) MessageType() actor.MessageCategory {
	return actor.MessageCategoryCoordinationEvent
}

func (e *MaintenanceTaskCreatedEvent) SourceActorID() string {
	return "DISPATCHER" // 调度中心 ID
}

// MaintenanceTaskAssignedEvent 检修任务分配事件
// 调度中心将检修任务分配给操作员时发射此事件
type MaintenanceTaskAssignedEvent struct {
	TaskID     string
	OperatorID string
	DeviceIDs  []string
	Reason     string
	Timestamp  time.Time
}

func (e *MaintenanceTaskAssignedEvent) MessageType() actor.MessageCategory {
	return actor.MessageCategoryCoordinationEvent
}

func (e *MaintenanceTaskAssignedEvent) SourceActorID() string {
	return "DISPATCHER" // 调度中心 ID
}

// MaintenanceTaskUpdatedEvent 检修任务更新事件
// 任务状态发生变化时发射此事件
type MaintenanceTaskUpdatedEvent struct {
	TaskID    string
	Status    string // "pending", "assigned", "in_progress", "completed", "failed"
	Timestamp time.Time
}

func (e *MaintenanceTaskUpdatedEvent) MessageType() actor.MessageCategory {
	return actor.MessageCategoryCoordinationEvent
}

func (e *MaintenanceTaskUpdatedEvent) SourceActorID() string {
	return "DISPATCHER" // 调度中心 ID
}

// MaintenanceCompletedEvent 检修完成事件
// 操作员完成检修操作后发射此事件
type MaintenanceCompletedEvent struct {
	TaskID     string
	OperatorID string
	DeviceIDs  []string
	Result     string // "success", "failed"
	Timestamp  time.Time
}

func (e *MaintenanceCompletedEvent) MessageType() actor.MessageCategory {
	return actor.MessageCategoryCoordinationEvent
}

func (e *MaintenanceCompletedEvent) SourceActorID() string {
	return e.OperatorID
}
