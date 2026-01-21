package main

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/uos-projects/uos-kernel/actors"
)

// DispatcherOperatorActor 调度操作员 Actor
// 代表操作员的数字化实体，只反映操作员的状态，不执行行为
// 实际行为由 SimulatedOperatorBinding 执行
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
	actor := &DispatcherOperatorActor{
		BaseResourceActor: actors.NewBaseResourceActor(id, "DispatcherOperator", nil),
		operatorID:        id,
		operatorName:      name,
		system:            system,
	}

	// 注册业务事件
	actor.registerBusinessEvents()

	return actor
}

// registerBusinessEvents 注册业务事件
func (o *DispatcherOperatorActor) registerBusinessEvents() {
	// 注册检修完成事件
	maintenanceCompletedEventDesc := actors.NewEventDescriptor(
		"MaintenanceCompletedEvent",
		actors.EventTypeCommandCompleted,
		reflect.TypeOf((*MaintenanceCompletedEvent)(nil)).Elem(),
		"检修完成事件",
		o.ResourceID(),
	)
	o.RegisterEvent(maintenanceCompletedEventDesc)
}

// Receive 重写消息处理逻辑
func (o *DispatcherOperatorActor) Receive(ctx context.Context, msg actors.Message) error {
	// 处理开始检修命令（通过 Binding 执行）
	switch cmd := msg.(type) {
	case *StartMaintenanceCommand:
		return o.handleStartMaintenanceCommand(ctx, cmd)
	}

	// 处理来自 Binding 的外部事件（操作员状态反馈）
	switch m := msg.(type) {
	case *actors.ExternalEventMessage:
		if m.BindingType == actors.BindingTypeHuman {
			return o.handleOperatorDeviceEvent(ctx, m.Event)
		}
	}

	// 其他消息交给基类处理
	return o.BaseResourceActor.Receive(ctx, msg)
}

// handleStartMaintenanceCommand 处理开始检修命令
// Actor 只负责状态管理：接收命令 -> 更新状态 -> 通过 Binding 执行行为
func (o *DispatcherOperatorActor) handleStartMaintenanceCommand(ctx context.Context, cmd *StartMaintenanceCommand) error {
	fmt.Printf("\n[操作员 Actor %s] 📋 收到开始检修命令：\n", o.operatorName)
	fmt.Printf("  任务ID：%s\n", cmd.TaskID)
	fmt.Printf("  类型：%s\n", cmd.Type)
	fmt.Printf("  设备：%v\n", cmd.Devices)
	fmt.Printf("  原因：%s\n", cmd.Reason)

	// 从命令构造任务对象
	task := &MaintenanceTask{
		TaskID:      cmd.TaskID,
		Type:        cmd.Type,
		Devices:     cmd.Devices,
		Description: cmd.Description,
		Reason:      cmd.Reason,
		AssignedTo:  cmd.OperatorID,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	// 更新状态：接收任务
	if err := o.acceptTask(task); err != nil {
		return err
	}

	// 检查 Binding 是否存在
	if _, exists := o.GetBinding(actors.BindingTypeHuman); !exists {
		return fmt.Errorf("操作员绑定未找到")
	}

	// 通过 Binding 执行实际行为（发送 ExecuteExternalCommandMessage）
	// 消息会被 BaseResourceActor.handleCoordinationEvent 处理，然后调用 Binding.ExecuteExternal
	executeMsg := &actors.ExecuteExternalCommandMessage{
		BindingType: actors.BindingTypeHuman,
		Command:     cmd,
	}
	// 发送消息到 Actor 的消息循环，由 BaseResourceActor 处理
	if !o.Send(executeMsg) {
		return fmt.Errorf("发送执行外部命令消息失败")
	}
	return nil
}

// ============================================================================
// 领域方法（Domain Methods）- 只负责状态管理
// ============================================================================

// acceptTask 接收任务并"接单"（状态更新）
func (o *DispatcherOperatorActor) acceptTask(task *MaintenanceTask) error {
	o.taskMu.Lock()
	o.currentTask = task
	o.currentTask.Status = "in_progress"
	o.taskMu.Unlock()

	fmt.Printf("[操作员 Actor %s] ✅ 已接受任务：%s\n", o.operatorName, task.TaskID)

	// 发射任务开始事件
	if emitter := o.GetEventEmitter(); emitter != nil {
		_ = emitter.Emit(actors.Event{
			Type: actors.EventTypeStateChanged,
			Payload: map[string]interface{}{
				"event_type":  "MaintenanceTaskStarted",
				"task_id":     task.TaskID,
				"operator_id": o.operatorID,
				"timestamp":   time.Now(),
			},
		})
	}

	return nil
}

// handleOperatorDeviceEvent 处理来自 Binding 的操作员状态反馈事件
// 这是"现实驱动事件"：Binding 执行行为后，反馈状态变化给 Actor
func (o *DispatcherOperatorActor) handleOperatorDeviceEvent(ctx context.Context, event interface{}) error {
	operatorEvent, ok := event.(*OperatorDeviceEvent)
	if !ok {
		return fmt.Errorf("invalid operator device event type: %T", event)
	}

	fmt.Printf("[操作员 Actor %s] 📨 收到状态反馈：%s (任务: %s)\n", o.operatorName, operatorEvent.Action, operatorEvent.TaskID)

	o.taskMu.RLock()
	task := o.currentTask
	o.taskMu.RUnlock()

	if task == nil || task.TaskID != operatorEvent.TaskID {
		return fmt.Errorf("任务不存在或任务ID不匹配")
	}

	// 根据事件类型更新状态并发射事件
	switch operatorEvent.Action {
	case "task_started":
		// 任务已开始（Binding 已开始执行）
		o.updateTaskStatus("in_progress")
		o.emitStepEvent("task_started", task.TaskID)

	case "power_outage_completed":
		// 停电完成
		o.emitStepEvent("power_outage_completed", task.TaskID)

	case "maintenance_completed":
		// 检修完成
		o.emitStepEvent("maintenance_completed", task.TaskID)

	case "power_restored":
		// 恢复供电完成
		o.emitStepEvent("power_restored", task.TaskID)

	case "task_completed":
		// 任务完成
		o.updateTaskStatus("completed")
		o.finishTask(ctx, task, operatorEvent.Result)
	}

	return nil
}

// updateTaskStatus 更新任务状态
func (o *DispatcherOperatorActor) updateTaskStatus(status string) {
	o.taskMu.Lock()
	if o.currentTask != nil {
		o.currentTask.Status = status
	}
	o.taskMu.Unlock()
}

// emitStepEvent 发射步骤完成事件
func (o *DispatcherOperatorActor) emitStepEvent(step string, taskID string) {
	if emitter := o.GetEventEmitter(); emitter != nil {
		_ = emitter.Emit(actors.Event{
			Type: actors.EventTypeStateChanged,
			Payload: map[string]interface{}{
				"event_type": "MaintenanceTaskStepCompleted",
				"task_id":    taskID,
				"step":       step,
				"timestamp":  time.Now(),
			},
		})
	}
}

// finishTask 完成任务（状态更新和事件通知）
func (o *DispatcherOperatorActor) finishTask(ctx context.Context, task *MaintenanceTask, result string) {
	// 发射检修完成事件
	if emitter := o.GetEventEmitter(); emitter != nil {
		_ = emitter.Emit(actors.Event{
			Type: actors.EventTypeCommandCompleted,
			Payload: &MaintenanceCompletedEvent{
				TaskID:     task.TaskID,
				OperatorID: o.operatorID,
				DeviceIDs:  task.Devices,
				Result:     result,
				Timestamp:  time.Now(),
			},
		})
	}

	// 通知调度中心
	completedEvent := &MaintenanceCompletedEvent{
		TaskID:     task.TaskID,
		OperatorID: o.operatorID,
		DeviceIDs:  task.Devices,
		Result:     result,
		Timestamp:  time.Now(),
	}
	_ = o.system.Send("DISPATCHER", completedEvent)

	// 清除当前任务
	o.taskMu.Lock()
	o.currentTask = nil
	o.taskMu.Unlock()

	fmt.Printf("[操作员 Actor %s] ✅ 任务完成：%s (结果: %s)\n", o.operatorName, task.TaskID, result)
}

// GetCurrentTask 获取当前任务
func (o *DispatcherOperatorActor) GetCurrentTask() *MaintenanceTask {
	o.taskMu.RLock()
	defer o.taskMu.RUnlock()
	return o.currentTask
}
