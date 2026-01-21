package actors

import (
	"github.com/uos-projects/uos-kernel/actors/capacities"
)

// ResourceActor 资源 Actor 接口
// 定义资源 Actor 的公共方法，供 ResourceManager 使用
type ResourceActor interface {
	Actor // 继承基础 Actor 接口

	// ResourceID 返回资源 ID
	ResourceID() string

	// ResourceType 返回资源类型
	ResourceType() string

	// ListCapabilities 返回所有能力名称列表
	ListCapabilities() []string

	// HasCapacity 检查是否具有某种能力
	HasCapacity(capacityName string) bool

	// GetCapacity 获取指定名称的能力
	GetCapacity(capacityName string) (capacities.Capacity, bool)

	// AddCapacity 添加一个能力
	AddCapacity(capacity capacities.Capacity)

	// Send 发送消息
	Send(msg Message) bool

	// ============================================================================
	// 事件管理方法（参考 Capacity 管理）
	// ============================================================================

	// ListEvents 返回所有事件名称列表
	ListEvents() []string

	// HasEvent 检查是否能发出某种事件
	HasEvent(eventName string) bool

	// GetEvent 获取指定名称的事件描述符
	GetEvent(eventName string) (*EventDescriptor, bool)

	// ListEventDescriptors 返回所有事件描述符列表
	ListEventDescriptors() []*EventDescriptor

	// CanEmitEvent 检查是否能发出指定类型的事件
	CanEmitEvent(eventType EventType, payload interface{}) bool
}

// PropertyHolder 属性持有者接口（业务中立）
// 提供属性访问功能，不绑定任何特定领域模型
type PropertyHolder interface {
	// GetProperty 获取属性值
	GetProperty(name string) (interface{}, bool)

	// GetAllProperties 获取所有属性
	GetAllProperties() map[string]interface{}
}
