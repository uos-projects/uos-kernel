package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uos-projects/uos-kernel/actor"
	"github.com/uos-projects/uos-kernel/app"
	"github.com/uos-projects/uos-kernel/kernel"
	"github.com/uos-projects/uos-kernel/meta"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n收到退出信号，正在关闭...")
		cancel()
	}()

	fmt.Println("=== 变电站停电检修操作场景演示 ===")
	fmt.Println()
	fmt.Println("架构说明：")
	fmt.Println("  Actor 层：BreakerActor（设备资源）持续监测状态，暴露开/合能力")
	fmt.Println("  Kernel 层：Open / Watch / Ioctl / Read 等 POSIX 风格 API")
	fmt.Println("  App 层：MaintenanceProcess（业务流程）通过 Kernel API 编排检修工作流")
	fmt.Println()

	// ========== 1. 创建 Actor System 和 Kernel ==========
	system := actor.NewSystem(ctx)
	defer system.Shutdown()
	k := kernel.NewKernel(system)

	// ========== 2. 注册资源类型 ==========
	k.GetTypeRegistry().Register(&meta.TypeDescriptor{Name: "Breaker"})

	// ========== 3. 创建设备 Actor（资源层） ==========
	fmt.Println("【初始化】创建设备 Actor...")
	breaker1 := NewBreakerActor("BREAKER-001", "主变压器进线断路器")
	system.Register(breaker1)
	breaker1.Start(ctx)

	// ========== 4. 显示初始状态 ==========
	fmt.Println("\n=== 初始状态 ===")
	displayDeviceStatus(breaker1)

	// ========== 5. 启动应用层业务流程（后台） ==========
	// 流程通过 Kernel API 操控设备，不直接接触 Actor
	application := app.New(k)

	// 创建带超时的上下文（演示用）
	processCtx, processCancel := context.WithTimeout(ctx, 30*time.Second)
	defer processCancel()

	// 在后台启动业务流程，先建立 Watch 监听
	done := make(chan error, 1)
	go func() {
		done <- application.Run(processCtx, &MaintenanceProcess{
			BreakerIDs: []string{"BREAKER-001"},
		})
	}()

	// 等待进程启动并建立 Watch 监听
	time.Sleep(500 * time.Millisecond)

	// ========== 6. 设备监测运行 ==========
	fmt.Println("\n【运行中】设备 Actor 持续监测状态...")
	fmt.Println("（等待触发检修事件：定期检修计划或异常检测）")
	time.Sleep(2 * time.Second)

	// ========== 7. 手动触发温度异常（演示） ==========
	fmt.Println("\n【演示】手动触发温度异常事件...")
	simulateTemperatureAbnormal(system, breaker1)

	// 等待业务流程完成
	if err := <-done; err != nil {
		fmt.Printf("业务流程错误: %v\n", err)
	}

	// ========== 8. 显示最终状态 ==========
	fmt.Println("\n=== 最终状态 ===")
	displayDeviceStatus(breaker1)

	fmt.Println("\n=== 演示完成 ===")
}

// simulateTemperatureAbnormal 模拟温度异常
func simulateTemperatureAbnormal(system *actor.System, breaker *BreakerActor) {
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
	// 通过 EventEmitter 发布事件，Watch 可以观察到
	if emitter := breaker.GetEventEmitter(); emitter != nil {
		_ = emitter.EmitEvent(event)
	}
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
