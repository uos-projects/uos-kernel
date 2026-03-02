package main

import (
	"context"
	"fmt"
	"time"

	"github.com/uos-projects/uos-kernel/app"
	"github.com/uos-projects/uos-kernel/kernel"
)

// MaintenanceProcess 变电站停电检修流程
// 实现 app.Process 接口，业务逻辑全部在应用层
type MaintenanceProcess struct {
	BreakerIDs []string // 需要监控的断路器 ID 列表
}

func (p *MaintenanceProcess) Name() string { return "变电站停电检修" }

// Run 运行流程：监听设备异常事件，触发检修工作流
func (p *MaintenanceProcess) Run(ctx context.Context, k *kernel.Kernel) error {
	// 打开所有断路器资源
	fds := make(map[string]kernel.ResourceDescriptor)
	for _, id := range p.BreakerIDs {
		fd, err := k.Open("Breaker", id, 0)
		if err != nil {
			return fmt.Errorf("打开断路器 %s 失败: %w", id, err)
		}
		defer k.Close(fd)
		fds[id] = fd
	}

	// 监听第一个断路器的事件（演示用）
	primaryID := p.BreakerIDs[0]
	primaryFD := fds[primaryID]

	ch, cancel, err := k.Watch(primaryFD, nil) // nil = 接收所有事件
	if err != nil {
		return fmt.Errorf("Watch 失败: %w", err)
	}
	defer cancel()

	fmt.Printf("[%s] 开始监听设备事件...\n", p.Name())

	for {
		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-ch:
			if !ok {
				return nil
			}
			// 只关注业务事件（设备异常）
			if event.Type != kernel.EventCustom {
				continue
			}
			if _, ok := event.Data.(*DeviceAbnormalEvent); ok {
				fmt.Printf("\n[%s] 收到设备异常事件，启动检修工作流...\n", p.Name())
				if err := p.executeWorkflow(ctx, k, fds); err != nil {
					fmt.Printf("[%s] 工作流执行失败: %v\n", p.Name(), err)
				}
				return nil // 演示模式：完成一次工作流后退出
			}
		}
	}
}

// executeWorkflow 执行停电检修工作流
func (p *MaintenanceProcess) executeWorkflow(ctx context.Context, k *kernel.Kernel, fds map[string]kernel.ResourceDescriptor) error {
	fmt.Println("\n========== 开始停电检修流程 ==========")

	// 步骤 1：停电 — 打开所有断路器
	fmt.Println("\n【步骤 1】停电操作")
	for id, fd := range fds {
		fmt.Printf("  打开断路器 %s...\n", id)
		_, err := k.Ioctl(ctx, fd, 0x1003, map[string]interface{}{
			"message": &OpenBreakerCommand{
				Reason:   "停电检修",
				Operator: "maintenance-process",
			},
		})
		if err != nil {
			return fmt.Errorf("打开断路器 %s 失败: %w", id, err)
		}

		// 等待断路器确认打开
		if err := app.WaitForProperty(ctx, k, fd, "isOpen", true); err != nil {
			return fmt.Errorf("等待断路器 %s 打开超时: %w", id, err)
		}
		fmt.Printf("  ✓ 断路器 %s 已打开\n", id)
	}

	// 步骤 2：执行检修
	fmt.Println("\n【步骤 2】执行检修操作")
	fmt.Println("  检修中...")
	time.Sleep(1 * time.Second) // 模拟检修耗时
	fmt.Println("  ✓ 检修完成")

	// 步骤 3：复电 — 关闭所有断路器
	fmt.Println("\n【步骤 3】恢复供电")
	for id, fd := range fds {
		fmt.Printf("  关闭断路器 %s...\n", id)
		_, err := k.Ioctl(ctx, fd, 0x1003, map[string]interface{}{
			"message": &CloseBreakerCommand{
				Reason:   "检修完成，恢复供电",
				Operator: "maintenance-process",
			},
		})
		if err != nil {
			return fmt.Errorf("关闭断路器 %s 失败: %w", id, err)
		}

		// 等待断路器确认关闭
		if err := app.WaitForProperty(ctx, k, fd, "isOpen", false); err != nil {
			return fmt.Errorf("等待断路器 %s 关闭超时: %w", id, err)
		}
		fmt.Printf("  ✓ 断路器 %s 已关闭\n", id)
	}

	// 步骤 4：完成检修记录
	fmt.Println("\n【步骤 4】更新检修记录")
	for id, fd := range fds {
		_, err := k.Ioctl(ctx, fd, 0x1003, map[string]interface{}{
			"message": &CompleteMaintenanceCommand{
				Operator: "maintenance-process",
				Result:   "success",
			},
		})
		if err != nil {
			fmt.Printf("  ⚠ 更新 %s 检修记录失败: %v\n", id, err)
		}
	}

	// 等待异步命令处理完成
	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n========== 停电检修流程完成 ==========")
	return nil
}
