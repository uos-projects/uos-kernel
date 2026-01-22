package main

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/uos-projects/uos-kernel/actor"
)

// MaintenanceTask 检修任务
// 实现 Message 接口，可以作为消息发送
type MaintenanceTask struct {
	TaskID      string
	Type        string // "scheduled", "emergency"
	Devices     []string
	Description string
	Reason      string
	AssignedTo  string
	Status      string // "pending", "in_progress", "completed", "failed"
	CreatedAt   time.Time
}

func (t *MaintenanceTask) MessageType() actor.MessageCategory {
	return actor.MessageCategoryCoordinationEvent
}

// MaintenancePlan 检修计划
type MaintenancePlan struct {
	PlanID              string
	DeviceID            string
	Interval            time.Duration // 检修间隔
	NextMaintenanceTime time.Time
}

// DispatcherActor 调度中心 Actor
// 接收设备异常事件和检修需求事件，制定检修计划并分配给操作员
type DispatcherActor struct {
	*actor.BaseResourceActor

	// 检修计划
	maintenancePlans []MaintenancePlan
	plansMu          sync.RWMutex

	// 待处理任务
	pendingTasks []MaintenanceTask
	tasksMu      sync.RWMutex

	// 操作员列表
	operators []string

	// 系统引用（用于发送消息）
	system *actor.System
}

// NewDispatcherActor 创建调度中心 Actor
func NewDispatcherActor(system *actor.System) *DispatcherActor {
	d := &DispatcherActor{
		BaseResourceActor: actor.NewBaseResourceActor("DISPATCHER", "Dispatcher"),
		maintenancePlans:  make([]MaintenancePlan, 0),
		pendingTasks:      make([]MaintenanceTask, 0),
		operators:         make([]string, 0),
		system:            system,
	}

	// 注册业务事件
	d.registerBusinessEvents()

	return d
}

// registerBusinessEvents 注册业务事件
func (d *DispatcherActor) registerBusinessEvents() {
	// 注册任务创建事件
	taskCreatedEventDesc := actor.NewEventDescriptor(
		"MaintenanceTaskCreatedEvent",
		actor.EventTypeStateChanged,
		reflect.TypeOf((*MaintenanceTaskCreatedEvent)(nil)).Elem(),
		"检修任务创建事件",
		d.ResourceID(),
	)
	d.RegisterEvent(taskCreatedEventDesc)

	// 注册任务分配事件
	taskAssignedEventDesc := actor.NewEventDescriptor(
		"MaintenanceTaskAssignedEvent",
		actor.EventTypeStateChanged,
		reflect.TypeOf((*MaintenanceTaskAssignedEvent)(nil)).Elem(),
		"检修任务分配事件",
		d.ResourceID(),
	)
	d.RegisterEvent(taskAssignedEventDesc)

	// 注册任务更新事件
	taskUpdatedEventDesc := actor.NewEventDescriptor(
		"MaintenanceTaskUpdatedEvent",
		actor.EventTypeStateChanged,
		reflect.TypeOf((*MaintenanceTaskUpdatedEvent)(nil)).Elem(),
		"检修任务更新事件",
		d.ResourceID(),
	)
	d.RegisterEvent(taskUpdatedEventDesc)
}

// RegisterOperator 注册操作员
func (d *DispatcherActor) RegisterOperator(operatorID string) {
	d.operators = append(d.operators, operatorID)
	fmt.Printf("[调度中心] 注册操作员：%s\n", operatorID)
}

// Receive 重写消息处理逻辑
func (d *DispatcherActor) Receive(ctx context.Context, msg actor.Message) error {
	// 优先处理世界事件（Coordination Events）
	// 注意：必须在基类 Receive() 之前处理，因为基类会根据 MessageType() 路由
	// 而 MaintenanceRequiredEvent 等是 MessageCategoryCoordinationEvent，
	// 会被基类的 handleCoordinationEvent() 的 default 分支忽略
	switch event := msg.(type) {
	case *DeviceAbnormalEvent:
		return d.handleWorldEvent(ctx, event)
	case *MaintenanceRequiredEvent:
		return d.handleWorldEvent(ctx, event)
	case *MaintenanceCompletedEvent:
		return d.handleMaintenanceCompletedEvent(ctx, event)
	}

	// 其他消息交给基类处理
	return d.BaseResourceActor.Receive(ctx, msg)
}

// handleWorldEvent 处理世界事件（统一入口）
func (d *DispatcherActor) handleWorldEvent(ctx context.Context, event actor.Message) error {
	var task MaintenanceTask

	switch e := event.(type) {
	case *DeviceAbnormalEvent:
		fmt.Printf("\n[调度中心] 📢 收到设备异常事件：\n")
		fmt.Printf("  设备：%s\n", e.DeviceID)
		fmt.Printf("  异常类型：%s\n", e.EventType)
		fmt.Printf("  严重程度：%s\n", e.Severity)
		fmt.Printf("  详情：%v\n", e.Details)

		// 通过领域方法创建任务
		task = d.createEmergencyTaskFrom(e)

	case *MaintenanceRequiredEvent:
		fmt.Printf("\n[调度中心] 📢 收到检修需求事件：\n")
		fmt.Printf("  设备：%s\n", e.DeviceID)
		fmt.Printf("  原因：%s\n", e.Reason)
		fmt.Printf("  运行小时数：%d\n", e.OperationHours)

		// 通过领域方法创建任务
		task = d.createScheduledTaskFrom(e)

	default:
		return fmt.Errorf("unknown world event type: %T", event)
	}

	// 应用任务创建（更新内部状态）
	d.applyTaskCreated(task)

	// 发射任务创建事件
	d.emitTaskCreatedEvent(task)

	fmt.Printf("[调度中心] ✅ 已创建检修任务：%s\n", task.TaskID)

	// 分配给操作员（通过命令）
	return d.assignTaskToOperator(task)
}

// handleMaintenanceCompletedEvent 处理检修完成事件
func (d *DispatcherActor) handleMaintenanceCompletedEvent(ctx context.Context, event *MaintenanceCompletedEvent) error {
	fmt.Printf("\n[调度中心] 📢 收到检修完成事件：\n")
	fmt.Printf("  任务ID：%s\n", event.TaskID)
	fmt.Printf("  操作员：%s\n", event.OperatorID)
	fmt.Printf("  结果：%s\n", event.Result)

	// 应用任务状态更新
	d.applyTaskStatusUpdate(event.TaskID, event.Result)

	// 发射任务更新事件
	d.emitTaskUpdatedEvent(event.TaskID, event.Result)

	return nil
}

// ============================================================================
// 领域方法（Domain Methods）
// ============================================================================

// createEmergencyTaskFrom 从异常事件创建紧急检修任务（领域方法）
func (d *DispatcherActor) createEmergencyTaskFrom(event *DeviceAbnormalEvent) MaintenanceTask {
	return MaintenanceTask{
		TaskID:      fmt.Sprintf("TASK-EMERGENCY-%d", time.Now().Unix()),
		Type:        "emergency",
		Devices:     []string{event.DeviceID},
		Description: fmt.Sprintf("紧急检修：%s - %s", event.DeviceID, event.EventType),
		Reason:      fmt.Sprintf("设备异常：%s", event.EventType),
		Status:      "pending",
		CreatedAt:   time.Now(),
	}
}

// createScheduledTaskFrom 从定期事件创建定期检修任务（领域方法）
func (d *DispatcherActor) createScheduledTaskFrom(event *MaintenanceRequiredEvent) MaintenanceTask {
	return MaintenanceTask{
		TaskID:      fmt.Sprintf("TASK-SCHEDULED-%d", time.Now().Unix()),
		Type:        "scheduled",
		Devices:     []string{event.DeviceID},
		Description: fmt.Sprintf("定期检修：%s", event.DeviceID),
		Reason:      fmt.Sprintf("运行时间达到检修间隔：%d 小时", event.OperationHours),
		Status:      "pending",
		CreatedAt:   time.Now(),
	}
}

// applyTaskCreated 应用任务创建（更新内部状态）
func (d *DispatcherActor) applyTaskCreated(task MaintenanceTask) {
	d.tasksMu.Lock()
	defer d.tasksMu.Unlock()
	d.pendingTasks = append(d.pendingTasks, task)
}

// applyTaskAssigned 应用任务分配（更新内部状态）
func (d *DispatcherActor) applyTaskAssigned(taskID string, operatorID string) {
	d.tasksMu.Lock()
	defer d.tasksMu.Unlock()
	for i, task := range d.pendingTasks {
		if task.TaskID == taskID {
			d.pendingTasks[i].AssignedTo = operatorID
			d.pendingTasks[i].Status = "assigned"
			break
		}
	}
}

// applyTaskStatusUpdate 应用任务状态更新
func (d *DispatcherActor) applyTaskStatusUpdate(taskID string, status string) {
	d.tasksMu.Lock()
	defer d.tasksMu.Unlock()
	for i, task := range d.pendingTasks {
		if task.TaskID == taskID {
			d.pendingTasks[i].Status = status
			break
		}
	}
}

// assignTaskToOperator 分配任务给操作员（通过命令）
func (d *DispatcherActor) assignTaskToOperator(task MaintenanceTask) error {
	if len(d.operators) == 0 {
		return fmt.Errorf("没有可用的操作员")
	}

	// 简单分配：选择第一个操作员（实际应用中可以实现更复杂的调度算法）
	operatorID := d.operators[0]

	// 应用任务分配（更新内部状态）
	d.applyTaskAssigned(task.TaskID, operatorID)

	// 发送 StartMaintenanceCommand 给操作员（而不是直接传 MaintenanceTask）
	cmd := &StartMaintenanceCommand{
		TaskID:      task.TaskID,
		Type:        task.Type,
		Devices:     task.Devices,
		Description: task.Description,
		Reason:      task.Reason,
		OperatorID:  operatorID,
	}

	if err := d.system.Send(operatorID, cmd); err != nil {
		return fmt.Errorf("发送开始检修命令给操作员失败: %w", err)
	}

	// 发射任务分配事件
	d.emitTaskAssignedEvent(task.TaskID, operatorID, task.Devices, task.Reason)

	fmt.Printf("[调度中心] 📤 已将任务 %s 分配给操作员 %s\n", task.TaskID, operatorID)

	return nil
}

// ============================================================================
// 事件发射方法（Event Emission）
// ============================================================================

// emitTaskCreatedEvent 发射任务创建事件
func (d *DispatcherActor) emitTaskCreatedEvent(task MaintenanceTask) {
	if emitter := d.GetEventEmitter(); emitter != nil {
		_ = emitter.Emit(actor.Event{
			Type: actor.EventTypeStateChanged,
			Payload: &MaintenanceTaskCreatedEvent{
				TaskID:      task.TaskID,
				Type:        task.Type,
				Devices:     task.Devices,
				Description: task.Description,
				Reason:      task.Reason,
				Timestamp:   task.CreatedAt,
			},
		})
	}
}

// emitTaskAssignedEvent 发射任务分配事件
func (d *DispatcherActor) emitTaskAssignedEvent(taskID string, operatorID string, deviceIDs []string, reason string) {
	if emitter := d.GetEventEmitter(); emitter != nil {
		_ = emitter.Emit(actor.Event{
			Type: actor.EventTypeStateChanged,
			Payload: &MaintenanceTaskAssignedEvent{
				TaskID:     taskID,
				OperatorID: operatorID,
				DeviceIDs:  deviceIDs,
				Reason:     reason,
				Timestamp:  time.Now(),
			},
		})
	}
}

// emitTaskUpdatedEvent 发射任务更新事件
func (d *DispatcherActor) emitTaskUpdatedEvent(taskID string, status string) {
	if emitter := d.GetEventEmitter(); emitter != nil {
		_ = emitter.Emit(actor.Event{
			Type: actor.EventTypeStateChanged,
			Payload: &MaintenanceTaskUpdatedEvent{
				TaskID:    taskID,
				Status:    status,
				Timestamp: time.Now(),
			},
		})
	}
}

// GetPendingTasks 获取待处理任务列表
func (d *DispatcherActor) GetPendingTasks() []MaintenanceTask {
	d.tasksMu.RLock()
	defer d.tasksMu.RUnlock()

	result := make([]MaintenanceTask, len(d.pendingTasks))
	copy(result, d.pendingTasks)
	return result
}
