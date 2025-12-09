package actors

import (
	"context"
	"fmt"
	"sync"
)

// Message 表示 Actor 之间传递的消息
type Message interface{}

// Actor 接口定义了 Actor 的基本行为
type Actor interface {
	// ID 返回 Actor 的唯一标识符
	ID() string

	// Receive 处理接收到的消息
	Receive(ctx context.Context, msg Message) error

	// Start 启动 Actor
	Start(ctx context.Context) error

	// Stop 停止 Actor
	Stop() error
}

// BaseActor 是 Actor 的基础实现
type BaseActor struct {
	id       string
	mailbox  chan Message
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	behavior ActorBehavior
	system   *System // System 引用，用于 Actor 之间通信
}

// ActorBehavior 定义 Actor 的行为逻辑
type ActorBehavior interface {
	// Handle 处理消息，返回是否继续处理
	Handle(ctx context.Context, msg Message) (bool, error)
}

// NewBaseActor 创建一个新的基础 Actor
func NewBaseActor(id string, behavior ActorBehavior) *BaseActor {
	ctx, cancel := context.WithCancel(context.Background())
	return &BaseActor{
		id:       id,
		mailbox:  make(chan Message, 100), // 带缓冲的邮箱
		ctx:      ctx,
		cancel:   cancel,
		behavior: behavior,
		system:   nil, // 将在注册到 System 时设置
	}
}

// SetSystem 设置 Actor 的 System 引用（由 System 在注册时调用）
func (a *BaseActor) SetSystem(system *System) {
	a.system = system
}

// GetRef 获取指定 Actor 的 ActorRef
func (a *BaseActor) GetRef(targetID string) (*ActorRef, error) {
	if a.system == nil {
		return nil, fmt.Errorf("actor %s is not registered in a system", a.id)
	}
	return a.system.GetRef(targetID)
}

// SendTo 向指定 ID 的 Actor 发送消息
func (a *BaseActor) SendTo(targetID string, msg Message) error {
	if a.system == nil {
		return fmt.Errorf("actor %s is not registered in a system", a.id)
	}
	return a.system.Send(targetID, msg)
}

// Tell 是 SendTo 的别名，符合 Actor 模型的命名习惯
func (a *BaseActor) Tell(targetID string, msg Message) error {
	return a.SendTo(targetID, msg)
}

// ID 返回 Actor 的 ID
func (a *BaseActor) ID() string {
	return a.id
}

// Send 向 Actor 发送消息（非阻塞）
func (a *BaseActor) Send(msg Message) bool {
	select {
	case a.mailbox <- msg:
		return true
	default:
		return false // 邮箱满了
	}
}

// SendAsync 异步发送消息（阻塞直到成功）
func (a *BaseActor) SendAsync(msg Message) {
	a.mailbox <- msg
}

// Receive 实现 Actor 接口
func (a *BaseActor) Receive(ctx context.Context, msg Message) error {
	if a.behavior != nil {
		_, err := a.behavior.Handle(ctx, msg)
		return err
	}
	return nil
}

// Start 启动 Actor
func (a *BaseActor) Start(ctx context.Context) error {
	a.wg.Add(1)
	go a.run()
	return nil
}

// run 是 Actor 的主循环
func (a *BaseActor) run() {
	defer a.wg.Done()

	for {
		select {
		case <-a.ctx.Done():
			// 处理剩余消息
			for {
				select {
				case msg := <-a.mailbox:
					a.Receive(a.ctx, msg)
				default:
					return
				}
			}
		case msg := <-a.mailbox:
			if err := a.Receive(a.ctx, msg); err != nil {
				// 错误处理逻辑可以在这里添加
				continue
			}
		}
	}
}

// Stop 停止 Actor
func (a *BaseActor) Stop() error {
	a.cancel()
	a.wg.Wait()
	close(a.mailbox)
	return nil
}
