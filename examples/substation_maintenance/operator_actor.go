package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/uos-projects/uos-kernel/actors"
)

// DispatcherOperatorActor 调度操作员 Actor
// 接收检修任务，执行停电检修操作
type DispatcherOperatorActor struct {
	*actors.BaseResourceActor

	operatorID   string
	operatorName string
	currentTask  *MaintenanceTask
	taskMu       sync.RWMutex

	// 系统引用（用于发送消息）
	system *actors.System
}

// NewDispatcherOperatorActor 创建调度操作员 Actor
func NewDispatcherOperatorActor(id string, name string, system *actors.System) *DispatcherOperatorActor {
	return &DispatcherOperatorActor{
		BaseResourceActor: actors.NewBaseResourceActor(id, "DispatcherOperator", nil),
		operatorID:        id,
		operatorName:      name,
		system:            system,
	}
}

// Receive 重写消息处理逻辑
func (o *DispatcherOperatorActor) Receive(ctx context.Context, msg actors.Message) error {
	// 处理检修任务
	switch task := msg.(type) {
	case *MaintenanceTask:
		return o.handleMaintenanceTask(ctx, task)
	}

	// 其他消息交给基类处理
	return o.BaseResourceActor.Receive(ctx, msg)
}

// handleMaintenanceTask 处理检修任务
func (o *DispatcherOperatorActor) handleMaintenanceTask(ctx context.Context, task *MaintenanceTask) error {
	o.taskMu.Lock()
	o.currentTask = task
	o.taskMu.Unlock()

	fmt.Printf("\n[操作员 %s] 📋 收到检修任务：\n", o.operatorName)
	fmt.Printf("  任务ID：%s\n", task.TaskID)
	fmt.Printf("  类型：%s\n", task.Type)
	fmt.Printf("  设备：%v\n", task.Devices)
	fmt.Printf("  原因：%s\n", task.Reason)

	// 执行检修操作
	return o.executeMaintenanceOperation(ctx, task)
}

// executeMaintenanceOperation 执行停电检修操作
func (o *DispatcherOperatorActor) executeMaintenanceOperation(
	ctx context.Context,
	task *MaintenanceTask,
) error {
	fmt.Printf("\n========== 开始停电检修操作 ==========\n")
	fmt.Printf("操作员：%s (%s)\n", o.operatorName, o.operatorID)
	fmt.Printf("任务ID：%s\n", task.TaskID)
	fmt.Printf("操作时间：%s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("涉及设备：%d 个\n", len(task.Devices))
	fmt.Println()

	// 步骤1：停电操作（按顺序打开所有断路器）
	fmt.Println("【步骤 1】停电操作")
	for i, deviceID := range task.Devices {
		fmt.Printf("  步骤 %d/%d: 打开断路器 %s\n", i+1, len(task.Devices), deviceID)

		cmd := &OpenBreakerCommand{
			commandID: fmt.Sprintf("%s_open_%d", task.TaskID, i),
			Reason:     fmt.Sprintf("检修操作：%s", task.Reason),
			Operator:   o.operatorName,
			TaskID:     task.TaskID,
		}

		// 发送命令到设备 Actor
		if err := o.system.Send(deviceID, cmd); err != nil {
			return fmt.Errorf("发送打开命令到 %s 失败: %w", deviceID, err)
		}

		// 等待操作完成（实际应用中应该通过事件或响应消息）
		time.Sleep(300 * time.Millisecond)

		fmt.Printf("  ✓ 步骤 %d 完成\n", i+1)
	}

	fmt.Println("\n【步骤 2】执行检修操作")
	// 模拟检修操作
	time.Sleep(1 * time.Second)
	fmt.Println("  ✓ 检修操作完成")

	// 步骤3：恢复供电（按顺序关闭所有断路器）
	fmt.Println("\n【步骤 3】恢复供电")
	for i, deviceID := range task.Devices {
		fmt.Printf("  步骤 %d/%d: 关闭断路器 %s\n", i+1, len(task.Devices), deviceID)

		cmd := &CloseBreakerCommand{
			commandID: fmt.Sprintf("%s_close_%d", task.TaskID, i),
			Reason:     fmt.Sprintf("检修完成，恢复供电"),
			Operator:   o.operatorName,
			TaskID:     task.TaskID,
		}

		// 发送命令到设备 Actor
		if err := o.system.Send(deviceID, cmd); err != nil {
			return fmt.Errorf("发送关闭命令到 %s 失败: %w", deviceID, err)
		}

		// 等待操作完成
		time.Sleep(300 * time.Millisecond)

		fmt.Printf("  ✓ 步骤 %d 完成\n", i+1)
	}

	fmt.Println("\n========== 停电检修操作完成 ==========\n")

	// 通知设备完成检修（更新检修时间）
	for _, deviceID := range task.Devices {
		// 这里应该通过消息通知设备更新检修时间
		// 为了简化，我们直接通过 System 获取 Actor 并调用方法
		if actor, exists := o.system.Get(deviceID); exists {
			if breaker, ok := actor.(*BreakerActor); ok {
				breaker.CompleteMaintenance()
			}
		}
	}

	// 发射检修完成事件
	if emitter := o.GetEventEmitter(); emitter != nil {
		_ = emitter.Emit(actors.Event{
			Type: actors.EventTypeCommandCompleted,
			Payload: &MaintenanceCompletedEvent{
				TaskID:     task.TaskID,
				OperatorID: o.operatorID,
				DeviceIDs:  task.Devices,
				Result:     "success",
				Timestamp:  time.Now(),
			},
		})
	}

	// 通知调度中心
	completedEvent := &MaintenanceCompletedEvent{
		TaskID:     task.TaskID,
		OperatorID: o.operatorID,
		DeviceIDs:  task.Devices,
		Result:     "success",
		Timestamp:  time.Now(),
	}
	_ = o.system.Send("DISPATCHER", completedEvent)

	// 清除当前任务
	o.taskMu.Lock()
	o.currentTask = nil
	o.taskMu.Unlock()

	return nil
}

// GetCurrentTask 获取当前任务
func (o *DispatcherOperatorActor) GetCurrentTask() *MaintenanceTask {
	o.taskMu.RLock()
	defer o.taskMu.RUnlock()
	return o.currentTask
}
