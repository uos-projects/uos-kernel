package kernel

import (
	"context"
	"fmt"
)

// ActorState Actor 状态
type ActorState struct {
	ResourceID   string
	ResourceType string
	Capabilities []string
	Properties map[string]interface{} // 属性（如果是 PropertyHolder）
}

// Read 读取 Actor 状态
// 类似 POSIX read()，但这里读取的是 Actor 的状态信息
func (rm *Manager) Read(ctx context.Context, fd ResourceDescriptor) (*ActorState, error) {
	handle, err := rm.GetHandle(fd)
	if err != nil {
		return nil, err
	}

	// 通过Resource访问Actor
	resource := handle.resource
	resource.mu.RLock()
	a := resource.actor
	resource.mu.RUnlock()

	// 构建基础状态
	state := &ActorState{
		ResourceID:   a.ResourceID(),
		ResourceType: a.ResourceType(),
		Capabilities: a.ListCapabilities(),
	}

	// 如果是 propertyHolder，获取属性
	if holder, ok := a.(propertyHolder); ok {
		state.Properties = holder.GetAllProperties()
	}

	return state, nil
}

// WriteRequest 写入请求（改变 Actor 状态）
type WriteRequest struct {
	// 可以包含状态更新信息
	// 例如：更新资源属性、配置等
	Updates map[string]interface{}
}

// propertyHolder 属性读取接口
type propertyHolder interface {
	GetAllProperties() map[string]interface{}
}

// propertySetter 属性设置接口
type propertySetter interface {
	SetProperty(name string, value interface{})
}

// Write 改变 Actor 状态
// 类似 POSIX write()，但这里改变的是 Actor 的状态
func (rm *Manager) Write(ctx context.Context, fd ResourceDescriptor, req *WriteRequest) error {
	handle, err := rm.GetHandle(fd)
	if err != nil {
		return err
	}

	// 通过Resource访问Actor
	resource := handle.resource
	resource.mu.RLock()
	a := resource.actor
	resource.mu.RUnlock()

	// 检查是否支持属性设置
	setter, ok := a.(propertySetter)
	if !ok {
		return fmt.Errorf("resource %s does not support property updates", a.ResourceID())
	}

	// 更新属性（直接方法调用，BaseResourceActor.SetProperty 内部有锁保护）
	if req.Updates != nil {
		for key, value := range req.Updates {
			setter.SetProperty(key, value)
		}
	}

	return nil
}
