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
}

// PropertyHolder 属性持有者接口（可选，CIMResourceActor 实现）
type PropertyHolder interface {
	// GetProperty 获取属性值
	GetProperty(name string) (interface{}, bool)

	// SetProperty 设置属性值
	SetProperty(name string, value interface{})

	// GetAllProperties 获取所有属性
	GetAllProperties() map[string]interface{}

	// GetOWLClassURI 获取 OWL 类 URI
	GetOWLClassURI() string
}
