package actor

import (
	"reflect"
)

// ============================================================================
// 事件描述符（Event Descriptor）
// ============================================================================

// EventDescriptor 事件描述符
// 描述 Actor 能发出的某种事件类型
type EventDescriptor struct {
	// Name 事件名称（唯一标识，如 "DeviceAbnormalEvent"）
	Name string

	// PayloadType Payload 的类型（反射类型）
	PayloadType reflect.Type

	// Description 事件描述
	Description string

	// ResourceID 关联的资源 ID
	ResourceID string
}

// NewEventDescriptor 创建事件描述符
func NewEventDescriptor(
	name string,
	payloadType reflect.Type,
	description string,
	resourceID string,
) *EventDescriptor {
	return &EventDescriptor{
		Name:        name,
		PayloadType: payloadType,
		Description: description,
		ResourceID:  resourceID,
	}
}

// CanEmit 检查是否能发出指定类型的事件
func (ed *EventDescriptor) CanEmit(payload interface{}) bool {
	if payload == nil {
		return true
	}
	payloadType := reflect.TypeOf(payload)
	return ed.PayloadType == payloadType || payloadType.AssignableTo(ed.PayloadType)
}

// ============================================================================
// 事件注册接口（可选实现）
// ============================================================================

// EventRegistry 事件注册接口
// Actor 可以实现此接口来声明自己能发出的事件
type EventRegistry interface {
	// DeclaredEvents 返回 Actor 声明能发出的事件描述符列表
	DeclaredEvents() []*EventDescriptor
}
