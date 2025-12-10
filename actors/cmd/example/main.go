package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/uos-projects/uos-kernel/actors"
)

// simpleBehavior 是一个简单的 Actor 行为实现
type simpleBehavior struct {
	name string
}

func (b *simpleBehavior) Handle(ctx context.Context, msg actors.Message) (bool, error) {
	switch m := msg.(type) {
	case string:
		fmt.Printf("[%s] Received message: %s\n", b.name, m)
	case int:
		fmt.Printf("[%s] Received number: %d\n", b.name, m)
	default:
		fmt.Printf("[%s] Received unknown message: %v\n", b.name, msg)
	}
	return true, nil
}

// pingMessage 是一个示例 ping 消息
type pingMessage struct {
	From string
	Time time.Time
}

func main() {
	ctx := context.Background()
	system := actors.NewSystem(ctx)
	defer system.Shutdown()

	// 创建两个 Actor
	_, err := system.Spawn("printer", &simpleBehavior{name: "Printer"})
	if err != nil {
		log.Fatal(err)
	}

	_, err = system.Spawn("calculator", &simpleBehavior{name: "Calculator"})
	if err != nil {
		log.Fatal(err)
	}

	// 发送各种消息
	fmt.Println("Sending messages...")
	system.Send("printer", "Hello from main!")
	system.Send("printer", 42)
	system.Send("calculator", 10)
	system.Send("calculator", 20)

	// 发送 ping 消息
	system.Send("printer", &pingMessage{
		From: "main",
		Time: time.Now(),
	})

	// 等待消息处理
	time.Sleep(200 * time.Millisecond)

	fmt.Printf("\nSystem has %d actors\n", system.Count())

	// 停止 Actor
	system.Stop("printer")
	system.Stop("calculator")

	fmt.Println("Example completed!")
}
