package actor

import (
	"context"
	"fmt"
	"sync"
)

// ============================================================================
// 事件发射器（Event Emission）
// ============================================================================

// EventEmitter 事件发射器接口
// 用于发射 Event（Message）类型的事件
type EventEmitter interface {
	// EmitEvent 发射事件（Event/Message）
	EmitEvent(event Event) error
}

// BaseEventEmitter 基础事件发射器实现
type BaseEventEmitter struct {
	actorID string
	system  *System // 用于发布事件到 System
	mu      sync.RWMutex
}

// NewBaseEventEmitter 创建基础事件发射器
func NewBaseEventEmitter(actorID string, system *System) *BaseEventEmitter {
	return &BaseEventEmitter{
		actorID: actorID,
		system:  system,
	}
}

// EmitEvent 发射事件（Event/Message）
// 事件会被发布到 System 的事件分发机制
func (e *BaseEventEmitter) EmitEvent(event Event) error {
	if e == nil {
		return fmt.Errorf("event emitter is nil")
	}
	if e.system == nil {
		return fmt.Errorf("system is nil, cannot emit event")
	}

	// 通过 System 发布事件
	ctx := context.Background()
	return e.system.PublishEvent(ctx, event)
}
