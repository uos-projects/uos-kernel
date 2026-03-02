package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uos-projects/uos-kernel/actor"
	"github.com/uos-projects/uos-kernel/kernel"
	"github.com/uos-projects/uos-kernel/meta"
	"github.com/uos-projects/uos-kernel/server"
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

	fmt.Println("=== 变电站停电检修 gRPC 服务端 ===")
	fmt.Println()

	// ========== 1. 创建 Actor System 和 Kernel ==========
	system := actor.NewSystem(ctx)
	defer system.Shutdown()
	k := kernel.NewKernel(system)

	// ========== 2. 注册资源类型 ==========
	k.GetTypeRegistry().Register(&meta.TypeDescriptor{Name: "Breaker"})

	// ========== 3. 创建设备 Actor ==========
	fmt.Println("【初始化】创建设备 Actor...")
	breaker1 := NewBreakerActor("BREAKER-001", "主变压器进线断路器")
	system.Register(breaker1)
	breaker1.Start(ctx)
	fmt.Println("  设备 BREAKER-001 已启动")

	// ========== 4. 创建 MessageRegistry ==========
	registry := server.NewMessageRegistry()
	registerMessages(registry)
	fmt.Println("  消息类型已注册")

	// ========== 5. 启动 gRPC 服务 ==========
	addr := ":50051"
	grpcServer, lis, err := server.Serve(addr, k, registry)
	if err != nil {
		fmt.Printf("启动 gRPC 服务失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\ngRPC 服务已启动，监听 %s\n", addr)
	fmt.Println("等待客户端连接...")
	fmt.Println()
	fmt.Println("提示：运行 Python 客户端连接此服务")
	fmt.Println("  cd sdk/python && python examples/substation_maintenance.py")

	// ========== 6. 延迟触发温度异常（演示用） ==========
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("\n【演示】自动触发温度异常事件...")
		simulateTemperatureAbnormal(breaker1)
	}()

	// ========== 7. 等待关闭 ==========
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("gRPC 服务错误: %v\n", err)
	}

	fmt.Println("服务已关闭")
}

// registerMessages 注册命令工厂和事件序列化器
func registerMessages(reg *server.MessageRegistry) {
	// 命令工厂
	reg.RegisterCommand("OpenBreakerCommand", func(fields map[string]interface{}) (actor.Message, error) {
		cmd := &OpenBreakerCommand{}
		if v, ok := fields["reason"].(string); ok {
			cmd.Reason = v
		}
		if v, ok := fields["operator"].(string); ok {
			cmd.Operator = v
		}
		return cmd, nil
	})

	reg.RegisterCommand("CloseBreakerCommand", func(fields map[string]interface{}) (actor.Message, error) {
		cmd := &CloseBreakerCommand{}
		if v, ok := fields["reason"].(string); ok {
			cmd.Reason = v
		}
		if v, ok := fields["operator"].(string); ok {
			cmd.Operator = v
		}
		return cmd, nil
	})

	reg.RegisterCommand("CompleteMaintenanceCommand", func(fields map[string]interface{}) (actor.Message, error) {
		cmd := &CompleteMaintenanceCommand{}
		if v, ok := fields["operator"].(string); ok {
			cmd.Operator = v
		}
		if v, ok := fields["result"].(string); ok {
			cmd.Result = v
		}
		return cmd, nil
	})

	// 事件序列化器
	reg.RegisterEvent("*main.DeviceAbnormalEvent", func(data interface{}) (map[string]interface{}, error) {
		e := data.(*DeviceAbnormalEvent)
		return map[string]interface{}{
			"_type":      "DeviceAbnormalEvent",
			"device_id":  e.DeviceID,
			"event_type": e.EventType,
			"severity":   e.Severity,
			"details":    e.Details,
			"timestamp":  e.Timestamp.Format(time.RFC3339),
		}, nil
	})

	reg.RegisterEvent("*main.MaintenanceRequiredEvent", func(data interface{}) (map[string]interface{}, error) {
		e := data.(*MaintenanceRequiredEvent)
		return map[string]interface{}{
			"_type":                "MaintenanceRequiredEvent",
			"device_id":            e.DeviceID,
			"reason":               e.Reason,
			"last_maintenance_time": e.LastMaintenanceTime.Format(time.RFC3339),
			"operation_hours":      e.OperationHours,
			"details":              e.Details,
			"timestamp":            e.Timestamp.Format(time.RFC3339),
		}, nil
	})
}

// simulateTemperatureAbnormal 模拟温度异常
func simulateTemperatureAbnormal(breaker *BreakerActor) {
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
	if emitter := breaker.GetEventEmitter(); emitter != nil {
		_ = emitter.EmitEvent(event)
	}
}
