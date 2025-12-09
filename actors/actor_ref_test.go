package actors

import (
	"context"
	"testing"
	"time"
)

func TestActorRef(t *testing.T) {
	ctx := context.Background()
	system := NewSystem(ctx)
	defer system.Shutdown()

	// 创建 Actor
	_, err := system.Spawn("actor1", newTestBehavior("Actor1"))
	if err != nil {
		t.Fatalf("Failed to spawn actor1: %v", err)
	}

	// 获取 ActorRef
	ref, err := system.GetRef("actor1")
	if err != nil {
		t.Fatalf("Failed to get ref: %v", err)
	}

	if ref.ID() != "actor1" {
		t.Errorf("Expected ID 'actor1', got '%s'", ref.ID())
	}

	// 通过 ActorRef 发送消息
	if err := ref.Send("Hello from ref"); err != nil {
		t.Errorf("Failed to send message: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
}

func TestActorToActorCommunication(t *testing.T) {
	ctx := context.Background()
	system := NewSystem(ctx)
	defer system.Shutdown()

	// 创建一个简单的 Behavior 用于测试
	behavior1 := newTestBehavior("Actor1")

	// 创建两个 Actor
	actor1, err := system.Spawn("actor1", behavior1)
	if err != nil {
		t.Fatalf("Failed to spawn actor1: %v", err)
	}

	_, err = system.Spawn("actor2", newTestBehavior("Actor2"))
	if err != nil {
		t.Fatalf("Failed to spawn actor2: %v", err)
	}

	// Actor1 向 Actor2 发送消息
	baseActor1 := actor1.(*BaseActor)
	if err := baseActor1.SendTo("actor2", "Hello from actor1"); err != nil {
		t.Errorf("Failed to send message: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	// Actor1 通过 GetRef 获取 Actor2 的引用并发送消息
	ref, err := baseActor1.GetRef("actor2")
	if err != nil {
		t.Fatalf("Failed to get ref: %v", err)
	}

	if err := ref.Send("Hello via ref"); err != nil {
		t.Errorf("Failed to send via ref: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
}
