package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/uos-projects/uos-kernel/actors"
	"github.com/uos-projects/uos-kernel/actors/capacities"
)

func main() {
	ctx := context.Background()
	system := actors.NewSystem(ctx)
	defer system.Shutdown()

	// 创建资源 Actor（对应 PowerSystemResource）
	// 示例1: BE-Line_2 (ACLineSegment) 具有 AccumulatorResetCapacity
	lineActor := actors.NewPowerSystemResourceActor(
		"BE-Line_2",
		"ACLineSegment",
		nil,
	)

	// 添加能力（基于 CIM 关联关系）
	factory := actors.NewCapacityFactory()

	// 假设从 CIM 数据中发现 BE-Line_2 关联了 AccumulatorReset
	accCapacity, err := factory.CreateCapacity("AccumulatorReset", "ACC_RESET_1")
	if err != nil {
		log.Fatal(err)
	}
	lineActor.AddCapacity(accCapacity)

	// 注册 Actor（会自动启动）
	if err := system.Register(lineActor); err != nil {
		log.Fatal(err)
	}

	// 示例2: BE-G4 (SynchronousMachine) 具有 RaiseLowerCommandCapacity
	genActor := actors.NewPowerSystemResourceActor(
		"BE-G4",
		"SynchronousMachine",
		nil,
	)

	rlCapacity, err := factory.CreateCapacity("RaiseLowerCommand", "CMD_RL_1")
	if err != nil {
		log.Fatal(err)
	}
	genActor.AddCapacity(rlCapacity)

	if err := system.Register(genActor); err != nil {
		log.Fatal(err)
	}

	// 示例3: CIRCB-1230991526 (Breaker) 具有 CommandCapacity
	breakerActor := actors.NewPowerSystemResourceActor(
		"CIRCB-1230991526",
		"Breaker",
		nil,
	)

	cmdCapacity, err := factory.CreateCapacity("Command", "CMD_1")
	if err != nil {
		log.Fatal(err)
	}
	breakerActor.AddCapacity(cmdCapacity)

	if err := system.Register(breakerActor); err != nil {
		log.Fatal(err)
	}

	// 发送消息到 Actor，会自动路由到对应的 Capacity
	fmt.Println("Sending messages to actors...")

	// 发送累加器复位消息
	system.Send("BE-Line_2", &capacities.AccumulatorResetMessage{Value: 100})

	// 发送升降命令消息
	system.Send("BE-G4", &capacities.RaiseLowerCommandMessage{Delta: 10.0})

	// 发送命令消息
	system.Send("CIRCB-1230991526", &capacities.CommandMessage{
		Command: "open",
		Value:   0,
	})

	// 等待消息处理
	time.Sleep(200 * time.Millisecond)

	// 检查能力
	fmt.Println("\nActor capabilities:")
	if lineActor.HasCapacity("AccumulatorResetCapacity") {
		fmt.Printf("  %s has AccumulatorResetCapacity\n", lineActor.ResourceID())
	}
	if genActor.HasCapacity("RaiseLowerCommandCapacity") {
		fmt.Printf("  %s has RaiseLowerCommandCapacity\n", genActor.ResourceID())
	}
	if breakerActor.HasCapacity("CommandCapacity") {
		fmt.Printf("  %s has CommandCapacity\n", breakerActor.ResourceID())
	}

	fmt.Printf("\nSystem has %d actors\n", system.Count())
	fmt.Println("Example completed!")
}
