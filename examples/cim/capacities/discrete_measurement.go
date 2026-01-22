package capacities



import (
	"context"
	"fmt"
	"github.com/uos-projects/uos-kernel/actors/capacities"
	"github.com/uos-projects/uos-kernel/actors/mq"
	"sync"
	"time"
)

// DiscreteState 离散状态
type DiscreteState int64

const (
	// DiscreteStateUnknown 未知状态
	DiscreteStateUnknown DiscreteState = -1
	// DiscreteStateOff 关闭状态（通常为 0）
	DiscreteStateOff DiscreteState = 0
	// DiscreteStateOn 开启状态（通常为 1）
	DiscreteStateOn DiscreteState = 1
)

// String 返回状态的字符串表示
func (s DiscreteState) String() string {
	switch s {
	case DiscreteStateOff:
		return "OFF"
	case DiscreteStateOn:
		return "ON"
	default:
		return "UNKNOWN"
	}
}

// StateTransition 状态转换记录
type StateTransition struct {
	FromState DiscreteState
	ToState   DiscreteState
	Timestamp time.Time
	Value     int64
	Quality   mq.QualityCode
	Source    string
}

// DiscreteMeasurementValueMessage 从 MQ 接收的离散量测量值消息
type DiscreteMeasurementValueMessage struct {
	MeasurementID string
	Value         int64
	Timestamp     time.Time
	Quality       mq.QualityCode
	Source        string
}

// DiscreteMeasurementValue 离散量测量值
type DiscreteMeasurementValue struct {
	Value     int64
	State     DiscreteState
	Timestamp time.Time
	Quality   mq.QualityCode
	Source    string
}

// DiscreteMeasurementCapacity 离散量测量订阅能力（带状态机）
type DiscreteMeasurementCapacity struct {
	*BaseMeasurementCapacity

	// 状态机相关
	currentValue *DiscreteMeasurementValue
	currentState DiscreteState
	// 状态转换历史（可选，用于审计和调试）
	stateHistory []StateTransition
	maxHistory   int // 最大历史记录数
	stateMu      sync.RWMutex

	// 状态转换回调（可选）
	onStateChange func(from, to DiscreteState, timestamp time.Time)
}

// NewDiscreteMeasurementCapacity 创建新的离散量测量能力
func NewDiscreteMeasurementCapacity(
	measurementID string,
	measurementType string,
	mqConsumer mq.MQConsumer,
) *DiscreteMeasurementCapacity {
	return &DiscreteMeasurementCapacity{
		BaseMeasurementCapacity: NewBaseMeasurementCapacity(
			measurementID,
			measurementType,
			"DiscreteMeasurementCapacity",
			mqConsumer,
		),
		currentState: DiscreteStateUnknown,
		maxHistory:   100, // 默认保留最近 100 条状态转换记录
		stateHistory: make([]StateTransition, 0, 100),
	}
}

// StartSubscription 启动 MQ 订阅（重写以启动消费循环）
func (c *DiscreteMeasurementCapacity) StartSubscription(ctx context.Context) error {
	if err := c.BaseMeasurementCapacity.StartSubscription(ctx); err != nil {
		return err
	}

	// 启动消费循环
	topic := fmt.Sprintf("measurement/%s", c.GetMeasurementID())
	c.BaseMeasurementCapacity.StartConsumeLoop(ctx, topic, c.convertMQMessage)

	return nil
}

// convertMQMessage 将 MQ 消息转换为 DiscreteMeasurementValueMessage
func (c *DiscreteMeasurementCapacity) convertMQMessage(mqMsg *mq.MQMessage) capacities.Message {
	return &DiscreteMeasurementValueMessage{
		MeasurementID: c.GetMeasurementID(),
		Value:         c.convertValue(mqMsg.Value),
		Timestamp:     mqMsg.Timestamp,
		Quality:       mqMsg.Quality,
		Source:        mqMsg.Source,
	}
}

// SetOnStateChange 设置状态转换回调函数
func (c *DiscreteMeasurementCapacity) SetOnStateChange(fn func(from, to DiscreteState, timestamp time.Time)) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.onStateChange = fn
}

// CanHandle 检查是否能处理消息
func (c *DiscreteMeasurementCapacity) CanHandle(msg capacities.Message) bool {
	_, ok := msg.(*DiscreteMeasurementValueMessage)
	return ok
}

// Execute 处理测量值消息（从 MQ 接收）
func (c *DiscreteMeasurementCapacity) Execute(ctx context.Context, msg capacities.Message) error {
	valueMsg, ok := msg.(*DiscreteMeasurementValueMessage)
	if !ok {
		return fmt.Errorf("invalid message type for DiscreteMeasurementCapacity")
	}

	// 转换值为状态
	newState := c.valueToState(valueMsg.Value)

	// 更新状态机
	c.stateMu.Lock()
	oldState := c.currentState
	stateChanged := oldState != newState

	if stateChanged {
		// 记录状态转换
		transition := StateTransition{
			FromState: oldState,
			ToState:   newState,
			Timestamp: valueMsg.Timestamp,
			Value:     valueMsg.Value,
			Quality:   valueMsg.Quality,
			Source:    valueMsg.Source,
		}

		// 添加到历史记录
		c.stateHistory = append(c.stateHistory, transition)
		if len(c.stateHistory) > c.maxHistory {
			// 移除最旧的记录
			c.stateHistory = c.stateHistory[1:]
		}

		// 更新当前状态
		c.currentState = newState
	}

	// 更新当前值
	c.currentValue = &DiscreteMeasurementValue{
		Value:     valueMsg.Value,
		State:     newState,
		Timestamp: valueMsg.Timestamp,
		Quality:   valueMsg.Quality,
		Source:    valueMsg.Source,
	}

	onStateChange := c.onStateChange
	c.stateMu.Unlock()

	// 触发状态转换回调（在锁外调用，避免死锁）
	if stateChanged && onStateChange != nil {
		onStateChange(oldState, newState, valueMsg.Timestamp)
	}

	// 可以触发业务逻辑（如状态变化告警、状态机验证等）
	if stateChanged {
		// TODO: 实现业务逻辑（如状态变化告警、状态机验证等）
	}

	return nil
}

// valueToState 将整数值转换为状态
func (c *DiscreteMeasurementCapacity) valueToState(value int64) DiscreteState {
	switch value {
	case 0:
		return DiscreteStateOff
	case 1:
		return DiscreteStateOn
	default:
		// 对于其他值，可以根据业务需求扩展
		// 这里返回未知状态，实际应用中可能需要更复杂的映射
		return DiscreteStateUnknown
	}
}

// GetCurrentState 获取当前状态
func (c *DiscreteMeasurementCapacity) GetCurrentState() DiscreteState {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.currentState
}

// GetCurrentValue 获取当前测量值
func (c *DiscreteMeasurementCapacity) GetCurrentValue() *DiscreteMeasurementValue {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.currentValue
}

// GetStateHistory 获取状态转换历史
func (c *DiscreteMeasurementCapacity) GetStateHistory() []StateTransition {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	// 返回副本，避免外部修改
	history := make([]StateTransition, len(c.stateHistory))
	copy(history, c.stateHistory)
	return history
}

// GetRecentTransitions 获取最近 N 次状态转换
func (c *DiscreteMeasurementCapacity) GetRecentTransitions(n int) []StateTransition {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	if n <= 0 || len(c.stateHistory) == 0 {
		return nil
	}
	start := len(c.stateHistory) - n
	if start < 0 {
		start = 0
	}
	history := make([]StateTransition, len(c.stateHistory)-start)
	copy(history, c.stateHistory[start:])
	return history
}

// convertValue 转换 MQ 消息的值类型为 int64
func (c *DiscreteMeasurementCapacity) convertValue(v interface{}) int64 {
	switch val := v.(type) {
	case int64:
		return val
	case int:
		return int64(val)
	case int32:
		return int64(val)
	case bool:
		if val {
			return 1
		}
		return 0
	case float64:
		return int64(val)
	case float32:
		return int64(val)
	default:
		return 0
	}
}
