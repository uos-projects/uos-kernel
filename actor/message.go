package actor

// ============================================================================
// 消息分类：三类消息
// ============================================================================

// Message 表示 Actor 之间传递的消息（基础接口）
type Message interface {
	// MessageType 返回消息类型（用于分类）
	MessageType() MessageCategory
}

// MessageCategory 消息类别
type MessageCategory string

const (
	// MessageCategoryCapabilityCommand 能力命令（对外、可调度）
	MessageCategoryCapabilityCommand MessageCategory = "capability_command"
	
	// MessageCategoryCoordinationEvent 协同事件（来自其他 Actor）
	MessageCategoryCoordinationEvent MessageCategory = "coordination_event"
	
	// MessageCategoryInternal 内部消息（定时、重试、生命周期）
	MessageCategoryInternal MessageCategory = "internal"
)

// ============================================================================
// 1. Capability Commands（对外、可调度）
// ============================================================================

// CapabilityCommand 能力命令消息（基础接口）
// 所有能力命令消息都应该实现此接口
type CapabilityCommand interface {
	Message
	// CommandID 返回命令 ID（用于跟踪）
	CommandID() string
}

// ============================================================================
// 2. Coordination Events（来自其他 Actor）
// ============================================================================

// CoordinationEvent 协同事件消息（基础接口）
// 用于 Actor 之间的协同通信
type CoordinationEvent interface {
	Message
	// SourceActorID 返回发送事件的 Actor ID
	SourceActorID() string
}

// ============================================================================
// 3. Internal Messages（定时、重试、生命周期）
// ============================================================================

// InternalMessage 内部消息（基础接口）
// 用于 Actor 内部状态管理、生命周期等
type InternalMessage interface {
	Message
}

// ============================================================================
// Actor 状态管理消息（Internal Messages）
// ============================================================================

// SetPropertyMessage 设置属性消息
// 用于通过消息驱动的方式修改 Actor 的属性
type SetPropertyMessage struct {
	Name  string
	Value interface{}
}

func (m *SetPropertyMessage) MessageType() MessageCategory {
	return MessageCategoryInternal
}

// GetPropertyMessage 获取属性消息
type GetPropertyMessage struct {
	Name string
}

func (m *GetPropertyMessage) MessageType() MessageCategory {
	return MessageCategoryInternal
}

// GetPropertyResponse 获取属性响应
type GetPropertyResponse struct {
	Name   string
	Value  interface{}
	Exists bool
}

func (m *GetPropertyResponse) MessageType() MessageCategory {
	return MessageCategoryInternal
}

// TransitionStateMessage 状态转换消息
// 用于通过消息驱动的方式转换 Actor 的生命周期状态
type TransitionStateMessage struct {
	ToState LifecycleState
	Trigger string
}

func (m *TransitionStateMessage) MessageType() MessageCategory {
	return MessageCategoryInternal
}

// GetStateMessage 获取状态消息
type GetStateMessage struct{}

func (m *GetStateMessage) MessageType() MessageCategory {
	return MessageCategoryInternal
}

// GetStateResponse 获取状态响应
type GetStateResponse struct {
	CurrentState LifecycleState
	StateHistory []StateTransition
}

func (m *GetStateResponse) MessageType() MessageCategory {
	return MessageCategoryInternal
}

// SuspendMessage 挂起消息
type SuspendMessage struct {
	Reason string
}

func (m *SuspendMessage) MessageType() MessageCategory {
	return MessageCategoryInternal
}

// ResumeMessage 恢复消息
type ResumeMessage struct {
	Reason string
}

func (m *ResumeMessage) MessageType() MessageCategory {
	return MessageCategoryInternal
}

// ============================================================================
// Binding 相关消息（Coordination Events）
// ============================================================================

// ExternalEventMessage 外部事件消息
// 由 Binding 接收外部事件后转换为消息发送给 Actor
type ExternalEventMessage struct {
	BindingType BindingType
	Event       interface{}
}

func (m *ExternalEventMessage) MessageType() MessageCategory {
	return MessageCategoryCoordinationEvent
}

func (m *ExternalEventMessage) SourceActorID() string {
	return "" // 外部事件没有源 Actor
}

// ExecuteExternalCommandMessage 执行外部命令消息
// Actor 需要执行外部操作时发送此消息
type ExecuteExternalCommandMessage struct {
	BindingType BindingType
	Command     interface{}
}

func (m *ExecuteExternalCommandMessage) MessageType() MessageCategory {
	return MessageCategoryCoordinationEvent
}

func (m *ExecuteExternalCommandMessage) SourceActorID() string {
	return "" // 内部消息，没有源 Actor
}
