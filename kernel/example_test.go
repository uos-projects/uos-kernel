package kernel_test

import (
	"context"
	"fmt"

	"github.com/uos-projects/uos-kernel/actors"
	"github.com/uos-projects/uos-kernel/kernel"
)

func ExampleManager() {
	ctx := context.Background()

	// 创建 Actor 系统
	system := actors.NewSystem(ctx)
	defer system.Shutdown()

	// 创建 CIMResourceActor
	actor := actors.NewCIMResourceActor("BE-G4", "http://www.iec.ch/TC57/CIM#SynchronousMachine", nil)

	// 添加 Capabilities
	factory := actors.NewCapacityFactory()
	setPointCap, _ := factory.CreateCapacity("SetPoint", "SET_PNT_1")
	commandCap, _ := factory.CreateCapacity("Command", "CMD_1")

	actor.AddCapacity(setPointCap)
	actor.AddCapacity(commandCap)

	system.Register(actor)

	// 创建资源管理器
	rm := kernel.NewManager(system)

	// 1. 打开资源（每次Open都会返回新的描述符）
	fd1, err := rm.Open("BE-G4")
	if err != nil {
		panic(err)
	}
	defer rm.Close(fd1)

	// 多次打开同一个资源，会返回不同的描述符
	fd2, err := rm.Open("BE-G4")
	if err != nil {
		panic(err)
	}
	defer rm.Close(fd2)

	fmt.Printf("First descriptor: %d\n", fd1)
	fmt.Printf("Second descriptor: %d\n", fd2)

	// 使用第一个描述符进行操作
	fd := fd1

	// 2. 读取 Actor 状态
	state, err := rm.Read(ctx, fd)
	if err == nil {
		fmt.Printf("Resource ID: %s\n", state.ResourceID)
		fmt.Printf("Resource Type: %s\n", state.ResourceType)
		fmt.Printf("Capabilities: %v\n", state.Capabilities)
	}

	// 3. 使用 rctl 获取资源信息
	info, _ := rm.RCtl(ctx, fd, kernel.CMD_GET_RESOURCE_INFO, nil)
	if info != nil {
		fmt.Printf("Resource Info: %+v\n", info)
	}

	// 4. 使用 rctl 列出所有能力
	caps, _ := rm.RCtl(ctx, fd, kernel.CMD_LIST_CAPABILITIES, nil)
	fmt.Printf("Capabilities: %v\n", caps)

	// 5. 使用 rctl 发送 SetPoint 控制命令
	setPointArg := map[string]interface{}{
		"value": 150.0,
	}
	rm.RCtl(ctx, fd, kernel.CMD_SET_POINT, setPointArg)

	// 6. 使用 rctl 发送 Command 控制命令
	commandArg := map[string]interface{}{
		"command": "START",
		"value":   1,
	}
	rm.RCtl(ctx, fd, kernel.CMD_COMMAND, commandArg)

	// Output:
	// First descriptor: 0
	// Second descriptor: 1
	// Resource ID: BE-G4
	// Resource Type: SynchronousMachine
	// Capabilities: [SetPointCapacity CommandCapacity]
	// Resource Info: &{ResourceID:BE-G4 ResourceType:SynchronousMachine}
	// Capabilities: [SetPointCapacity CommandCapacity]
}

func ExampleManager_Exclusive() {
	ctx := context.Background()

	// 创建 Actor 系统
	system := actors.NewSystem(ctx)
	defer system.Shutdown()

	// 创建 CIMResourceActor
	actor := actors.NewCIMResourceActor("EXCLUSIVE-RESOURCE", "http://www.iec.ch/TC57/CIM#SynchronousMachine", nil)
	system.Register(actor)

	// 创建资源管理器
	rm := kernel.NewManager(system)

	// 1. 以排他性方式打开资源
	fd1, err := rm.OpenWithExclusive("EXCLUSIVE-RESOURCE", true)
	if err != nil {
		panic(err)
	}
	defer rm.Close(fd1)

	fmt.Printf("First exclusive descriptor: %d\n", fd1)

	// 2. 尝试再次打开排他性资源，应该失败
	fd2, err := rm.OpenWithExclusive("EXCLUSIVE-RESOURCE", true)
	if err != nil {
		fmt.Printf("Failed to open exclusive resource again: %v\n", err)
	} else {
		defer rm.Close(fd2)
		fmt.Printf("Second descriptor: %d\n", fd2)
	}

	// Output:
	// First exclusive descriptor: 0
	// Failed to open exclusive resource again: resource EXCLUSIVE-RESOURCE is exclusive and already opened
}
