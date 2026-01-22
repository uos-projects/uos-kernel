package mq

import (
	"context"
	"fmt"
	"sync"
)

// MockMQConsumer 模拟 MQ 消费者（用于测试和开发）
type MockMQConsumer struct {
	mu          sync.RWMutex
	subscribers map[string][]chan *MQMessage
	closed      bool
}

// NewMockMQConsumer 创建新的模拟 MQ 消费者
func NewMockMQConsumer() *MockMQConsumer {
	return &MockMQConsumer{
		subscribers: make(map[string][]chan *MQMessage),
	}
}

// Consume 消费一条消息（阻塞，直到有消息或 context 取消）
func (m *MockMQConsumer) Consume(ctx context.Context, topic string) (*MQMessage, error) {
	ch, err := m.Subscribe(ctx, topic)
	if err != nil {
		return nil, err
	}

	select {
	case msg := <-ch:
		return msg, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Subscribe 订阅主题，通过 channel 返回消息
func (m *MockMQConsumer) Subscribe(ctx context.Context, topic string) (<-chan *MQMessage, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil, fmt.Errorf("consumer is closed")
	}

	ch := make(chan *MQMessage, 10) // 缓冲 10 条消息
	m.subscribers[topic] = append(m.subscribers[topic], ch)

	// 启动 goroutine 清理订阅（当 context 取消时）
	go func() {
		<-ctx.Done()
		m.mu.Lock()
		defer m.mu.Unlock()
		// 从订阅列表中移除
		subs := m.subscribers[topic]
		for i, sub := range subs {
			if sub == ch {
				m.subscribers[topic] = append(subs[:i], subs[i+1:]...)
				close(ch)
				break
			}
		}
	}()

	return ch, nil
}

// Close 关闭消费者
func (m *MockMQConsumer) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.closed = true
	for _, subs := range m.subscribers {
		for _, ch := range subs {
			close(ch)
		}
	}
	m.subscribers = make(map[string][]chan *MQMessage)
	return nil
}

// Publish 发布消息（用于测试，模拟外部系统发布）
func (m *MockMQConsumer) Publish(topic string, msg *MQMessage) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return fmt.Errorf("consumer is closed")
	}

	subs, exists := m.subscribers[topic]
	if !exists {
		return nil // 没有订阅者，忽略
	}

	// 向所有订阅者发送消息
	for _, ch := range subs {
		select {
		case ch <- msg:
		default:
			// channel 已满，跳过
		}
	}

	return nil
}
