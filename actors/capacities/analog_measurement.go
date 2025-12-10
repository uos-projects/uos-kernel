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
	*BaseMeasurementCapacity

	// 当前值缓存
	currentValue *AnalogMeasurementValue
	valueMu      sync.RWMutex
}

// NewAnalogMeasurementCapacity 创建新的模拟量测量能力
func NewAnalogMeasurementCapacity(
	measurementID string,
	measurementType string,
	mqConsumer mq.MQConsumer,
) *AnalogMeasurementCapacity {
	return &AnalogMeasurementCapacity{
		BaseMeasurementCapacity: NewBaseMeasurementCapacity(
			measurementID,
			measurementType,
			"AnalogMeasurementCapacity",
			mqConsumer,
		),
	}
}

// StartSubscription 启动 MQ 订阅（重写以启动消费循环）
func (c *AnalogMeasurementCapacity) StartSubscription(ctx context.Context) error {
	if err := c.BaseMeasurementCapacity.StartSubscription(ctx); err != nil {
		return err
	}

	// 启动消费循环
	topic := fmt.Sprintf("measurement/%s", c.GetMeasurementID())
	c.BaseMeasurementCapacity.StartConsumeLoop(ctx, topic, c.convertMQMessage)

	return nil
}

// convertMQMessage 将 MQ 消息转换为 AnalogMeasurementValueMessage
func (c *AnalogMeasurementCapacity) convertMQMessage(mqMsg *mq.MQMessage) Message {
	return &AnalogMeasurementValueMessage{
		MeasurementID: c.GetMeasurementID(),
		Value:         c.convertValue(mqMsg.Value),
		Timestamp:     mqMsg.Timestamp,
		Quality:       mqMsg.Quality,
		Source:        mqMsg.Source,
	}
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
	c.valueMu.Lock()
	c.currentValue = &AnalogMeasurementValue{
		Value:     valueMsg.Value,
		Timestamp: valueMsg.Timestamp,
		Quality:   valueMsg.Quality,
		Source:    valueMsg.Source,
	}
	c.valueMu.Unlock()

	// 可以触发业务逻辑（如阈值检查、告警等）
	// TODO: 实现业务逻辑

	return nil
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
	c.valueMu.RLock()
	defer c.valueMu.RUnlock()
	return c.currentValue
}
