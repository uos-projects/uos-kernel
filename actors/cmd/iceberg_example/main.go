package main

import (
	"context"
	"fmt"
	"log"

	"github.com/uos-projects/uos-kernel/actors"
	"github.com/uos-projects/uos-kernel/actors/capacities"
	"github.com/uos-projects/uos-kernel/actors/state"
)

func main() {
	ctx := context.Background()

	// 1. 创建 Actor 系统
	system := actors.NewSystem(ctx)
	defer system.Shutdown()

	// 2. 创建 IcebergSnapshotStore（优先使用 Thrift 协议，失败则回退到 beeline）
	icebergStore, err := state.NewIcebergSnapshotStore(state.IcebergSnapshotStoreConfig{
		Host:      "localhost",
		Port:      10000, // Spark Thrift Server 二进制端口（Thrift 模式）
		Database:  "default",
		Namespace: "ontology.grid",
		UseThrift: true, // 优先使用 Thrift，失败会自动回退到 beeline
	})
	if err != nil {
		log.Fatalf("Failed to create IcebergSnapshotStore: %v", err)
	}

	// 3. 创建 CIMResourceActor
	breakerURI := "http://www.iec.ch/TC57/CIM#Breaker"
	actor := actors.NewCIMResourceActor("BREAKER_001", breakerURI, nil)

	// 设置属性
	actor.SetProperty("mRID", "BREAKER_001")
	actor.SetProperty("name", "Main Breaker")
	actor.SetProperty("normalOpen", false)
	actor.SetProperty("ratedCurrent", 1000.0)

	// 添加能力
	commandCapacity := capacities.NewCommandCapacity("breaker-command")
	actor.AddCapacity(commandCapacity)

	// 设置 Iceberg 快照存储
	actor.SetSnapshotStore(icebergStore)

	// 注册到系统
	if err := system.Register(actor); err != nil {
		log.Fatalf("Failed to register actor: %v", err)
	}

	fmt.Println("Actor created and registered successfully")

	// 4. 创建快照并保存到 Iceberg
	fmt.Println("\nCreating snapshot...")
	if err := actor.SaveSnapshotAsync(ctx); err != nil {
		log.Fatalf("Failed to save snapshot: %v", err)
	}

	// 等待异步保存完成（实际应该使用更好的同步机制）
	fmt.Println("Snapshot saved to Iceberg (async)")

	// 5. 修改属性
	fmt.Println("\nUpdating properties...")
	actor.SetProperty("ratedCurrent", 2000.0)
	actor.SetProperty("description", "Updated breaker")

	// 6. 再次保存快照
	fmt.Println("Creating second snapshot...")
	if err := actor.SaveSnapshotAsync(ctx); err != nil {
		log.Fatalf("Failed to save second snapshot: %v", err)
	}

	fmt.Println("\nExample completed!")
	fmt.Println("\nNote: This example requires:")
	fmt.Println("  1. Nessie Catalog running at http://localhost:19120")
	fmt.Println("  2. MinIO running at http://localhost:19000")
	fmt.Println("  3. PySpark installed and available in PATH")
	fmt.Println("  4. Iceberg Spark dependencies configured")
}
