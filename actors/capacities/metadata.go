package capacities

import (
	"context"
	"fmt"
	"reflect"
)

// Metadata 是 Capacity 元数据的聚合接口
// 实现了所有四个维度的声明式接口
type Metadata interface {
	CommandShape
	Guards
	ExecutionSemantics
	Effects
}

// CapacityMetadata Capacity 的完整元数据
type CapacityMetadata struct {
	// CommandShape
	commands []CommandType

	// Guards
	guards *BaseGuards

	// ExecutionSemantics
	execInfo ExecutionInfo

	// Effects
	effects []Effect
}

// NewCapacityMetadata 创建 Capacity 元数据
func NewCapacityMetadata(
	commands []CommandType,
	requirements []Requirement,
	execInfo ExecutionInfo,
	effects []Effect,
) *CapacityMetadata {
	return &CapacityMetadata{
		commands: commands,
		guards:   NewBaseGuards(requirements),
		execInfo: execInfo,
		effects:  effects,
	}
}

// AcceptableCommands 实现 CommandShape
func (m *CapacityMetadata) AcceptableCommands() []CommandType {
	return m.commands
}

// Requires 实现 Guards
func (m *CapacityMetadata) Requires() []Requirement {
	return m.guards.Requires()
}

// CheckGuards 实现 Guards
func (m *CapacityMetadata) CheckGuards(ctx context.Context, actor interface{}) (bool, []string, error) {
	return m.guards.CheckGuards(ctx, actor)
}

// ExecutionInfo 实现 ExecutionSemantics
func (m *CapacityMetadata) ExecutionInfo() ExecutionInfo {
	return m.execInfo
}

// Effects 实现 Effects
func (m *CapacityMetadata) Effects() []Effect {
	return m.effects
}

// ============================================================================
// 辅助函数：创建常用 Requirement
// ============================================================================

// RequireState 创建状态要求（要求 Actor 的某个状态属性等于指定值）
func RequireState(propertyName string, expectedValue interface{}) Requirement {
	return Requirement{
		Name:        fmt.Sprintf("state.%s == %v", propertyName, expectedValue),
		Description: fmt.Sprintf("%s == %v", propertyName, expectedValue),
		Check: func(ctx context.Context, actor interface{}) (bool, error) {
			// 尝试通过 PropertyHolder 接口获取属性
			if holder, ok := actor.(interface {
				GetProperty(name string) (interface{}, bool)
			}); ok {
				value, exists := holder.GetProperty(propertyName)
				if !exists {
					return false, nil
				}
				return value == expectedValue, nil
			}
			// 尝试通过反射获取字段
			actorVal := reflect.ValueOf(actor)
			if actorVal.Kind() == reflect.Ptr {
				actorVal = actorVal.Elem()
			}
			field := actorVal.FieldByName(propertyName)
			if !field.IsValid() {
				return false, fmt.Errorf("property %s not found", propertyName)
			}
			return field.Interface() == expectedValue, nil
		},
	}
}

// RequireLifecycleState 创建生命周期状态要求（要求 Actor 的生命周期状态为指定值）
func RequireLifecycleState(expectedState string) Requirement {
	return Requirement{
		Name:        fmt.Sprintf("lifecycle_state == %s", expectedState),
		Description: fmt.Sprintf("lifecycle_state == %s", expectedState),
		Check: func(ctx context.Context, actor interface{}) (bool, error) {
			// 尝试通过 CurrentState 方法获取状态
			if stateGetter, ok := actor.(interface {
				CurrentState() interface{}
			}); ok {
				state := stateGetter.CurrentState()
				return fmt.Sprintf("%v", state) == expectedState, nil
			}
			// 尝试通过属性获取
			if holder, ok := actor.(interface {
				GetProperty(name string) (interface{}, bool)
			}); ok {
				value, exists := holder.GetProperty("lifecycle_state")
				if !exists {
					return false, nil
				}
				return fmt.Sprintf("%v", value) == expectedState, nil
			}
			return false, fmt.Errorf("cannot check lifecycle state: actor does not implement CurrentState() or PropertyHolder")
		},
	}
}

// RequirePermission 创建权限要求（示例：要求操作员确认）
func RequirePermission(permissionName string) Requirement {
	return Requirement{
		Name:        fmt.Sprintf("permission.%s", permissionName),
		Description: fmt.Sprintf("permission == %s", permissionName),
		Check: func(ctx context.Context, actor interface{}) (bool, error) {
			// 这里可以从 context 中获取权限信息
			// 或者从 Actor 的属性中获取
			if holder, ok := actor.(interface {
				GetProperty(name string) (interface{}, bool)
			}); ok {
				value, exists := holder.GetProperty("permission")
				if !exists {
					return false, nil
				}
				return fmt.Sprintf("%v", value) == permissionName, nil
			}
			// 也可以从 context 中获取（如果权限信息存储在 context 中）
			if perm := ctx.Value("permission"); perm != nil {
				return fmt.Sprintf("%v", perm) == permissionName, nil
			}
			return false, nil
		},
	}
}

// ============================================================================
// 辅助函数：创建常用 CommandType
// ============================================================================

// CommandTypeFor 根据消息类型创建 CommandType
func CommandTypeFor(msgType Message, typeName string) CommandType {
	return CommandType{
		TypeName:    typeName,
		MessageType: reflect.TypeOf(msgType),
	}
}

// ============================================================================
// 辅助函数：创建常用 Effect
// ============================================================================

// EffectStateChange 创建状态变化效果
func EffectStateChange(property string, toValue interface{}) Effect {
	return Effect{
		StateChanges: []StateChange{
			{
				Property: property,
				ToValue:  toValue,
			},
		},
		Events: nil,
	}
}

// EffectEvent 创建事件效果
func EffectEvent(eventType string, payload interface{}) Effect {
	return Effect{
		StateChanges: nil,
		Events: []Event{
			{
				Type:    eventType,
				Payload: payload,
			},
		},
	}
}

// EffectStateChangeAndEvent 创建状态变化和事件效果
func EffectStateChangeAndEvent(property string, toValue interface{}, eventType string, eventPayload interface{}) Effect {
	return Effect{
		StateChanges: []StateChange{
			{
				Property: property,
				ToValue:  toValue,
			},
		},
		Events: []Event{
			{
				Type:    eventType,
				Payload: eventPayload,
			},
		},
	}
}
