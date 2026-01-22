package actor

import (
	"context"
	"fmt"
	"sync"
)

// System 是 Actor 系统的管理器
type System struct {
	actors map[string]Actor
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

// NewSystem 创建一个新的 Actor 系统
func NewSystem(ctx context.Context) *System {
	sysCtx, cancel := context.WithCancel(ctx)
	return &System{
		actors: make(map[string]Actor),
		ctx:    sysCtx,
		cancel: cancel,
	}
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

	// 设置 System 引用
	// 优先处理 BaseResourceActor（因为它重写了 SetSystem）
	if resourceActor, ok := actor.(*BaseResourceActor); ok {
		resourceActor.SetSystem(s)
	} else if baseActor, ok := actor.(*BaseActor); ok {
		baseActor.SetSystem(s)
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

	// 回退到 BaseActor
	if baseActor, ok := actor.(*BaseActor); ok {
		if !baseActor.Send(msg) {
			return fmt.Errorf("failed to send message to actor %s: mailbox full", id)
		}
		return nil
	}

	return fmt.Errorf("actor %s has unsupported type", id)
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
