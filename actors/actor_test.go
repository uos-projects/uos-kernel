package actors

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// testBehavior 是测试用的 Actor 行为实现
type testBehavior struct {
	name string
}

func newTestBehavior(name string) *testBehavior {
	return &testBehavior{name: name}
}

func (b *testBehavior) Handle(ctx context.Context, msg Message) (bool, error) {
	switch m := msg.(type) {
	case string:
		fmt.Printf("[%s] Received message: %s\n", b.name, m)
	case int:
		fmt.Printf("[%s] Received number: %d\n", b.name, m)
	case *testPingMessage:
		fmt.Printf("[%s] Received ping from %s\n", b.name, m.From)
		return true, nil
	case *testStopMessage:
		fmt.Printf("[%s] Received stop signal\n", b.name)
		return false, nil
	default:
		fmt.Printf("[%s] Received unknown message: %v\n", b.name, msg)
	}
	return true, nil
}

// testPingMessage 是测试用的 ping 消息
type testPingMessage struct {
	From string
	Time time.Time
}

// testStopMessage 是测试用的停止消息
type testStopMessage struct{}

// NewExampleBehavior 创建示例行为（用于测试和向后兼容）
// 注意：这个函数保留是为了向后兼容，实际应该使用自己的实现
func NewExampleBehavior(name string) *ExampleBehavior {
	return &ExampleBehavior{name: name}
}

// ExampleBehavior 是一个示例 Actor 行为实现（保留用于向后兼容）
type ExampleBehavior struct {
	name string
}

func (b *ExampleBehavior) Handle(ctx context.Context, msg Message) (bool, error) {
	switch m := msg.(type) {
	case string:
		fmt.Printf("[%s] Received message: %s\n", b.name, m)
	case int:
		fmt.Printf("[%s] Received number: %d\n", b.name, m)
	case *PingMessage:
		fmt.Printf("[%s] Received ping from %s\n", b.name, m.From)
		return true, nil
	case *StopMessage:
		fmt.Printf("[%s] Received stop signal\n", b.name)
		return false, nil
	default:
		fmt.Printf("[%s] Received unknown message: %v\n", b.name, msg)
	}
	return true, nil
}

// PingMessage 是一个示例消息类型（保留用于向后兼容）
type PingMessage struct {
	From string
	Time time.Time
}

// StopMessage 是停止消息（保留用于向后兼容）
type StopMessage struct{}

func TestBaseActor(t *testing.T) {
	behavior := newTestBehavior("test-actor")
	actor := NewBaseActor("test-actor", behavior)

	if actor.ID() != "test-actor" {
		t.Errorf("Expected ID 'test-actor', got '%s'", actor.ID())
	}

	ctx := context.Background()
	if err := actor.Start(ctx); err != nil {
		t.Fatalf("Failed to start actor: %v", err)
	}

	// 发送消息
	actor.Send("Hello")
	actor.Send(42)

	// 等待处理
	time.Sleep(50 * time.Millisecond)

	// 停止
	if err := actor.Stop(); err != nil {
		t.Fatalf("Failed to stop actor: %v", err)
	}
}

func TestSystem(t *testing.T) {
	ctx := context.Background()
	system := NewSystem(ctx)

	// 创建 Actor
	actor1, err := system.Spawn("actor1", newTestBehavior("Actor1"))
	if err != nil {
		t.Fatalf("Failed to spawn actor1: %v", err)
	}

	if actor1 == nil {
		t.Fatal("Actor1 is nil")
	}

	// 测试重复创建
	_, err = system.Spawn("actor1", newTestBehavior("Actor1"))
	if err == nil {
		t.Error("Expected error when spawning duplicate actor")
	}

	// 获取 Actor
	retrieved, exists := system.Get("actor1")
	if !exists {
		t.Error("Actor1 should exist")
	}
	if retrieved != actor1 {
		t.Error("Retrieved actor should be the same instance")
	}

	// 发送消息
	if err := system.Send("actor1", "test message"); err != nil {
		t.Errorf("Failed to send message: %v", err)
	}

	// 等待处理
	time.Sleep(50 * time.Millisecond)

	// 检查数量
	if system.Count() != 1 {
		t.Errorf("Expected 1 actor, got %d", system.Count())
	}

	// 停止 Actor
	if err := system.Stop("actor1"); err != nil {
		t.Errorf("Failed to stop actor: %v", err)
	}

	// 检查数量
	if system.Count() != 0 {
		t.Errorf("Expected 0 actors after stop, got %d", system.Count())
	}

	// 关闭系统
	if err := system.Shutdown(); err != nil {
		t.Errorf("Failed to shutdown system: %v", err)
	}
}

func TestSystemShutdown(t *testing.T) {
	ctx := context.Background()
	system := NewSystem(ctx)

	// 创建多个 Actor
	for i := 0; i < 5; i++ {
		id := fmt.Sprintf("actor%d", i)
		_, err := system.Spawn(id, newTestBehavior(id))
		if err != nil {
			t.Fatalf("Failed to spawn %s: %v", id, err)
		}
	}

	if system.Count() != 5 {
		t.Errorf("Expected 5 actors, got %d", system.Count())
	}

	// 关闭系统
	if err := system.Shutdown(); err != nil {
		t.Errorf("Failed to shutdown system: %v", err)
	}

	if system.Count() != 0 {
		t.Errorf("Expected 0 actors after shutdown, got %d", system.Count())
	}
}
