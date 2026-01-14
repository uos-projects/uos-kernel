package kernel

import (
	"context"
	"fmt"

	"github.com/uos-projects/uos-kernel/actors"
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
	actor := resource.actor
	resource.mu.RUnlock()

	// 构建基础状态
	state := &ActorState{
		ResourceID:   actor.ResourceID(),
		ResourceType: actor.ResourceType(),
		Capabilities: actor.ListCapabilities(),
	}

	// 如果是 PropertyHolder，获取属性
	if holder, ok := actor.(actors.PropertyHolder); ok {
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
	actor := resource.actor
	resource.mu.RUnlock()

	// 检查是否支持属性更新（通过 PropertyHolder 接口）
	if _, ok := actor.(actors.PropertyHolder); !ok {
		// 如果 Actor 不支持属性更新，返回错误
		return fmt.Errorf("resource %s does not support property updates", actor.ResourceID())
	}

	// 更新属性（通过消息驱动）
	if req.Updates != nil {
		for key, value := range req.Updates {
			// 通过消息设置属性，符合 Actor 设计理念
			setPropMsg := &actors.SetPropertyMessage{
				Name:  key,
				Value: value,
			}
			if !actor.Send(setPropMsg) {
				return fmt.Errorf("failed to send SetProperty message for key %s", key)
			}
		}
	}

	return nil
}
