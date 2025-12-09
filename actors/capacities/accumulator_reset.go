package capacities

import (
	"context"
	"fmt"
)

// Message 表示 Actor 之间传递的消息
type Message interface{}

// Capacity 定义了一个能力接口
// 每个 Control 类型对应一个 Capacity 实现
type Capacity interface {
	// Name 返回能力的名称（对应 Control 类型名）
	Name() string

	// CanHandle 检查是否能处理某种消息
	CanHandle(msg Message) bool

	// Execute 执行能力相关的操作
	Execute(ctx context.Context, msg Message) error

	// ControlID 返回关联的 Control 对象的 ID
	ControlID() string
}

// BaseCapacity 提供 Capacity 的基础实现
type BaseCapacity struct {
	controlID string
	name      string
}

func (c *BaseCapacity) Name() string {
	return c.name
}

func (c *BaseCapacity) ControlID() string {
	return c.controlID
}

// AccumulatorResetMessage 累加器复位消息
type AccumulatorResetMessage struct {
	Value int
}

// AccumulatorResetCapacity 实现 AccumulatorReset 控制能力
type AccumulatorResetCapacity struct {
	BaseCapacity
}

// NewAccumulatorResetCapacity 创建新的 AccumulatorResetCapacity
func NewAccumulatorResetCapacity(controlID string) *AccumulatorResetCapacity {
	return &AccumulatorResetCapacity{
		BaseCapacity: BaseCapacity{
			controlID: controlID,
			name:      "AccumulatorResetCapacity",
		},
	}
}

func (c *AccumulatorResetCapacity) CanHandle(msg Message) bool {
	_, ok := msg.(*AccumulatorResetMessage)
	return ok
}

func (c *AccumulatorResetCapacity) Execute(ctx context.Context, msg Message) error {
	resetMsg, ok := msg.(*AccumulatorResetMessage)
	if !ok {
		return fmt.Errorf("invalid message type for AccumulatorResetCapacity")
	}

	// 执行累加器复位逻辑
	// TODO: 实现具体的业务逻辑
	_ = resetMsg
	return nil
}
