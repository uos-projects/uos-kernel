package main

import (
	"context"
	"fmt"
	"time"

	"github.com/uos-projects/uos-kernel/actor"
)

func main() {
	ctx := context.Background()

	// 1. 创建 Actor 系统
	system := actor.NewSystem(ctx)
	defer system.Shutdown()

	fmt.Println("=== 变电站停电检修操作场景演示 ===\n")
	fmt.Println("场景说明：")
	fmt.Println("1. 设备 Actor 持续运行，监测状态（温度、运行时间）")
	fmt.Println("2. 定期检修计划：运行时间达到阈值时触发检修")
	fmt.Println("3. 异常检测：温度异常时触发紧急检修")
	fmt.Println("4. 调度中心接收事件，制定检修计划并分配给操作员")
	fmt.Println("5. 调度操作员执行停电检修操作\n")

	// 2. 创建设备 Actor（长期存在的运行时实体）
	fmt.Println("【初始化】创建设备 Actor...")

	// 主变压器进线断路器
	breaker1 := NewBreakerActor("BREAKER-001", "主变压器进线断路器")
	system.Register(breaker1)

	// 主变压器出线断路器
	breaker2 := NewBreakerActor("BREAKER-002", "主变压器出线断路器")
	system.Register(breaker2)

	// 备用线路断路器
	breaker3 := NewBreakerActor("BREAKER-003", "备用线路断路器")
	system.Register(breaker3)

	// 3. 创建调度中心 Actor
	fmt.Println("【初始化】创建调度中心 Actor...")
	dispatcher := NewDispatcherActor(system)
	system.Register(dispatcher)

	// 4. 创建调度操作员 Actor
	fmt.Println("【初始化】创建调度操作员 Actor...")
	operator := NewDispatcherOperatorActor("OP-001", "张三", system)
	system.Register(operator)

	// 创建并绑定操作员 Binding（模拟操作员行为）
	operatorBinding := NewSimulatedOperatorBinding("OP-001", system)
	if err := operator.AddBinding(operatorBinding); err != nil {
		panic(fmt.Errorf("绑定操作员 Binding 失败: %w", err))
	}
	if err := operatorBinding.Start(ctx); err != nil {
		panic(fmt.Errorf("启动操作员 Binding 失败: %w", err))
	}

	// 注册操作员到调度中心
	dispatcher.RegisterOperator("OP-001")

	// 5. 启动所有 Actor（必须在注册事件处理器之前启动，确保 eventEmitter 已初始化）
	fmt.Println("\n【启动】启动所有 Actor...")
	if err := breaker1.Start(ctx); err != nil {
		panic(err)
	}
	if err := breaker2.Start(ctx); err != nil {
		panic(err)
	}
	if err := breaker3.Start(ctx); err != nil {
		panic(err)
	}
	if err := dispatcher.Start(ctx); err != nil {
		panic(err)
	}
	if err := operator.Start(ctx); err != nil {
		panic(err)
	}

	// 6. 设置事件订阅（调度中心订阅接收所有事件）
	// 调度中心 Actor 订阅事件，所有事件会直接发送给它
	if err := system.SubscribeEvent("DISPATCHER"); err != nil {
		panic(fmt.Errorf("订阅事件失败: %w", err))
	}
	fmt.Printf("[事件订阅] ✅ 调度中心已订阅事件\n")

	// 7. 显示初始状态
	fmt.Println("=== 初始状态 ===")
	displayDeviceStatus(breaker1, breaker2, breaker3)

	// 8. 等待设备监测触发事件（定期检修或异常检测）
	fmt.Println("\n【运行中】设备 Actor 持续监测状态...")
	fmt.Println("（等待触发检修事件：定期检修计划或异常检测）\n")

	// 模拟运行一段时间，让设备监测触发事件
	time.Sleep(3 * time.Second)

	// 9. 手动触发一个异常事件（演示异常检测）
	fmt.Println("【演示】手动触发温度异常事件...")
	simulateTemperatureAbnormal(system, breaker1)

	// 等待任务执行
	time.Sleep(5 * time.Second)

	// 10. 显示最终状态
	fmt.Println("\n=== 最终状态 ===")
	displayDeviceStatus(breaker1, breaker2, breaker3)

	// 11. 显示任务列表
	fmt.Println("\n=== 调度中心任务列表 ===")
	tasks := dispatcher.GetPendingTasks()
	for _, task := range tasks {
		fmt.Printf("任务ID：%s, 类型：%s, 状态：%s, 操作员：%s\n",
			task.TaskID, task.Type, task.Status, task.AssignedTo)
	}

	fmt.Println("\n=== 演示完成 ===")
	fmt.Println("（实际应用中，Actor 会一直运行，持续监测和响应事件）")
}

// simulateTemperatureAbnormal 模拟温度异常
func simulateTemperatureAbnormal(system *actor.System, breaker *BreakerActor) {
	// 模拟温度突然升高
	event := &DeviceAbnormalEvent{
		DeviceID:  breaker.ResourceID(),
		EventType: "temperature_high",
		Severity:  "critical",
		Details: map[string]interface{}{
			"temperature": 85.0,
			"threshold":   80.0,
		},
		Timestamp: time.Now(),
	}

	// 发送事件到调度中心
	_ = system.Send("DISPATCHER", event)
}

// displayDeviceStatus 显示设备状态
func displayDeviceStatus(breakers ...*BreakerActor) {
	for _, breaker := range breakers {
		status := breaker.GetStatus()
		fmt.Printf("%s (%s):\n", status["id"], status["name"])
		fmt.Printf("  状态：%s\n", map[bool]string{true: "打开", false: "关闭"}[status["isOpen"].(bool)])
		fmt.Printf("  电压：%.2f kV\n", status["voltage"].(float64))
		fmt.Printf("  电流：%.2f A\n", status["current"].(float64))
		fmt.Printf("  温度：%.2f °C\n", status["temperature"].(float64))
		fmt.Printf("  运行小时数：%d 小时\n", status["operationHours"].(int64))
		fmt.Printf("  上次检修：%s\n", status["lastMaintenanceTime"].(time.Time).Format("2006-01-02 15:04:05"))
		fmt.Println()
	}
}
