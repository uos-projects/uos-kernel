package resource

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
	OWLClassURI  string                 // OWL 类 URI（如果是 CIMResourceActor）
	Properties   map[string]interface{} // 属性（如果是 CIMResourceActor）
}

// Read 读取 Actor 状态
// 类似 POSIX read()，但这里读取的是 Actor 的状态信息
func (rm *ResourceManager) Read(ctx context.Context, fd ResourceDescriptor) (*ActorState, error) {
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

	// 如果是 PropertyHolder（CIMResourceActor），获取属性
	if holder, ok := actor.(actors.PropertyHolder); ok {
		state.OWLClassURI = holder.GetOWLClassURI()
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
func (rm *ResourceManager) Write(ctx context.Context, fd ResourceDescriptor, req *WriteRequest) error {
	handle, err := rm.GetHandle(fd)
	if err != nil {
		return err
	}

	// 通过Resource访问Actor
	resource := handle.resource
	resource.mu.RLock()
	actor := resource.actor
	resource.mu.RUnlock()

	// 检查是否支持属性更新
	holder, ok := actor.(actors.PropertyHolder)
	if !ok {
		// 如果 Actor 不支持属性更新，这里怎么处理？
		// 可以返回错误，表示该资源不支持写入属性
		return fmt.Errorf("resource %s does not support property updates", actor.ResourceID())
	}

	// 更新属性
	if req.Updates != nil {
		for key, value := range req.Updates {
			holder.SetProperty(key, value)
		}
	}

	return nil
}
