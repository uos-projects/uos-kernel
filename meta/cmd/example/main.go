package main

import (
	"context"
	"fmt"
	"log"

	"github.com/uos-projects/uos-kernel/actors"
	"github.com/uos-projects/uos-kernel/kernel"
)

func main() {
	ctx := context.Background()

	// 1. 创建Actor系统
	system := actors.NewSystem(ctx)
	defer system.Shutdown()

	// 2. 创建资源内核
	k := kernel.NewKernel(system)

	// 3. 加载类型系统定义
	if err := k.LoadTypeSystem("../../typesystem.yaml"); err != nil {
		log.Fatalf("Failed to load type system: %v", err)
	}

	// 4. 打开资源（如果不存在则自动创建）
	// 使用 O_CREAT 标志
	fd, err := k.Open("Breaker", "BREAKER_001", kernel.O_CREAT)
	if err != nil {
		log.Fatalf("Failed to open resource: %v", err)
	}
	defer k.Close(fd)

	fmt.Printf("Opened resource with descriptor: %d\n", fd)

	// [验证类型安全] 尝试用错误的类型打开同一个资源
	fmt.Println("\nAttempting to open existing resource with wrong type...")
	_, err = k.Open("SynchronousMachine", "BREAKER_001", 0)
	if err != nil {
		fmt.Printf("Expected error caught: %v\n", err)
	} else {
		log.Fatal("Error: Expected type mismatch error but got success")
	}

	// [验证属性修改]
	fmt.Println("\nAttempting to modify resource properties...")
	updateReq := &kernel.WriteRequest{
		Updates: map[string]interface{}{
			"ratedCurrent": 2000.0,
			"description":  "Updated description via Write",
		},
	}
	if err := k.Write(ctx, fd, updateReq); err != nil {
		log.Fatalf("Failed to write resource: %v", err)
	}
	fmt.Println("Properties updated via k.Write")

	// 重新读取验证
	updatedState, err := k.Read(ctx, fd)
	if err != nil {
		log.Fatalf("Failed to read resource: %v", err)
	}
	if val, ok := updatedState.Properties["ratedCurrent"]; ok && val == 2000.0 {
		fmt.Println("Verification Success: ratedCurrent updated to 2000.0")
	} else {
		fmt.Printf("Verification Failed: ratedCurrent is %v\n", val)
	}

	// [验证持久化同步]
	fmt.Println("\nAttempting to sync resource state (Snapshot)...")
	_, err = k.Ioctl(ctx, fd, int(kernel.CMD_SYNC), nil)
	if err != nil {
		log.Fatalf("Failed to sync resource: %v", err)
	}
	fmt.Println("Resource state synced successfully (Snapshot created)")

	// 5. 查询资源信息（Stat）
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
