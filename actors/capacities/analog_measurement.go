package capacities

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/uos-projects/uos-kernel/actors/mq"
)

// AnalogMeasurementValueMessage 从 MQ 接收的模拟量测量值消息
type AnalogMeasurementValueMessage struct {
	MeasurementID string
	Value         float64
	Timestamp     time.Time
	Quality       mq.QualityCode
	Source        string
}

// AnalogMeasurementValue 模拟量测量值
type AnalogMeasurementValue struct {
	Value     float64
	Timestamp time.Time
	Quality   mq.QualityCode
	Source    string
}

// AnalogMeasurementCapacity 模拟量测量订阅能力
type AnalogMeasurementCapacity struct {
	BaseCapacity
	measurementID   string
	measurementType string // "ThreePhaseActivePower" 等

	// MQ 订阅相关
	mqConsumer mq.MQConsumer
	subscribed bool
	subChan    <-chan *mq.MQMessage

	// 当前值缓存
	currentValue *AnalogMeasurementValue
	mu           sync.RWMutex

	// Actor 引用（用于发送消息）
	actorRef interface {
		Send(msg Message) bool
	}
}

// NewAnalogMeasurementCapacity 创建新的模拟量测量能力
func NewAnalogMeasurementCapacity(
	measurementID string,
	measurementType string,
	mqConsumer mq.MQConsumer,
) *AnalogMeasurementCapacity {
	return &AnalogMeasurementCapacity{
		BaseCapacity: BaseCapacity{
			resourceID: measurementID,
			name:       "AnalogMeasurementCapacity",
		},
		measurementID:   measurementID,
		measurementType: measurementType,
		mqConsumer:      mqConsumer,
		subscribed:      false,
	}
}

// SetActorRef 设置 Actor 引用（用于发送消息）
func (c *AnalogMeasurementCapacity) SetActorRef(actorRef interface {
	Send(msg Message) bool
}) {
	c.actorRef = actorRef
}

// MeasurementID 返回测量 ID
func (c *AnalogMeasurementCapacity) MeasurementID() string {
	return c.measurementID
}

// MeasurementType 返回测量类型
func (c *AnalogMeasurementCapacity) MeasurementType() string {
	return c.measurementType
}

// CanHandle 检查是否能处理消息
func (c *AnalogMeasurementCapacity) CanHandle(msg Message) bool {
	_, ok := msg.(*AnalogMeasurementValueMessage)
	return ok
}

// Execute 处理测量值消息（从 MQ 接收）
func (c *AnalogMeasurementCapacity) Execute(ctx context.Context, msg Message) error {
	valueMsg, ok := msg.(*AnalogMeasurementValueMessage)
	if !ok {
		return fmt.Errorf("invalid message type for AnalogMeasurementCapacity")
	}

	// 更新当前值
	c.mu.Lock()
	c.currentValue = &AnalogMeasurementValue{
		Value:     valueMsg.Value,
		Timestamp: valueMsg.Timestamp,
		Quality:   valueMsg.Quality,
		Source:    valueMsg.Source,
	}
	c.mu.Unlock()

	// 可以触发业务逻辑（如阈值检查、告警等）
	// TODO: 实现业务逻辑

	return nil
}

// StartSubscription 启动 MQ 订阅
func (c *AnalogMeasurementCapacity) StartSubscription(ctx context.Context) error {
	if c.subscribed {
		return nil // 已订阅
	}

	// 订阅 MQ topic
	topic := fmt.Sprintf("measurement/%s", c.measurementID)

	ch, err := c.mqConsumer.Subscribe(ctx, topic)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}

	c.subChan = ch
	c.subscribed = true

	// 启动消费循环
	go c.consumeLoop(ctx, topic)

	return nil
}

// StopSubscription 停止 MQ 订阅
func (c *AnalogMeasurementCapacity) StopSubscription() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.subscribed = false
}

// consumeLoop MQ 消费循环
func (c *AnalogMeasurementCapacity) consumeLoop(ctx context.Context, topic string) {
	for {
		select {
		case <-ctx.Done():
			return
		case mqMsg, ok := <-c.subChan:
			if !ok {
				// channel 已关闭
				return
			}

			// 转换为内部消息格式
			valueMsg := &AnalogMeasurementValueMessage{
				MeasurementID: c.measurementID,
				Value:         c.convertValue(mqMsg.Value),
				Timestamp:     mqMsg.Timestamp,
				Quality:       mqMsg.Quality,
				Source:        mqMsg.Source,
			}

			// 发送给 Actor 处理（通过 Execute）
			if c.actorRef != nil {
				if !c.actorRef.Send(valueMsg) {
					// 邮箱满了，可以记录日志或采取其他措施
					// TODO: 添加日志
				}
			}
		}
	}
}

// convertValue 转换 MQ 消息的值类型
func (c *AnalogMeasurementCapacity) convertValue(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	default:
		return 0.0
	}
}

// GetCurrentValue 获取当前测量值
func (c *AnalogMeasurementCapacity) GetCurrentValue() *AnalogMeasurementValue {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentValue
}
