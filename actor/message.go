package actor

// ============================================================================
// 消息分类：两类消息
// ============================================================================

// Message 表示 Actor 之间传递的消息（基础接口）
type Message interface {
	// MessageType 返回消息类型（用于分类）
	MessageType() MessageCategory
}

// MessageCategory 消息类别
type MessageCategory string

const (
	// MessageCategoryCommand 命令（对外、可调度）
	MessageCategoryCommand MessageCategory = "command"

	// MessageCategoryEvent 事件（来自其他 Actor）
	MessageCategoryEvent MessageCategory = "event"
)

// ============================================================================
// 1. Command（对外、可调度）
// ============================================================================

// Command 命令消息（基础接口）
// 所有命令消息都应该实现此接口
type Command interface {
	Message
	// CommandID 返回命令 ID（用于跟踪）
	CommandID() string
}

// ============================================================================
// 2. Event（来自其他 Actor）
// ============================================================================

// Event 事件消息（基础接口）
// 用于 Actor 之间的协同通信
type Event interface {
	Message
	// SourceActorID 返回发送事件的 Actor ID
	SourceActorID() string
}

// ============================================================================
// Binding 相关消息（Events）
// ============================================================================

// ExternalEventMessage 外部事件消息
// 由 Binding 接收外部事件后转换为消息发送给 Actor
type ExternalEventMessage struct {
	BindingType BindingType
	Event       interface{}
}

func (m *ExternalEventMessage) MessageType() MessageCategory {
	return MessageCategoryEvent
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
	return MessageCategoryEvent
}

func (m *ExecuteExternalCommandMessage) SourceActorID() string {
	return "" // 内部消息，没有源 Actor
}
