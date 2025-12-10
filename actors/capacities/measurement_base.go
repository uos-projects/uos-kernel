package capacities

import (
	"context"
	"fmt"
	"sync"

	"github.com/uos-projects/uos-kernel/actors/mq"
)

// SubscriptionStarter 需要订阅的 Capacity 接口
type SubscriptionStarter interface {
	StartSubscription(ctx context.Context) error
}

// ActorRefSetter 需要设置 Actor 引用的 Capacity 接口
type ActorRefSetter interface {
	SetActorRef(actorRef interface {
		Send(msg Message) bool
	})
}

// BaseMeasurementCapacity 提供 Measurement Capacity 的基础实现
// 包含 MQ 订阅、Actor 引用等共同功能
type BaseMeasurementCapacity struct {
	BaseCapacity
	measurementID   string
	measurementType string

	// MQ 订阅相关
	mqConsumer mq.MQConsumer
	subscribed bool
	subChan    <-chan *mq.MQMessage
	mu         sync.RWMutex

	// Actor 引用（用于发送消息，可选）
	actorRef interface {
		Send(msg Message) bool
	}

	// Capacity 引用（用于直接调用 Execute，优先使用）
	capacityRef Capacity
	ctx         context.Context // 用于 Execute 调用
}

// NewBaseMeasurementCapacity 创建基础 Measurement Capacity
func NewBaseMeasurementCapacity(
	measurementID string,
	measurementType string,
	name string,
	mqConsumer mq.MQConsumer,
) *BaseMeasurementCapacity {
	return &BaseMeasurementCapacity{
		BaseCapacity: BaseCapacity{
			resourceID: measurementID,
			name:       name,
		},
		measurementID:   measurementID,
		measurementType: measurementType,
		mqConsumer:      mqConsumer,
		subscribed:      false,
	}
}

// SetActorRef 设置 Actor 引用（用于发送消息）
func (b *BaseMeasurementCapacity) SetActorRef(actorRef interface {
	Send(msg Message) bool
}) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.actorRef = actorRef
}

// SetCapacityRef 设置 Capacity 引用和 context（用于直接调用 Execute）
func (b *BaseMeasurementCapacity) SetCapacityRef(capacity Capacity, ctx context.Context) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.capacityRef = capacity
	b.ctx = ctx
}

// MeasurementID 返回测量 ID
func (b *BaseMeasurementCapacity) MeasurementID() string {
	return b.measurementID
}

// MeasurementType 返回测量类型
func (b *BaseMeasurementCapacity) MeasurementType() string {
	return b.measurementType
}

// StartSubscription 启动 MQ 订阅
func (b *BaseMeasurementCapacity) StartSubscription(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.subscribed {
		return nil // 已订阅
	}

	// 订阅 MQ topic
	topic := fmt.Sprintf("measurement/%s", b.measurementID)

	ch, err := b.mqConsumer.Subscribe(ctx, topic)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}

	b.subChan = ch
	b.subscribed = true

	return nil
}

// StopSubscription 停止 MQ 订阅
func (b *BaseMeasurementCapacity) StopSubscription() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribed = false
}

// IsSubscribed 检查是否已订阅
func (b *BaseMeasurementCapacity) IsSubscribed() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.subscribed
}

// GetSubChan 获取订阅 channel（用于消费循环）
func (b *BaseMeasurementCapacity) GetSubChan() <-chan *mq.MQMessage {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.subChan
}

// GetActorRef 获取 Actor 引用
func (b *BaseMeasurementCapacity) GetActorRef() interface {
	Send(msg Message) bool
} {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.actorRef
}

// GetMQConsumer 获取 MQ 消费者
func (b *BaseMeasurementCapacity) GetMQConsumer() mq.MQConsumer {
	return b.mqConsumer
}

// GetMeasurementID 获取测量 ID（内部使用）
func (b *BaseMeasurementCapacity) GetMeasurementID() string {
	return b.measurementID
}

// StartConsumeLoop 启动消费循环（通用实现）
// messageConverter 用于将 MQ 消息转换为具体的消息类型
// 优先直接调用 Capacity.Execute()，如果没有设置 Capacity 引用则通过 Actor mailbox
func (b *BaseMeasurementCapacity) StartConsumeLoop(
	ctx context.Context,
	topic string,
	messageConverter func(*mq.MQMessage) Message,
) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case mqMsg, ok := <-b.GetSubChan():
				if !ok {
					// channel 已关闭
					return
				}

				// 转换为内部消息格式
				valueMsg := messageConverter(mqMsg)

				// 优先直接调用 Capacity.Execute()（更高效，避免额外的消息传递）
				b.mu.RLock()
				capacityRef := b.capacityRef
				execCtx := b.ctx
				b.mu.RUnlock()

				if capacityRef != nil && execCtx != nil {
					// 直接调用 Execute，避免通过 Actor mailbox
					if err := capacityRef.Execute(execCtx, valueMsg); err != nil {
						// TODO: 添加错误日志
					}
				} else if actorRef := b.GetActorRef(); actorRef != nil {
					// 回退到通过 Actor mailbox（向后兼容）
					if !actorRef.Send(valueMsg) {
						// 邮箱满了，可以记录日志或采取其他措施
						// TODO: 添加日志
					}
				}
			}
		}
	}()
}
