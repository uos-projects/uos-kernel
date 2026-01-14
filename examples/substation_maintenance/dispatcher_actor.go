package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/uos-projects/uos-kernel/actors"
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

func (t *MaintenanceTask) MessageType() actors.MessageCategory {
	return actors.MessageCategoryCoordinationEvent
}

// MaintenancePlan 检修计划
type MaintenancePlan struct {
	PlanID             string
	DeviceID           string
	Interval           time.Duration // 检修间隔
	NextMaintenanceTime time.Time
}

// DispatcherActor 调度中心 Actor
// 接收设备异常事件和检修需求事件，制定检修计划并分配给操作员
type DispatcherActor struct {
	*actors.BaseResourceActor

	// 检修计划
	maintenancePlans []MaintenancePlan
	plansMu          sync.RWMutex

	// 待处理任务
	pendingTasks []MaintenanceTask
	tasksMu      sync.RWMutex

	// 操作员列表
	operators []string

	// 系统引用（用于发送消息）
	system *actors.System
}

// NewDispatcherActor 创建调度中心 Actor
func NewDispatcherActor(system *actors.System) *DispatcherActor {
	actor := &DispatcherActor{
		BaseResourceActor: actors.NewBaseResourceActor("DISPATCHER", "Dispatcher", nil),
		maintenancePlans:  make([]MaintenancePlan, 0),
		pendingTasks:      make([]MaintenanceTask, 0),
		operators:         make([]string, 0),
		system:            system,
	}

	return actor
}

// RegisterOperator 注册操作员
func (d *DispatcherActor) RegisterOperator(operatorID string) {
	d.operators = append(d.operators, operatorID)
	fmt.Printf("[调度中心] 注册操作员：%s\n", operatorID)
}

// Receive 重写消息处理逻辑
func (d *DispatcherActor) Receive(ctx context.Context, msg actors.Message) error {
	// 处理设备事件
	switch event := msg.(type) {
	case *DeviceAbnormalEvent:
		return d.handleDeviceAbnormalEvent(ctx, event)
	case *MaintenanceRequiredEvent:
		return d.handleMaintenanceRequiredEvent(ctx, event)
	case *MaintenanceCompletedEvent:
		return d.handleMaintenanceCompletedEvent(ctx, event)
	}

	// 其他消息交给基类处理
	return d.BaseResourceActor.Receive(ctx, msg)
}

// handleDeviceAbnormalEvent 处理设备异常事件
func (d *DispatcherActor) handleDeviceAbnormalEvent(ctx context.Context, event *DeviceAbnormalEvent) error {
	fmt.Printf("\n[调度中心] 📢 收到设备异常事件：\n")
	fmt.Printf("  设备：%s\n", event.DeviceID)
	fmt.Printf("  异常类型：%s\n", event.EventType)
	fmt.Printf("  严重程度：%s\n", event.Severity)
	fmt.Printf("  详情：%v\n", event.Details)

	// 创建紧急检修任务
	task := d.createEmergencyMaintenanceTask(event)
	
	d.tasksMu.Lock()
	d.pendingTasks = append(d.pendingTasks, task)
	d.tasksMu.Unlock()

	fmt.Printf("[调度中心] ✅ 已创建紧急检修任务：%s\n", task.TaskID)

	// 分配给操作员
	return d.assignTaskToOperator(task)
}

// handleMaintenanceRequiredEvent 处理需要检修事件
func (d *DispatcherActor) handleMaintenanceRequiredEvent(ctx context.Context, event *MaintenanceRequiredEvent) error {
	fmt.Printf("\n[调度中心] 📢 收到检修需求事件：\n")
	fmt.Printf("  设备：%s\n", event.DeviceID)
	fmt.Printf("  原因：%s\n", event.Reason)
	fmt.Printf("  运行小时数：%d\n", event.OperationHours)

	// 创建定期检修任务
	task := d.createScheduledMaintenanceTask(event)
	
	d.tasksMu.Lock()
	d.pendingTasks = append(d.pendingTasks, task)
	d.tasksMu.Unlock()

	fmt.Printf("[调度中心] ✅ 已创建定期检修任务：%s\n", task.TaskID)

	// 分配给操作员
	return d.assignTaskToOperator(task)
}

// handleMaintenanceCompletedEvent 处理检修完成事件
func (d *DispatcherActor) handleMaintenanceCompletedEvent(ctx context.Context, event *MaintenanceCompletedEvent) error {
	fmt.Printf("\n[调度中心] 📢 收到检修完成事件：\n")
	fmt.Printf("  任务ID：%s\n", event.TaskID)
	fmt.Printf("  操作员：%s\n", event.OperatorID)
	fmt.Printf("  结果：%s\n", event.Result)

	// 更新任务状态
	d.tasksMu.Lock()
	for i, task := range d.pendingTasks {
		if task.TaskID == event.TaskID {
			d.pendingTasks[i].Status = event.Result
			break
		}
	}
	d.tasksMu.Unlock()

	return nil
}

// createEmergencyMaintenanceTask 创建紧急检修任务
func (d *DispatcherActor) createEmergencyMaintenanceTask(event *DeviceAbnormalEvent) MaintenanceTask {
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

// createScheduledMaintenanceTask 创建定期检修任务
func (d *DispatcherActor) createScheduledMaintenanceTask(event *MaintenanceRequiredEvent) MaintenanceTask {
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

// assignTaskToOperator 分配任务给操作员
func (d *DispatcherActor) assignTaskToOperator(task MaintenanceTask) error {
	if len(d.operators) == 0 {
		return fmt.Errorf("没有可用的操作员")
	}

	// 简单分配：选择第一个操作员（实际应用中可以实现更复杂的调度算法）
	operatorID := d.operators[0]
	task.AssignedTo = operatorID
	task.Status = "assigned"

	// 更新任务列表
	d.tasksMu.Lock()
	for i, t := range d.pendingTasks {
		if t.TaskID == task.TaskID {
			d.pendingTasks[i] = task
			break
		}
	}
	d.tasksMu.Unlock()

	// 发送任务给操作员
	if err := d.system.Send(operatorID, &task); err != nil {
		return fmt.Errorf("发送任务给操作员失败: %w", err)
	}

	// 发射任务分配事件
	if emitter := d.GetEventEmitter(); emitter != nil {
		_ = emitter.Emit(actors.Event{
			Type: actors.EventTypeStateChanged,
			Payload: &MaintenanceTaskAssignedEvent{
				TaskID:     task.TaskID,
				OperatorID: operatorID,
				DeviceIDs:  task.Devices,
				Reason:     task.Reason,
				Timestamp:  time.Now(),
			},
		})
	}

	fmt.Printf("[调度中心] 📤 已将任务 %s 分配给操作员 %s\n", task.TaskID, operatorID)

	return nil
}

// GetPendingTasks 获取待处理任务列表
func (d *DispatcherActor) GetPendingTasks() []MaintenanceTask {
	d.tasksMu.RLock()
	defer d.tasksMu.RUnlock()

	result := make([]MaintenanceTask, len(d.pendingTasks))
	copy(result, d.pendingTasks)
	return result
}
