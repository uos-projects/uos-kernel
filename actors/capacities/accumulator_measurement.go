package capacities

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/uos-projects/uos-kernel/actors/mq"
)

// AccumulatorMeasurementValueMessage 从 MQ 接收的累加器测量值消息
type AccumulatorMeasurementValueMessage struct {
	MeasurementID string
	Value         int64
	Timestamp     time.Time
	Quality       mq.QualityCode
	Source        string
}

// AccumulatorMeasurementValue 累加器测量值
type AccumulatorMeasurementValue struct {
	Value     int64
	Timestamp time.Time
	Quality   mq.QualityCode
	Source    string
}

// AccumulatorMeasurementCapacity 累加器测量订阅能力
type AccumulatorMeasurementCapacity struct {
	*BaseMeasurementCapacity

	// 当前值缓存
	currentValue *AccumulatorMeasurementValue
	valueMu      sync.RWMutex
}

// NewAccumulatorMeasurementCapacity 创建新的累加器测量能力
func NewAccumulatorMeasurementCapacity(
	measurementID string,
	measurementType string,
	mqConsumer mq.MQConsumer,
) *AccumulatorMeasurementCapacity {
	return &AccumulatorMeasurementCapacity{
		BaseMeasurementCapacity: NewBaseMeasurementCapacity(
			measurementID,
			measurementType,
			"AccumulatorMeasurementCapacity",
			mqConsumer,
		),
	}
}

// StartSubscription 启动 MQ 订阅（重写以启动消费循环）
func (c *AccumulatorMeasurementCapacity) StartSubscription(ctx context.Context) error {
	if err := c.BaseMeasurementCapacity.StartSubscription(ctx); err != nil {
		return err
	}

	// 启动消费循环
	topic := fmt.Sprintf("measurement/%s", c.GetMeasurementID())
	c.BaseMeasurementCapacity.StartConsumeLoop(ctx, topic, c.convertMQMessage)

	return nil
}

// convertMQMessage 将 MQ 消息转换为 AccumulatorMeasurementValueMessage
func (c *AccumulatorMeasurementCapacity) convertMQMessage(mqMsg *mq.MQMessage) Message {
	return &AccumulatorMeasurementValueMessage{
		MeasurementID: c.GetMeasurementID(),
		Value:         c.convertValue(mqMsg.Value),
		Timestamp:     mqMsg.Timestamp,
		Quality:       mqMsg.Quality,
		Source:        mqMsg.Source,
	}
}

// CanHandle 检查是否能处理消息
func (c *AccumulatorMeasurementCapacity) CanHandle(msg Message) bool {
	_, ok := msg.(*AccumulatorMeasurementValueMessage)
	return ok
}

// Execute 处理测量值消息（从 MQ 接收）
func (c *AccumulatorMeasurementCapacity) Execute(ctx context.Context, msg Message) error {
	valueMsg, ok := msg.(*AccumulatorMeasurementValueMessage)
	if !ok {
		return fmt.Errorf("invalid message type for AccumulatorMeasurementCapacity")
	}

	// 更新当前值
	c.valueMu.Lock()
	oldValue := c.currentValue
	c.currentValue = &AccumulatorMeasurementValue{
		Value:     valueMsg.Value,
		Timestamp: valueMsg.Timestamp,
		Quality:   valueMsg.Quality,
		Source:    valueMsg.Source,
	}
	c.valueMu.Unlock()

	// 可以触发业务逻辑（如累计值变化检测、阈值检查等）
	if oldValue != nil && oldValue.Value != valueMsg.Value {
		// 累计值发生变化，可以触发相关业务逻辑
		// TODO: 实现业务逻辑（如累计值变化告警、统计等）
	}

	return nil
}

// convertValue 转换 MQ 消息的值类型为 int64
func (c *AccumulatorMeasurementCapacity) convertValue(v interface{}) int64 {
	switch val := v.(type) {
	case int64:
		return val
	case int:
		return int64(val)
	case int32:
		return int64(val)
	case float64:
		return int64(val)
	case float32:
		return int64(val)
	default:
		return 0
	}
}

// GetCurrentValue 获取当前测量值
func (c *AccumulatorMeasurementCapacity) GetCurrentValue() *AccumulatorMeasurementValue {
	c.valueMu.RLock()
	defer c.valueMu.RUnlock()
	return c.currentValue
}
