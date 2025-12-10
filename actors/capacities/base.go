package capacities

import (
	"context"
)

// Message 表示 Actor 之间传递的消息
type Message interface{}

// Capacity 定义了一个能力接口
// 可以是 Control（执行能力）或 Measurement（订阅能力）
type Capacity interface {
	// Name 返回能力的名称
	Name() string

	// CanHandle 检查是否能处理某种消息
	CanHandle(msg Message) bool

	// Execute 执行能力相关的操作
	// 对于 Control：执行控制命令
	// 对于 Measurement：处理测量值更新
	Execute(ctx context.Context, msg Message) error

	// ResourceID 返回关联的资源 ID
	// 对于 Control：返回 ControlID
	// 对于 Measurement：返回 MeasurementID
	ResourceID() string
}

// BaseCapacity 提供 Capacity 的基础实现
type BaseCapacity struct {
	resourceID string // 通用字段名，可以是 ControlID 或 MeasurementID
	name       string
}

// Name 返回能力的名称
func (c *BaseCapacity) Name() string {
	return c.name
}

// ResourceID 返回关联的资源 ID
func (c *BaseCapacity) ResourceID() string {
	return c.resourceID
}

// ControlID 返回 Control ID（向后兼容）
func (c *BaseCapacity) ControlID() string {
	return c.resourceID
}

// MeasurementID 返回 Measurement ID（向后兼容）
func (c *BaseCapacity) MeasurementID() string {
	return c.resourceID
}

