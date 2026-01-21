package main

import (
	"time"

	"github.com/uos-projects/uos-kernel/actors"
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

func (e *DeviceAbnormalEvent) MessageType() actors.MessageCategory {
	return actors.MessageCategoryCoordinationEvent
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

func (e *MaintenanceRequiredEvent) MessageType() actors.MessageCategory {
	return actors.MessageCategoryCoordinationEvent
}

func (e *MaintenanceRequiredEvent) SourceActorID() string {
	return e.DeviceID
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

func (e *MaintenanceTaskAssignedEvent) MessageType() actors.MessageCategory {
	return actors.MessageCategoryCoordinationEvent
}

func (e *MaintenanceTaskAssignedEvent) SourceActorID() string {
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

func (e *MaintenanceCompletedEvent) MessageType() actors.MessageCategory {
	return actors.MessageCategoryCoordinationEvent
}

func (e *MaintenanceCompletedEvent) SourceActorID() string {
	return e.OperatorID
}
