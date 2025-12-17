package main

import (
	"context"
	"fmt"
	"log"

	"github.com/uos-projects/uos-kernel/actors"
	"github.com/uos-projects/uos-kernel/actors/capacities"
	"github.com/uos-projects/uos-kernel/resource"
)

func main() {
	ctx := context.Background()

	// 1. 创建Actor系统
	system := actors.NewSystem(ctx)
	defer system.Shutdown()

	// 2. 创建 CIMResourceActor（带属性和 OWL 引用）
	breakerURI := "http://www.iec.ch/TC57/CIM#Breaker"
	actor := actors.NewCIMResourceActor("BREAKER_001", breakerURI, nil)

	// 设置 Breaker 的属性（包括继承的属性）
	actor.SetProperty("mRID", "BREAKER_001")
	actor.SetProperty("name", "Main Breaker")
	actor.SetProperty("description", "Main circuit breaker for substation A")
	actor.SetProperty("normalOpen", false)
	actor.SetProperty("open", false)
	actor.SetProperty("locked", false)
	actor.SetProperty("ratedCurrent", 1000.0)

	// 添加 Command Capacity
	commandCapacity := capacities.NewCommandCapacity("breaker-command")
	actor.AddCapacity(commandCapacity)

	// 注册到 System
	if err := system.Register(actor); err != nil {
		log.Fatalf("Failed to register actor: %v", err)
	}

	// 3. 创建资源内核
	k := resource.NewResourceKernel(system)

	// 4. 加载类型系统定义
	if err := k.LoadTypeSystem("../../typesystem.yaml"); err != nil {
		log.Fatalf("Failed to load type system: %v", err)
	}

	// 5. 打开资源
	fd, err := k.Open("Breaker", "BREAKER_001", 0)
	if err != nil {
		log.Fatalf("Failed to open resource: %v", err)
	}
	defer k.Close(fd)

	fmt.Printf("Opened resource with descriptor: %d\n", fd)

	// 6. 查询资源信息（Stat）
	stat, err := k.Stat(ctx, fd)
	if err != nil {
		log.Fatalf("Failed to stat resource: %v", err)
	}

	fmt.Printf("\nResource Stat:\n")
	fmt.Printf("  ResourceID: %s\n", stat.ResourceID)
	fmt.Printf("  ResourceType: %s\n", stat.ResourceType)
	if stat.TypeDescriptor != nil {
		fmt.Printf("  Type Name: %s\n", stat.TypeDescriptor.Name)
		fmt.Printf("  Base Type: %v\n", stat.TypeDescriptor.BaseType)
		fmt.Printf("  Attributes: %d\n", len(stat.TypeDescriptor.GetAllAttributes()))
		fmt.Printf("  Capabilities: %d\n", len(stat.TypeDescriptor.GetAllCapabilities()))
	}
	fmt.Printf("  Capabilities:\n")
	for _, cap := range stat.Capabilities {
		fmt.Printf("    - %s: %v\n", cap.Name, cap.Operations)
	}

	// 7. 读取资源状态（Read）
	state, err := k.Read(ctx, fd)
	if err != nil {
		log.Fatalf("Failed to read resource: %v", err)
	}

	fmt.Printf("\nResource State:\n")
	fmt.Printf("  ResourceID: %s\n", state.ResourceID)
	fmt.Printf("  ResourceType: %s\n", state.ResourceType)
	fmt.Printf("  OWLClassURI: %s\n", state.OWLClassURI)
	fmt.Printf("  Capabilities: %v\n", state.Capabilities)

	// 显示属性
	fmt.Printf("\nProperties:\n")
	for k, v := range state.Properties {
		fmt.Printf("  %s: %v\n", k, v)
	}

	fmt.Println("\nExample completed successfully!")
}
