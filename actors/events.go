package actors

import (
	"context"
	"sync"
	"time"
)

// ============================================================================
// 事件输出（Event Emission）
// ============================================================================

// EventType 事件类型
type EventType string

const (
	// 状态变化事件
	EventTypeStateChanged EventType = "state_changed"
	
	// 执行完成/失败事件
	EventTypeCommandCompleted EventType = "command_completed"
	EventTypeCommandFailed    EventType = "command_failed"
	
	// 生命周期事件
	EventTypeActorCreated  EventType = "actor_created"
	EventTypeActorStarted  EventType = "actor_started"
	EventTypeActorStopped  EventType = "actor_stopped"
	EventTypeActorSuspended EventType = "actor_suspended"
	EventTypeActorResumed  EventType = "actor_resumed"
)

// Event 事件
type Event struct {
	// Type 事件类型
	Type EventType
	
	// ActorID 产生事件的 Actor ID
	ActorID string
	
	// Timestamp 事件时间戳
	Timestamp time.Time
	
	// Payload 事件负载（具体事件数据）
	Payload interface{}
}

// EventEmitter 事件发射器接口
type EventEmitter interface {
	// Emit 发射事件
	Emit(event Event) error
	
	// EmitStateChanged 发射状态变化事件
	EmitStateChanged(from, to LifecycleState, trigger string) error
	
	// EmitCommandCompleted 发射命令完成事件
	EmitCommandCompleted(commandID string, result interface{}) error
	
	// EmitCommandFailed 发射命令失败事件
	EmitCommandFailed(commandID string, err error) error
}

// EventHandler 事件处理器接口（用于订阅事件）
type EventHandler interface {
	// HandleEvent 处理事件
	HandleEvent(ctx context.Context, event Event) error
}

// BaseEventEmitter 基础事件发射器实现
type BaseEventEmitter struct {
	actorID    string
	handlers   []EventHandler
	system     *System // 用于向其他 Actor 发送事件
	mu         sync.RWMutex
}

// NewBaseEventEmitter 创建基础事件发射器
func NewBaseEventEmitter(actorID string, system *System) *BaseEventEmitter {
	return &BaseEventEmitter{
		actorID:  actorID,
		handlers: make([]EventHandler, 0),
		system:   system,
	}
}

// AddHandler 添加事件处理器
func (e *BaseEventEmitter) AddHandler(handler EventHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers = append(e.handlers, handler)
}

// RemoveHandler 移除事件处理器
func (e *BaseEventEmitter) RemoveHandler(handler EventHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for i, h := range e.handlers {
		if h == handler {
			e.handlers = append(e.handlers[:i], e.handlers[i+1:]...)
			break
		}
	}
}

// Emit 发射事件
func (e *BaseEventEmitter) Emit(event Event) error {
	// 设置 ActorID 和时间戳（如果未设置）
	if event.ActorID == "" {
		event.ActorID = e.actorID
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	
	// 通知所有处理器
	e.mu.RLock()
	handlers := make([]EventHandler, len(e.handlers))
	copy(handlers, e.handlers)
	e.mu.RUnlock()
	
	for _, handler := range handlers {
		// 异步处理事件（避免阻塞）
		go func(h EventHandler) {
			ctx := context.Background()
			_ = h.HandleEvent(ctx, event)
		}(handler)
	}
	
	return nil
}

// EmitStateChanged 发射状态变化事件
func (e *BaseEventEmitter) EmitStateChanged(from, to LifecycleState, trigger string) error {
	return e.Emit(Event{
		Type:    EventTypeStateChanged,
		Payload: map[string]interface{}{
			"from":    from,
			"to":      to,
			"trigger": trigger,
		},
	})
}

// EmitCommandCompleted 发射命令完成事件
func (e *BaseEventEmitter) EmitCommandCompleted(commandID string, result interface{}) error {
	return e.Emit(Event{
		Type:    EventTypeCommandCompleted,
		Payload: map[string]interface{}{
			"command_id": commandID,
			"result":     result,
		},
	})
}

// EmitCommandFailed 发射命令失败事件
func (e *BaseEventEmitter) EmitCommandFailed(commandID string, err error) error {
	return e.Emit(Event{
		Type:    EventTypeCommandFailed,
		Payload: map[string]interface{}{
			"command_id": commandID,
			"error":      err.Error(),
		},
	})
}
