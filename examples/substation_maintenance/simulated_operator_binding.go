package main

import (
	"context"
	"fmt"
	"time"

	"github.com/uos-projects/uos-kernel/actors"
)

// OperatorDeviceEvent 表示来自操作员侧的状态反馈事件
// 这是 Binding 从真实世界（或模拟操作员）感知到的结果
type OperatorDeviceEvent struct {
	Action   string // "task_started", "power_outage_completed", "maintenance_completed", "power_restored", "task_completed"
	TaskID   string
	Step     string // "power_outage", "maintenance", "power_restored"
	Result   string // "success", "failed"
	Operator string
}

// SimulatedOperatorBinding 模拟操作员绑定
// 代表真实操作员的行为执行，Binding 执行实际行为，Actor 只反映状态
//
// 设计原则：
// - Binding 执行实际行为（发送命令到设备、等待、执行检修等）
// - Actor 只反映状态变化（任务开始、步骤完成、任务完成等）
// - 通过 ExternalEventMessage 反馈操作员的状态变化给 Actor
type SimulatedOperatorBinding struct {
	*actors.BaseBinding

	// 系统引用（用于发送命令到设备）
	system *actors.System
}

// NewSimulatedOperatorBinding 创建模拟操作员绑定
func NewSimulatedOperatorBinding(resourceID string, system *actors.System) *SimulatedOperatorBinding {
	base := actors.NewBaseBinding(actors.BindingTypeHuman, resourceID)
	return &SimulatedOperatorBinding{
		BaseBinding: base,
		system:      system,
	}
}

// Start 启动绑定（这里不需要额外逻辑）
func (b *SimulatedOperatorBinding) Start(ctx context.Context) error {
	return nil
}

// Stop 停止绑定
func (b *SimulatedOperatorBinding) Stop() error {
	return b.BaseBinding.Stop()
}

// OnExternalEvent 外部事件入口（本示例未使用，保留以符合接口）
func (b *SimulatedOperatorBinding) OnExternalEvent(ctx context.Context, event interface{}) error {
	// 直接转发为 ExternalEventMessage
	msg := &actors.ExternalEventMessage{
		BindingType: actors.BindingTypeHuman,
		Event:       event,
	}
	return b.SendToActor(msg)
}

// ExecuteExternal 执行外部操作（接收来自 Actor 的命令）
// 在这个模拟中，Binding 执行实际的检修流程：
// 1. 接收 StartMaintenanceCommand
// 2. 执行检修流程（停电、检修、恢复供电）
// 3. 通过 ExternalEventMessage 反馈每个步骤的完成情况
func (b *SimulatedOperatorBinding) ExecuteExternal(ctx context.Context, command interface{}) error {
	// 只处理 StartMaintenanceCommand
	cmd, ok := command.(*StartMaintenanceCommand)
	if !ok {
		return fmt.Errorf("unsupported command type: %T", command)
	}

	// 异步执行检修流程
	go b.executeMaintenanceFlow(ctx, cmd)

	return nil
}

// executeMaintenanceFlow 执行检修流程（Binding 执行实际行为）
func (b *SimulatedOperatorBinding) executeMaintenanceFlow(ctx context.Context, cmd *StartMaintenanceCommand) {
	fmt.Printf("\n========== [操作员 Binding] 开始执行检修流程 ==========\n")
	fmt.Printf("任务ID：%s\n", cmd.TaskID)
	fmt.Printf("操作时间：%s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("涉及设备：%d 个\n", len(cmd.Devices))
	fmt.Println()

	// 步骤1：反馈任务开始
	b.emitOperatorEvent(&OperatorDeviceEvent{
		Action:   "task_started",
		TaskID:   cmd.TaskID,
		Operator: cmd.OperatorID,
	})

	// 步骤2：执行停电操作
	if err := b.executePowerOutage(ctx, cmd); err != nil {
		b.emitOperatorEvent(&OperatorDeviceEvent{
			Action:   "task_completed",
			TaskID:   cmd.TaskID,
			Result:   "failed",
			Operator: cmd.OperatorID,
		})
		return
	}

	// 步骤3：执行检修操作
	if err := b.executeMaintenance(ctx, cmd); err != nil {
		b.emitOperatorEvent(&OperatorDeviceEvent{
			Action:   "task_completed",
			TaskID:   cmd.TaskID,
			Result:   "failed",
			Operator: cmd.OperatorID,
		})
		return
	}

	// 步骤4：恢复供电
	if err := b.executePowerRestore(ctx, cmd); err != nil {
		b.emitOperatorEvent(&OperatorDeviceEvent{
			Action:   "task_completed",
			TaskID:   cmd.TaskID,
			Result:   "failed",
			Operator: cmd.OperatorID,
		})
		return
	}

	// 步骤5：通知设备完成检修
	b.notifyDevicesMaintenanceCompleted(cmd)

	// 步骤6：反馈任务完成
	b.emitOperatorEvent(&OperatorDeviceEvent{
		Action:   "task_completed",
		TaskID:   cmd.TaskID,
		Result:   "success",
		Operator: cmd.OperatorID,
	})

	fmt.Println("\n========== [操作员 Binding] 检修流程执行完成 ==========")
}

// executePowerOutage 执行停电操作
func (b *SimulatedOperatorBinding) executePowerOutage(ctx context.Context, cmd *StartMaintenanceCommand) error {
	fmt.Println("【Binding 执行】停电操作")
	for i, deviceID := range cmd.Devices {
		fmt.Printf("  步骤 %d/%d: 打开断路器 %s\n", i+1, len(cmd.Devices), deviceID)

		openCmd := &OpenBreakerCommand{
			commandID: fmt.Sprintf("%s_open_%d", cmd.TaskID, i),
			Reason:    fmt.Sprintf("检修操作：%s", cmd.Reason),
			Operator:  cmd.OperatorID,
			TaskID:    cmd.TaskID,
		}

		// 发送命令到设备 Actor
		if err := b.system.Send(deviceID, openCmd); err != nil {
			return fmt.Errorf("发送打开命令到 %s 失败: %w", deviceID, err)
		}

		// 等待操作完成（实际应用中应该通过事件或响应消息）
		time.Sleep(300 * time.Millisecond)

		fmt.Printf("  ✓ 步骤 %d 完成\n", i+1)
	}

	// 反馈停电完成
	b.emitOperatorEvent(&OperatorDeviceEvent{
		Action:   "power_outage_completed",
		TaskID:   cmd.TaskID,
		Step:     "power_outage",
		Operator: cmd.OperatorID,
	})

	return nil
}

// executeMaintenance 执行检修操作
func (b *SimulatedOperatorBinding) executeMaintenance(ctx context.Context, cmd *StartMaintenanceCommand) error {
	fmt.Println("\n【Binding 执行】检修操作")
	// 模拟检修操作
	time.Sleep(1 * time.Second)
	fmt.Println("  ✓ 检修操作完成")

	// 反馈检修完成
	b.emitOperatorEvent(&OperatorDeviceEvent{
		Action:   "maintenance_completed",
		TaskID:   cmd.TaskID,
		Step:     "maintenance",
		Operator: cmd.OperatorID,
	})

	return nil
}

// executePowerRestore 执行恢复供电操作
func (b *SimulatedOperatorBinding) executePowerRestore(ctx context.Context, cmd *StartMaintenanceCommand) error {
	fmt.Println("\n【Binding 执行】恢复供电")
	for i, deviceID := range cmd.Devices {
		fmt.Printf("  步骤 %d/%d: 关闭断路器 %s\n", i+1, len(cmd.Devices), deviceID)

		closeCmd := &CloseBreakerCommand{
			commandID: fmt.Sprintf("%s_close_%d", cmd.TaskID, i),
			Reason:    "检修完成，恢复供电",
			Operator:  cmd.OperatorID,
			TaskID:    cmd.TaskID,
		}

		// 发送命令到设备 Actor
		if err := b.system.Send(deviceID, closeCmd); err != nil {
			return fmt.Errorf("发送关闭命令到 %s 失败: %w", deviceID, err)
		}

		// 等待操作完成
		time.Sleep(300 * time.Millisecond)

		fmt.Printf("  ✓ 步骤 %d 完成\n", i+1)
	}

	fmt.Println("\n========== 停电检修操作完成 ==========")

	// 反馈恢复供电完成
	b.emitOperatorEvent(&OperatorDeviceEvent{
		Action:   "power_restored",
		TaskID:   cmd.TaskID,
		Step:     "power_restored",
		Operator: cmd.OperatorID,
	})

	return nil
}

// notifyDevicesMaintenanceCompleted 通知设备完成检修
func (b *SimulatedOperatorBinding) notifyDevicesMaintenanceCompleted(cmd *StartMaintenanceCommand) {
	for _, deviceID := range cmd.Devices {
		completeCmd := &CompleteMaintenanceCommand{
			TaskID:   cmd.TaskID,
			Operator: cmd.OperatorID,
			Result:   "success",
		}
		if err := b.system.Send(deviceID, completeCmd); err != nil {
			fmt.Printf("[操作员 Binding] ⚠️  发送完成检修命令到 %s 失败: %v\n", deviceID, err)
		}
	}
}

// emitOperatorEvent 发射操作员事件（反馈给 Actor）
func (b *SimulatedOperatorBinding) emitOperatorEvent(event *OperatorDeviceEvent) {
	_ = b.OnExternalEvent(context.Background(), event)
}
