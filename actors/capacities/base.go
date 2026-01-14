package capacities

import (
	"context"
	"fmt"
	"reflect"
)

// Message 表示 Actor 之间传递的消息
type Message interface{}

// ============================================================================
// 1. Command Shape - 可接受的命令类型
// ============================================================================

// CommandType 命令类型描述
type CommandType struct {
	// TypeName 类型名称（如 "SwitchOffCommand"）
	TypeName string
	// MessageType 消息类型的反射类型（用于运行时类型检查）
	MessageType reflect.Type
}

// CommandShape 定义 Capacity 可接受的命令类型
type CommandShape interface {
	// AcceptableCommands 返回可接受的命令类型列表
	AcceptableCommands() []CommandType
}

// ============================================================================
// 2. Guards / Requires - 前置条件
// ============================================================================

// Requirement 前置条件
type Requirement struct {
	// Name 条件名称（用于错误报告）
	Name string
	// Description 条件描述（如 "live_state == online"）
	Description string
	// Check 检查函数，返回是否满足条件
	Check func(ctx context.Context, actor interface{}) (bool, error)
}

// Guards 定义 Capacity 的前置条件
type Guards interface {
	// Requires 返回前置条件列表
	Requires() []Requirement
	// CheckGuards 检查所有前置条件是否满足
	// 返回 (是否全部满足, 不满足的条件名称列表, 错误)
	CheckGuards(ctx context.Context, actor interface{}) (bool, []string, error)
}

// ============================================================================
// 3. Execution Semantics - 执行语义
// ============================================================================

// ExecutionMode 执行模式
type ExecutionMode string

const (
	ExecutionModeSync      ExecutionMode = "sync"      // 同步执行
	ExecutionModeAsync     ExecutionMode = "async"     // 异步执行
	ExecutionModeFireForget ExecutionMode = "fire_forget" // 发射后忘记
)

// ExecutionInfo 执行语义信息
type ExecutionInfo struct {
	// Mode 执行模式（同步/异步/发射后忘记）
	Mode ExecutionMode
	// Interruptible 是否可中断
	Interruptible bool
	// ResourceBound 是否会占用资源（如锁定设备）
	ResourceBound bool
	// Timeout 超时时间（秒，0 表示无超时）
	Timeout int
	// Priority 优先级（数字越大优先级越高）
	Priority int
}

// ExecutionSemantics 定义 Capacity 的执行语义
type ExecutionSemantics interface {
	// ExecutionInfo 返回执行语义信息
	ExecutionInfo() ExecutionInfo
}

// ============================================================================
// 4. Effects - 状态与事件效果
// ============================================================================

// StateChange 状态变化
type StateChange struct {
	// Property 属性名（如 "live_state"）
	Property string
	// FromValue 原值（nil 表示不关心原值）
	FromValue interface{}
	// ToValue 目标值（如 "offline"）
	ToValue interface{}
}

// Event 事件
type Event struct {
	// Type 事件类型（如 "DevicePoweredOff"）
	Type string
	// Payload 事件负载
	Payload interface{}
}

// Effect 效果
type Effect struct {
	// StateChanges 状态变化列表
	StateChanges []StateChange
	// Events 事件列表
	Events []Event
}

// Effects 定义 Capacity 执行后的效果
type Effects interface {
	// Effects 返回执行后的预期效果
	Effects() []Effect
}

// ============================================================================
// Capacity 接口（扩展版，保持向后兼容）
// ============================================================================

// Capacity 定义了一个能力接口
// 这是一个业务中立的接口，不绑定任何特定领域模型
//
// 可选接口（可选实现）：
//   - CommandShape: 声明可接受的命令类型
//   - Guards: 声明前置条件
//   - ExecutionSemantics: 声明执行语义
//   - Effects: 声明执行效果
type Capacity interface {
	// Name 返回能力的名称
	Name() string

	// CanHandle 检查是否能处理某种消息
	CanHandle(msg Message) bool

	// Execute 执行能力相关的操作
	Execute(ctx context.Context, msg Message) error

	// ResourceID 返回关联的资源 ID
	ResourceID() string
}

// ============================================================================
// BaseCapacity 基础实现
// ============================================================================

// BaseCapacity 提供 Capacity 的基础实现
type BaseCapacity struct {
	resourceID string
	name       string
}

// NewBaseCapacity 创建基础 Capacity
func NewBaseCapacity(name, resourceID string) *BaseCapacity {
	return &BaseCapacity{
		resourceID: resourceID,
		name:       name,
	}
}

// Name 返回能力的名称
func (c *BaseCapacity) Name() string {
	return c.name
}

// ResourceID 返回关联的资源 ID
func (c *BaseCapacity) ResourceID() string {
	return c.resourceID
}

// CanHandle 默认实现：检查消息类型是否匹配 AcceptableCommands
// 如果 Capacity 实现了 CommandShape，则使用 CommandShape 进行类型检查
// 否则返回 false（子类需要重写）
// 注意：这个方法应该被子类重写，或者子类应该实现 CommandShape 接口
func (c *BaseCapacity) CanHandle(msg Message) bool {
	// 默认返回 false，子类需要重写
	// 如果子类实现了 CommandShape，可以在子类的 CanHandle 中使用 CommandShape
	return false
}

// ============================================================================
// Guards 默认实现（可选混入）
// ============================================================================

// BaseGuards 提供 Guards 的基础实现
type BaseGuards struct {
	requirements []Requirement
}

// NewBaseGuards 创建基础 Guards
func NewBaseGuards(requirements []Requirement) *BaseGuards {
	return &BaseGuards{
		requirements: requirements,
	}
}

// Requires 返回前置条件列表
func (g *BaseGuards) Requires() []Requirement {
	return g.requirements
}

// CheckGuards 检查所有前置条件是否满足
func (g *BaseGuards) CheckGuards(ctx context.Context, actor interface{}) (bool, []string, error) {
	var failed []string
	for _, req := range g.requirements {
		satisfied, err := req.Check(ctx, actor)
		if err != nil {
			return false, nil, fmt.Errorf("error checking requirement %s: %w", req.Name, err)
		}
		if !satisfied {
			failed = append(failed, req.Name)
		}
	}
	return len(failed) == 0, failed, nil
}
