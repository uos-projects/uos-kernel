package capacities

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/uos-projects/uos-kernel/actors/mq"
)

// StringMeasurementValueMessage 从 MQ 接收的字符串测量值消息
type StringMeasurementValueMessage struct {
	MeasurementID string
	Value         string
	Timestamp     time.Time
	Quality       mq.QualityCode
	Source        string
}

// StringMeasurementValue 字符串测量值
type StringMeasurementValue struct {
	Value     string
	Timestamp time.Time
	Quality   mq.QualityCode
	Source    string
}

// StringMeasurementCapacity 字符串测量订阅能力
type StringMeasurementCapacity struct {
	*BaseMeasurementCapacity

	// 当前值缓存
	currentValue *StringMeasurementValue
	// 值变化历史（可选，用于审计）
	valueHistory []StringMeasurementValue
	maxHistory   int // 最大历史记录数
	valueMu      sync.RWMutex
}

// NewStringMeasurementCapacity 创建新的字符串测量能力
func NewStringMeasurementCapacity(
	measurementID string,
	measurementType string,
	mqConsumer mq.MQConsumer,
) *StringMeasurementCapacity {
	return &StringMeasurementCapacity{
		BaseMeasurementCapacity: NewBaseMeasurementCapacity(
			measurementID,
			measurementType,
			"StringMeasurementCapacity",
			mqConsumer,
		),
		maxHistory:   100, // 默认保留最近 100 条历史记录
		valueHistory: make([]StringMeasurementValue, 0, 100),
	}
}

// StartSubscription 启动 MQ 订阅（重写以启动消费循环）
func (c *StringMeasurementCapacity) StartSubscription(ctx context.Context) error {
	if err := c.BaseMeasurementCapacity.StartSubscription(ctx); err != nil {
		return err
	}

	// 启动消费循环
	topic := fmt.Sprintf("measurement/%s", c.GetMeasurementID())
	c.BaseMeasurementCapacity.StartConsumeLoop(ctx, topic, c.convertMQMessage)

	return nil
}

// convertMQMessage 将 MQ 消息转换为 StringMeasurementValueMessage
func (c *StringMeasurementCapacity) convertMQMessage(mqMsg *mq.MQMessage) Message {
	return &StringMeasurementValueMessage{
		MeasurementID: c.GetMeasurementID(),
		Value:         c.convertValue(mqMsg.Value),
		Timestamp:     mqMsg.Timestamp,
		Quality:       mqMsg.Quality,
		Source:        mqMsg.Source,
	}
}

// CanHandle 检查是否能处理消息
func (c *StringMeasurementCapacity) CanHandle(msg Message) bool {
	_, ok := msg.(*StringMeasurementValueMessage)
	return ok
}

// Execute 处理测量值消息（从 MQ 接收）
func (c *StringMeasurementCapacity) Execute(ctx context.Context, msg Message) error {
	valueMsg, ok := msg.(*StringMeasurementValueMessage)
	if !ok {
		return fmt.Errorf("invalid message type for StringMeasurementCapacity")
	}

	// 更新当前值
	c.valueMu.Lock()
	oldValue := c.currentValue
	c.currentValue = &StringMeasurementValue{
		Value:     valueMsg.Value,
		Timestamp: valueMsg.Timestamp,
		Quality:   valueMsg.Quality,
		Source:    valueMsg.Source,
	}

	// 记录值变化历史
	if oldValue == nil || oldValue.Value != valueMsg.Value {
		c.valueHistory = append(c.valueHistory, *c.currentValue)
		if len(c.valueHistory) > c.maxHistory {
			// 移除最旧的记录
			c.valueHistory = c.valueHistory[1:]
		}
	}
	c.valueMu.Unlock()

	// 可以触发业务逻辑（如字符串值变化检测、模式匹配等）
	if oldValue != nil && oldValue.Value != valueMsg.Value {
		// 字符串值发生变化，可以触发相关业务逻辑
		// TODO: 实现业务逻辑（如字符串模式匹配、告警等）
	}

	return nil
}

// convertValue 转换 MQ 消息的值类型为 string
func (c *StringMeasurementCapacity) convertValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case fmt.Stringer:
		return val.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}

// GetCurrentValue 获取当前测量值
func (c *StringMeasurementCapacity) GetCurrentValue() *StringMeasurementValue {
	c.valueMu.RLock()
	defer c.valueMu.RUnlock()
	return c.currentValue
}

// GetValueHistory 获取值变化历史
func (c *StringMeasurementCapacity) GetValueHistory() []StringMeasurementValue {
	c.valueMu.RLock()
	defer c.valueMu.RUnlock()
	// 返回副本，避免外部修改
	history := make([]StringMeasurementValue, len(c.valueHistory))
	copy(history, c.valueHistory)
	return history
}

// GetRecentValues 获取最近 N 个值
func (c *StringMeasurementCapacity) GetRecentValues(n int) []StringMeasurementValue {
	c.valueMu.RLock()
	defer c.valueMu.RUnlock()
	if n <= 0 || len(c.valueHistory) == 0 {
		return nil
	}
	start := len(c.valueHistory) - n
	if start < 0 {
		start = 0
	}
	history := make([]StringMeasurementValue, len(c.valueHistory)-start)
	copy(history, c.valueHistory[start:])
	return history
}
