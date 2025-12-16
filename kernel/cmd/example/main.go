package main

import (
	"context"
	"fmt"
	"log"

	"github.com/uos-projects/uos-kernel/actors"
	"github.com/uos-projects/uos-kernel/resource"
)

func main() {
	ctx := context.Background()

	// 1. 创建Actor系统
	system := actors.NewSystem(ctx)
	defer system.Shutdown()

	// 2. 创建资源Actor
	actor := actors.NewPowerSystemResourceActor("BREAKER_001", "Breaker", nil)
	system.Register(actor)

	// 3. 创建资源内核
	k := resource.NewResourceKernel(system)

	// 4. 加载类型系统定义
	if err := k.LoadTypeSystem("../../kernel/typesystem.yaml"); err != nil {
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
	fmt.Printf("  Capabilities: %v\n", state.Capabilities)

	// 8. 执行控制操作（Ioctl）
	fmt.Printf("\nExecuting control operations...\n")

	// 示例：SetPoint操作
	setPointArg := map[string]interface{}{"value": 150.0}
	result, err := k.Ioctl(ctx, fd, 0x1004, setPointArg) // CMD_SET_POINT
	if err != nil {
		log.Printf("Failed to execute SetPoint: %v", err)
	} else {
		fmt.Printf("SetPoint result: %v\n", result)
	}

	// 示例：查询资源信息
	info, err := k.Ioctl(ctx, fd, 0x1000, nil) // CMD_GET_RESOURCE_INFO
	if err != nil {
		log.Printf("Failed to get resource info: %v", err)
	} else {
		fmt.Printf("Resource info: %v\n", info)
	}

	fmt.Println("\nExample completed successfully!")
}
