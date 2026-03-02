package main

import (
	"time"

	"github.com/uos-projects/uos-kernel/actor"
)

// ============================================================================
// 事件定义（Events）— 设备级事件
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
	return actor.MessageCategoryEvent
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
	return actor.MessageCategoryEvent
}

func (e *MaintenanceRequiredEvent) SourceActorID() string {
	return e.DeviceID
}
