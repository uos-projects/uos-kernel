package actors

// ActorRef 是对 Actor 的引用，用于 Actor 之间发送消息
type ActorRef struct {
	system *System
	id     string
}

// NewActorRef 创建一个新的 ActorRef
func NewActorRef(system *System, id string) *ActorRef {
	return &ActorRef{
		system: system,
		id:     id,
	}
}

// ID 返回 Actor 的 ID
func (ref *ActorRef) ID() string {
	return ref.id
}

// Send 向引用的 Actor 发送消息（非阻塞）
func (ref *ActorRef) Send(msg Message) error {
	return ref.system.Send(ref.id, msg)
}

// SendAsync 向引用的 Actor 异步发送消息（阻塞直到成功）
func (ref *ActorRef) SendAsync(msg Message) error {
	return ref.system.SendAsync(ref.id, msg)
}

// Tell 是 Send 的别名，符合 Actor 模型的命名习惯
func (ref *ActorRef) Tell(msg Message) error {
	return ref.Send(msg)
}
