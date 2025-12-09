package actors

import (
	"context"
	"fmt"

	"github.com/uos-projects/uos-kernel/actors/capacities"
)

// PowerSystemResourceActor 代表一个电力系统资源
type PowerSystemResourceActor struct {
	*BaseActor
	resourceID   string
	resourceType string
	capabilities map[string]capacities.Capacity // 能力名称 -> Capacity 实现
}

// GetRef 获取指定 Actor 的 ActorRef
func (a *PowerSystemResourceActor) GetRef(targetID string) (*ActorRef, error) {
	return a.BaseActor.GetRef(targetID)
}

// SendTo 向指定 ID 的 Actor 发送消息
func (a *PowerSystemResourceActor) SendTo(targetID string, msg Message) error {
	return a.BaseActor.SendTo(targetID, msg)
}

// Tell 是 SendTo 的别名
func (a *PowerSystemResourceActor) Tell(targetID string, msg Message) error {
	return a.BaseActor.Tell(targetID, msg)
}

// ID 实现 Actor 接口，返回资源 ID
func (a *PowerSystemResourceActor) ID() string {
	return a.resourceID
}

// NewPowerSystemResourceActor 创建一个新的资源 Actor
func NewPowerSystemResourceActor(
	id string,
	resourceType string,
	behavior ActorBehavior,
) *PowerSystemResourceActor {
	baseActor := NewBaseActor(id, behavior)
	return &PowerSystemResourceActor{
		BaseActor:    baseActor,
		resourceID:   id,
		resourceType: resourceType,
		capabilities: make(map[string]capacities.Capacity),
	}
}

// ResourceID 返回资源 ID
func (a *PowerSystemResourceActor) ResourceID() string {
	return a.resourceID
}

// ResourceType 返回资源类型
func (a *PowerSystemResourceActor) ResourceType() string {
	return a.resourceType
}

// AddCapacity 添加一个能力
func (a *PowerSystemResourceActor) AddCapacity(capacity capacities.Capacity) {
	a.capabilities[capacity.Name()] = capacity
}

// RemoveCapacity 移除一个能力
func (a *PowerSystemResourceActor) RemoveCapacity(capacityName string) {
	delete(a.capabilities, capacityName)
}

// HasCapacity 检查是否具有某种能力
func (a *PowerSystemResourceActor) HasCapacity(capacityName string) bool {
	_, exists := a.capabilities[capacityName]
	return exists
}

// GetCapacity 获取指定名称的能力
func (a *PowerSystemResourceActor) GetCapacity(capacityName string) (capacities.Capacity, bool) {
	cap, exists := a.capabilities[capacityName]
	return cap, exists
}

// ListCapabilities 返回所有能力名称列表
func (a *PowerSystemResourceActor) ListCapabilities() []string {
	names := make([]string, 0, len(a.capabilities))
	for name := range a.capabilities {
		names = append(names, name)
	}
	return names
}

// Receive 重写消息接收逻辑，路由到对应的 Capacity
func (a *PowerSystemResourceActor) Receive(ctx context.Context, msg Message) error {
	// 尝试找到能处理此消息的 Capacity
	// 将 actors.Message 转换为 capacities.Message
	var capMsg capacities.Message = msg
	for _, capacity := range a.capabilities {
		if capacity.CanHandle(capMsg) {
			return capacity.Execute(ctx, capMsg)
		}
	}

	// 如果没有找到对应的 Capacity，使用默认行为
	if a.BaseActor.behavior != nil {
		_, err := a.BaseActor.behavior.Handle(ctx, msg)
		return err
	}

	// 如果也没有默认行为，返回错误
	return fmt.Errorf("no capacity can handle message type %T", msg)
}

// Start 启动 Actor（扩展基类方法，启动所有需要订阅的 Capacity）
func (a *PowerSystemResourceActor) Start(ctx context.Context) error {
	// 先调用基类的 Start
	if err := a.BaseActor.Start(ctx); err != nil {
		return err
	}

	// 启动所有需要订阅的 Capacity
	for _, capacity := range a.capabilities {
		// 检查是否有 StartSubscription 方法（可选接口）
		if starter, ok := capacity.(interface {
			StartSubscription(context.Context) error
		}); ok {
			// 设置 Actor 引用（如果 Capacity 支持）
			if setter, ok := capacity.(interface {
				SetActorRef(interface {
					Send(Message) bool
				})
			}); ok {
				setter.SetActorRef(a)
			}

			// 启动订阅
			if err := starter.StartSubscription(ctx); err != nil {
				return fmt.Errorf("failed to start subscription for capacity %s: %w",
					capacity.Name(), err)
			}
		}
	}

	return nil
}
