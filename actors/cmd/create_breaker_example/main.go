package main

import (
	"context"
	"fmt"
	"log"

	"github.com/uos-projects/uos-kernel/actors"
	"github.com/uos-projects/uos-kernel/actors/capacities"
)

func main() {
	ctx := context.Background()
	
	// 创建 Actor System
	system := actors.NewSystem(ctx)
	defer system.Shutdown()
	
	// 创建 Breaker Actor（带 Control 的 PowerSystemResource）
	breakerURI := "http://www.iec.ch/TC57/CIM#Breaker"
	breakerActor := actors.NewCIMResourceActor(
		"breaker-001",
		breakerURI,
		nil, // behavior
	)
	
	// 设置 Breaker 的属性（包括继承的属性）
	breakerActor.SetProperty("mRID", "breaker-001")
	breakerActor.SetProperty("name", "Main Breaker")
	breakerActor.SetProperty("normalOpen", false)
	breakerActor.SetProperty("open", false)
	breakerActor.SetProperty("locked", false)
	breakerActor.SetProperty("ratedCurrent", 1000.0)
	breakerActor.SetProperty("description", "Main circuit breaker")
	
	// 添加 Control Capacity（Command Capacity 用于控制开关）
	commandCapacity := capacities.NewCommandCapacity("breaker-command-001")
	breakerActor.AddCapacity(commandCapacity)
	
	// 注册到 System
	if err := system.Register(breakerActor); err != nil {
		log.Fatalf("Failed to register breaker actor: %v", err)
	}
	
	// 启动 Actor
	if err := breakerActor.Start(ctx); err != nil {
		log.Fatalf("Failed to start breaker actor: %v", err)
	}
	
	fmt.Printf("Breaker Actor created successfully:\n")
	fmt.Printf("  ID: %s\n", breakerActor.ID())
	fmt.Printf("  OWL Class URI: %s\n", breakerActor.GetOWLClassURI())
	fmt.Printf("  Resource Type: %s\n", breakerActor.ResourceType())
	
	// 显示所有属性
	props := breakerActor.GetAllProperties()
	fmt.Printf("\nProperties:\n")
	for k, v := range props {
		fmt.Printf("  %s: %v\n", k, v)
	}
	
	// 创建快照（保存到 Iceberg）
	snapshot, err := breakerActor.CreateSnapshot(1)
	if err != nil {
		log.Fatalf("Failed to create snapshot: %v", err)
	}
	
	fmt.Printf("\nSnapshot created:\n")
	fmt.Printf("  Actor ID: %s\n", snapshot.ActorID)
	fmt.Printf("  OWL Class URI: %s\n", snapshot.OWLClassURI)
	fmt.Printf("  Sequence: %d\n", snapshot.Sequence)
	fmt.Printf("  Timestamp: %s\n", snapshot.Timestamp)
	
	// 显示快照中的属性（已转换为 Iceberg 字段名）
	fmt.Printf("\nSnapshot properties (Iceberg field names):\n")
	for k, v := range snapshot.Properties {
		fmt.Printf("  %s: %v\n", k, v)
	}
}
