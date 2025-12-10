package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/uos-projects/uos-kernel/actors"
)

// coordinatorBehavior 协调者 Actor，会向其他 Actor 发送消息
type coordinatorBehavior struct {
	name string
}

func (b *coordinatorBehavior) Handle(ctx context.Context, msg actors.Message) (bool, error) {
	switch m := msg.(type) {
	case *startWorkMessage:
		fmt.Printf("[%s] Received start work message, target: %s\n", b.name, m.TargetID)
		// 这里可以通过 System 发送消息，但需要 System 引用
		// 实际使用中，可以通过 GetRef 或 SendTo 方法
		return true, nil
	case string:
		fmt.Printf("[%s] Received: %s\n", b.name, m)
	default:
		fmt.Printf("[%s] Received: %v\n", b.name, msg)
	}
	return true, nil
}

// workerBehavior 工作者 Actor
type workerBehavior struct {
	name string
}

func (b *workerBehavior) Handle(ctx context.Context, msg actors.Message) (bool, error) {
	switch m := msg.(type) {
	case *workMessage:
		fmt.Printf("[%s] Processing work: %s\n", b.name, m.Content)
		// 完成后可以向协调者发送完成消息
		return true, nil
	case string:
		fmt.Printf("[%s] Received: %s\n", b.name, m)
	default:
		fmt.Printf("[%s] Received: %v\n", b.name, msg)
	}
	return true, nil
}

// startWorkMessage 启动工作消息
type startWorkMessage struct {
	TargetID string
	Work     string
}

// workMessage 工作消息
type workMessage struct {
	Content string
}

func main() {
	ctx := context.Background()
	system := actors.NewSystem(ctx)
	defer system.Shutdown()

	// 创建协调者 Actor
	coordinator, err := system.Spawn("coordinator", &coordinatorBehavior{name: "Coordinator"})
	if err != nil {
		log.Fatal(err)
	}

	// 创建工作者 Actor
	worker1, err := system.Spawn("worker1", &workerBehavior{name: "Worker1"})
	if err != nil {
		log.Fatal(err)
	}

	_, err = system.Spawn("worker2", &workerBehavior{name: "Worker2"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Actor 之间通信示例 ===")
	fmt.Println()

	// 方式1: 通过 System 发送（外部代码）
	fmt.Println("方式1: 通过 System 发送消息")
	system.Send("coordinator", "Start coordination")
	system.Send("worker1", &workMessage{Content: "Task 1"})
	system.Send("worker2", &workMessage{Content: "Task 2"})

	time.Sleep(100 * time.Millisecond)

	// 方式2: Actor 之间直接通信（通过 BaseActor 的方法）
	fmt.Println("\n方式2: Actor 之间直接通信")
	coordinatorBase := coordinator.(*actors.BaseActor)
	
	// 协调者向工作者发送消息
	if err := coordinatorBase.SendTo("worker1", &workMessage{Content: "Task from coordinator"}); err != nil {
		log.Printf("Error sending message: %v", err)
	}

	// 通过 ActorRef 发送
	worker1Ref, err := coordinatorBase.GetRef("worker1")
	if err != nil {
		log.Printf("Error getting ref: %v", err)
	} else {
		worker1Ref.Tell(&workMessage{Content: "Task via ActorRef"})
	}

	// Worker1 向 Worker2 发送消息
	worker1Base := worker1.(*actors.BaseActor)
	if err := worker1Base.SendTo("worker2", "Hello from worker1"); err != nil {
		log.Printf("Error sending message: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	fmt.Printf("\nSystem has %d actors\n", system.Count())
	fmt.Println("Example completed!")
}

