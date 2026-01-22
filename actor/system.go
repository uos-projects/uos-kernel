package actor

import (
	"context"
	"fmt"
	"sync"
)

// System 是 Actor 系统的管理器
type System struct {
	actors           map[string]Actor
	eventSubscribers []string // 订阅事件的 Actor ID 列表
	mu               sync.RWMutex
	ctx              context.Context
	cancel           context.CancelFunc
}

// NewSystem 创建一个新的 Actor 系统
func NewSystem(ctx context.Context) *System {
	sysCtx, cancel := context.WithCancel(ctx)
	return &System{
		actors:           make(map[string]Actor),
		eventSubscribers: make([]string, 0),
		ctx:              sysCtx,
		cancel:           cancel,
	}
}

// SubscribeEvent 订阅事件（Actor 订阅接收所有事件）
func (s *System) SubscribeEvent(actorID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查 Actor 是否存在
	if _, exists := s.actors[actorID]; !exists {
		return fmt.Errorf("actor %s not found", actorID)
	}

	// 检查是否已经订阅
	for _, id := range s.eventSubscribers {
		if id == actorID {
			return nil // 已经订阅，直接返回
		}
	}

	s.eventSubscribers = append(s.eventSubscribers, actorID)
	return nil
}

// UnsubscribeEvent 取消订阅事件
func (s *System) UnsubscribeEvent(actorID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, id := range s.eventSubscribers {
		if id == actorID {
			s.eventSubscribers = append(s.eventSubscribers[:i], s.eventSubscribers[i+1:]...)
			break
		}
	}
}

// PublishEvent 发布事件（直接发送给所有订阅的 Actor）
func (s *System) PublishEvent(ctx context.Context, event Event) error {
	s.mu.RLock()
	subscribers := make([]string, len(s.eventSubscribers))
	copy(subscribers, s.eventSubscribers)
	s.mu.RUnlock()

	if len(subscribers) == 0 {
		// 没有订阅者，直接返回
		return nil
	}

	// 直接通过 system.Send() 发送给所有订阅的 Actor
	for _, actorID := range subscribers {
		go func(id string) {
			if err := s.Send(id, event); err != nil {
				// 静默处理错误，避免影响其他订阅者
				_ = err
			}
		}(actorID)
	}

	return nil
}

// Register 注册一个已创建的 Actor（用于 ResourceActor 等）
// 注意：Register 不会自动启动 Actor，需要显式调用 Start()
func (s *System) Register(actor Actor) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := actor.ID()
	if _, exists := s.actors[id]; exists {
		return fmt.Errorf("actor %s already exists", id)
	}

	// 设置 System 引用和 self 引用
	// 优先处理 BaseResourceActor（因为它重写了 SetSystem）
	// 支持直接类型和嵌入类型（通过嵌入字段的方法提升）
	if resourceActor, ok := actor.(*BaseResourceActor); ok {
		resourceActor.SetSystem(s)
		// 设置 self 引用
		if baseActor := resourceActor.BaseActor; baseActor != nil {
			baseActor.SetSelf(actor)
		}
	} else if baseActor, ok := actor.(*BaseActor); ok {
		baseActor.SetSystem(s)
		baseActor.SetSelf(actor)
	} else {
		// 对于嵌入 BaseResourceActor 的类型，由于 Go 的方法提升机制，
		// 如果 BaseResourceActor 有 SetSystem 方法，嵌入类型可以直接调用
		// 但类型断言不会成功，所以我们需要通过接口来调用
		// 尝试调用 SetSystem（如果 Actor 实现了相应接口）
		if setter, ok := actor.(interface{ SetSystem(*System) }); ok {
			setter.SetSystem(s)
		}
		// 尝试设置 self 引用
		if selfSetter, ok := actor.(interface{ SetSelf(Actor) }); ok {
			selfSetter.SetSelf(actor)
		} else {
			// 如果 Actor 嵌入了 BaseResourceActor 或 BaseActor，尝试通过反射设置
			// 这里简化处理：如果 Actor 有 BaseActor 字段，尝试设置
			if baseActorGetter, ok := actor.(interface{ GetBaseActor() *BaseActor }); ok {
				if baseActor := baseActorGetter.GetBaseActor(); baseActor != nil {
					baseActor.SetSelf(actor)
				}
			}
		}
	}

	s.actors[id] = actor
	return nil
}

// Get 获取指定 ID 的 Actor
func (s *System) Get(id string) (Actor, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	actor, exists := s.actors[id]
	return actor, exists
}

// GetRef 获取指定 ID 的 ActorRef，用于 Actor 之间发送消息
func (s *System) GetRef(id string) (*ActorRef, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.actors[id]; !exists {
		return nil, fmt.Errorf("actor %s not found", id)
	}

	return NewActorRef(s, id), nil
}

// MustGetRef 获取指定 ID 的 ActorRef，如果不存在则 panic
func (s *System) MustGetRef(id string) *ActorRef {
	ref, err := s.GetRef(id)
	if err != nil {
		panic(err)
	}
	return ref
}

// Send 向指定 ID 的 Actor 发送消息
func (s *System) Send(id string, msg Message) error {
	s.mu.RLock()
	actor, exists := s.actors[id]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("actor %s not found", id)
	}

	// 优先使用 ResourceActor 接口（包含 Send 方法）
	if resourceActor, ok := actor.(ResourceActor); ok {
		if !resourceActor.Send(msg) {
			return fmt.Errorf("failed to send message to actor %s: mailbox full", id)
		}
		return nil
	}

	// 回退到 BaseResourceActor（直接类型）
	if resourceActor, ok := actor.(*BaseResourceActor); ok {
		if !resourceActor.Send(msg) {
			return fmt.Errorf("failed to send message to actor %s: mailbox full", id)
		}
		return nil
	}

	// 回退到 BaseActor
	if baseActor, ok := actor.(*BaseActor); ok {
		if !baseActor.Send(msg) {
			return fmt.Errorf("failed to send message to actor %s: mailbox full", id)
		}
		return nil
	}

	// 最后尝试：如果 Actor 有 Send 方法（通过嵌入获得，如 BreakerActor, DispatcherActor）
	if sender, ok := actor.(interface{ Send(Message) bool }); ok {
		if !sender.Send(msg) {
			return fmt.Errorf("failed to send message to actor %s: mailbox full", id)
		}
		return nil
	}

	return fmt.Errorf("actor %s has unsupported type: %T", id, actor)
}

// SendAsync 异步发送消息（阻塞直到成功）
func (s *System) SendAsync(id string, msg Message) error {
	s.mu.RLock()
	actor, exists := s.actors[id]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("actor %s not found", id)
	}

	// 优先使用 BaseResourceActor（包含 SendAsync 方法）
	if resourceActor, ok := actor.(*BaseResourceActor); ok {
		resourceActor.SendAsync(msg)
		return nil
	}

	// 回退到 BaseActor
	if baseActor, ok := actor.(*BaseActor); ok {
		baseActor.SendAsync(msg)
		return nil
	}

	return fmt.Errorf("actor %s has unsupported type", id)
}

// Stop 停止指定 ID 的 Actor
func (s *System) Stop(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	actor, exists := s.actors[id]
	if !exists {
		return fmt.Errorf("actor %s not found", id)
	}

	if err := actor.Stop(); err != nil {
		return err
	}

	delete(s.actors, id)
	return nil
}

// Shutdown 关闭整个 Actor 系统
func (s *System) Shutdown() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cancel()

	var errs []error
	for id, actor := range s.actors {
		if err := actor.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("failed to stop actor %s: %w", id, err))
		}
	}

	s.actors = make(map[string]Actor)

	if len(errs) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errs)
	}

	return nil
}

// Count 返回当前系统中的 Actor 数量
func (s *System) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.actors)
}
