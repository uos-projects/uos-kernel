package capacities

import (
	"context"
	"fmt"
)

// ============================================================================
// 示例：SwitchOffCommand - 展示四个维度的完整实现
// ============================================================================

// SwitchOffCommand 关机命令消息
// 注意：实际使用时，SwitchOffCommand 应该在 actors 包中定义并实现 actors.CapabilityCommand
// 这里仅作为示例展示 Capacity 的实现
type SwitchOffCommand struct {
	Reason string // 关机原因
}

// SwitchOffCapacity 实现关机能力
// 展示如何实现四个维度：Command Shape, Guards, Execution Semantics, Effects
type SwitchOffCapacity struct {
	*BaseCapacity
	*CapacityMetadata // 嵌入元数据
}

// NewSwitchOffCapacity 创建关机能力
func NewSwitchOffCapacity(resourceID string) *SwitchOffCapacity {
	// 1. Command Shape: 声明可接受的命令类型
	commands := []CommandType{
		CommandTypeFor((*SwitchOffCommand)(nil), "SwitchOffCommand"),
	}

	// 2. Guards: 声明前置条件
	requirements := []Requirement{
		RequireLifecycleState("online"), // 要求设备在线
		RequirePermission("operator_confirmed"), // 要求操作员确认
	}

	// 3. Execution Semantics: 声明执行语义
	execInfo := ExecutionInfo{
		Mode:          ExecutionModeAsync, // 异步执行
		Interruptible: false,              // 不可中断
		ResourceBound: true,               // 会占用资源（锁定设备）
		Timeout:       30,                 // 30秒超时
		Priority:      5,                  // 优先级5
	}

	// 4. Effects: 声明执行效果
	effects := []Effect{
		EffectStateChangeAndEvent(
			"live_state",           // 状态属性
			"offline",              // 目标值
			"DevicePoweredOff",     // 事件类型
			map[string]interface{}{ // 事件负载
				"resource_id": resourceID,
				"reason":      "command",
			},
		),
	}

	metadata := NewCapacityMetadata(commands, requirements, execInfo, effects)

	return &SwitchOffCapacity{
		BaseCapacity:   NewBaseCapacity("SwitchOffCapacity", resourceID),
		CapacityMetadata: metadata,
	}
}

// CanHandle 实现 Capacity 接口
// 使用 CommandShape 进行类型检查
func (c *SwitchOffCapacity) CanHandle(msg Message) bool {
	return CanHandleByCommandShape(c, msg)
}

// Execute 实现 Capacity 接口
func (c *SwitchOffCapacity) Execute(ctx context.Context, msg Message) error {
	cmd, ok := msg.(*SwitchOffCommand)
	if !ok {
		return fmt.Errorf("invalid message type for SwitchOffCapacity")
	}

	// 注意：在实际实现中，这里应该：
	// 1. 检查前置条件（虽然调度器可能已经检查过，但这里可以再次验证）
	// 2. 执行实际的关机逻辑
	// 3. 应用 Effects（状态变化、事件发射）

	// 示例：执行关机逻辑
	fmt.Printf("Executing SwitchOff for resource %s, reason: %s\n", c.ResourceID(), cmd.Reason)

	// TODO: 实际业务逻辑
	// - 调用设备绑定执行关机
	// - 更新 Actor 状态
	// - 发射事件

	return nil
}

// ============================================================================
// 示例：如何在 Actor 中使用带 Guards 的 Capacity
// ============================================================================

// ExampleUsageInActor 展示如何在 Actor 的 Receive 方法中使用 Guards
// 注意：这个函数是示例代码，实际使用时需要导入 actors 包
// func ExampleUsageInActor(ctx context.Context, actor actors.ResourceActor, capacity Capacity, msg Message) error {
// 	// 1. 检查 Capacity 是否实现了 Guards
// 	if guards, ok := capacity.(Guards); ok {
// 		// 2. 检查前置条件
// 		satisfied, failed, err := guards.CheckGuards(ctx, actor)
// 		if err != nil {
// 			return fmt.Errorf("error checking guards: %w", err)
// 		}
// 		if !satisfied {
// 			return fmt.Errorf("guards not satisfied: %v", failed)
// 		}
// 	}
//
// 	// 3. 检查执行语义（调度器可以使用这个信息）
// 	if semantics, ok := capacity.(ExecutionSemantics); ok {
// 		execInfo := semantics.ExecutionInfo()
// 		// 调度器可以根据 execInfo 决定是否立即执行、排队、或拒绝
// 		_ = execInfo
// 	}
//
// 	// 4. 执行 Capacity
// 	if err := capacity.Execute(ctx, msg); err != nil {
// 		return err
// 	}
//
// 	// 5. 应用 Effects（如果 Capacity 实现了 Effects）
// 	if effects, ok := capacity.(Effects); ok {
// 		for _, effect := range effects.Effects() {
// 			// 应用状态变化
// 			for _, stateChange := range effect.StateChanges {
// 				// 通过消息驱动更新状态
// 				// setPropMsg := &actors.SetPropertyMessage{
// 				// 	Name:  stateChange.Property,
// 				// 	Value: stateChange.ToValue,
// 				// }
// 				// actor.Send(setPropMsg)
// 			}
//
// 			// 发射事件
// 			for _, event := range effect.Events {
// 				// 可以通过事件总线或消息系统发射事件
// 				fmt.Printf("Emitting event: %s, payload: %v\n", event.Type, event.Payload)
// 				// TODO: 实际的事件发射逻辑
// 			}
// 		}
// 	}
//
// 	return nil
// }
