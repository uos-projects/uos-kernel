package actors

import (
	"context"
	"fmt"
)

// BindingType 绑定类型
type BindingType string

const (
	BindingTypeDevice  BindingType = "device"  // 设备协议
	BindingTypeHuman  BindingType = "human"   // 人机交互
	BindingTypeService BindingType = "service" // 服务 API
	BindingTypeMQ     BindingType = "mq"       // 消息队列
)

// Binding 外部绑定接口
// 用于与真实世界交互的适配层，可插拔（设备协议/人机交互/服务API）
type Binding interface {
	// Type 返回绑定类型
	Type() BindingType

	// ResourceID 返回关联的资源ID
	ResourceID() string

	// Start 启动绑定
	Start(ctx context.Context) error

	// Stop 停止绑定
	Stop() error

	// OnExternalEvent 接收外部事件（转换为消息发送到 Actor）
	// 外部事件到达时，Binding 将其转换为消息并发送给关联的 Actor
	OnExternalEvent(ctx context.Context, event interface{}) error

	// ExecuteExternal 执行外部操作（从 Actor 消息转换为外部协议）
	// Actor 需要执行外部操作时，通过此方法转换为外部协议并执行
	ExecuteExternal(ctx context.Context, command interface{}) error
}

// BaseBinding 基础绑定实现
// 提供通用的绑定功能，具体绑定类型可以嵌入此结构
// 注意：BaseBinding 本身不实现 Binding 接口，需要具体实现类来实现
type BaseBinding struct {
	bindingType BindingType
	resourceID  string
	actorRef    interface {
		Send(Message) bool
	}
	ctx    context.Context
	cancel context.CancelFunc
}

// NewBaseBinding 创建基础绑定
func NewBaseBinding(bindingType BindingType, resourceID string) *BaseBinding {
	ctx, cancel := context.WithCancel(context.Background())
	return &BaseBinding{
		bindingType: bindingType,
		resourceID:  resourceID,
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Type 返回绑定类型
func (b *BaseBinding) Type() BindingType {
	return b.bindingType
}

// ResourceID 返回资源ID
func (b *BaseBinding) ResourceID() string {
	return b.resourceID
}

// SetActorRef 设置 Actor 引用（由 Actor 在添加绑定时调用）
func (b *BaseBinding) SetActorRef(actorRef interface {
	Send(Message) bool
}) {
	b.actorRef = actorRef
}

// GetContext 获取绑定的上下文
func (b *BaseBinding) GetContext() context.Context {
	return b.ctx
}

// Stop 停止绑定（基础实现，子类可以重写）
func (b *BaseBinding) Stop() error {
	if b.cancel != nil {
		b.cancel()
	}
	return nil
}

// SendToActor 发送消息到 Actor（辅助方法）
func (b *BaseBinding) SendToActor(msg Message) error {
	if b.actorRef == nil {
		return fmt.Errorf("actor reference not set for binding %s", b.bindingType)
	}
	if !b.actorRef.Send(msg) {
		return fmt.Errorf("failed to send message to actor: mailbox full")
	}
	return nil
}

// ActorRefSetter 接口，用于设置 Actor 引用
type ActorRefSetter interface {
	SetActorRef(actorRef interface {
		Send(Message) bool
	})
}
