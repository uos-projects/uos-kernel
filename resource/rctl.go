package resource

import (
	"context"
	"fmt"

	"github.com/uos-projects/uos-kernel/actors"
	"github.com/uos-projects/uos-kernel/actors/capacities"
)

// ControlCommand 控制命令（类似 ioctl 的命令）
type ControlCommand int32

const (
	// General Commands
	CMD_GET_RESOURCE_INFO ControlCommand = iota + 0x1000
	CMD_LIST_CAPABILITIES // 0x1001

	// Control Commands (Mapped to Capacities)
	CMD_ACCUMULATOR_RESET // 0x1002
	CMD_COMMAND           // 0x1003
	CMD_RAISE_LOWER       // 0x1004
	CMD_SET_POINT         // 0x1005
)

// System Commands
const (
	CMD_SYNC ControlCommand = 0x2001 // 强制持久化同步
)

// ControlArg 控制参数（类似 ioctl 的参数）
type ControlArg interface{}

// ResourceInfo 资源信息
type ResourceInfo struct {
	ResourceID   string
	ResourceType string
}

// RCtl 资源控制（类似 POSIX ioctl）
// 构造 Control 消息并发送给 Actor
func (rm *ResourceManager) RCtl(ctx context.Context, fd ResourceDescriptor, cmd ControlCommand, arg ControlArg) (interface{}, error) {
	handle, err := rm.GetHandle(fd)
	if err != nil {
		return nil, err
	}

	// 通过Resource访问Actor
	resource := handle.resource
	resource.mu.RLock()
	actor := resource.actor
	resource.mu.RUnlock()

	// 根据命令构造对应的 Control 消息
	var msg capacities.Message
	var buildErr error

	switch cmd {
	case CMD_ACCUMULATOR_RESET:
		msg, buildErr = rm.buildAccumulatorResetMessage(arg)

	case CMD_COMMAND:
		msg, buildErr = rm.buildCommandMessage(arg)

	case CMD_RAISE_LOWER:
		msg, buildErr = rm.buildRaiseLowerMessage(arg)

	case CMD_SET_POINT:
		msg, buildErr = rm.buildSetPointMessage(arg)

	case CMD_GET_RESOURCE_INFO:
		// 查询类命令，直接返回结果
		return &ResourceInfo{
			ResourceID:   actor.ResourceID(),
			ResourceType: actor.ResourceType(),
		}, nil

	case CMD_LIST_CAPABILITIES:
		return actor.ListCapabilities(), nil

	case CMD_SYNC:
		// 强制同步状态到持久化层
		// 这里虽然使用了具体的类型断言，但作为内核层，了解 CIMResourceActor 是可以接受的
		// 或者可以定义一个 Snapshotable 接口
		if cimActor, ok := actor.(*actors.CIMResourceActor); ok {
			if err := cimActor.SaveSnapshotAsync(ctx); err != nil {
				return nil, fmt.Errorf("failed to sync snapshot: %w", err)
			}
			return nil, nil // API 返回成功，异步保存
		}
		return nil, fmt.Errorf("resource does not support sync")

	default:
		return nil, fmt.Errorf("unknown control command: %d", cmd)
	}

	if buildErr != nil {
		return nil, buildErr
	}

	// 通过 Actor.Send() 发送消息
	// Actor 的 Receive() 会自动路由到对应的 Capacity
	if !actor.Send(msg) {
		return nil, fmt.Errorf("failed to send control message: mailbox full")
	}

	return nil, nil
}

// 辅助函数：构造各种 Control 消息
func (rm *ResourceManager) buildAccumulatorResetMessage(arg ControlArg) (capacities.Message, error) {
	req, ok := arg.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid argument type for CMD_ACCUMULATOR_RESET")
	}

	value := 0
	if v, ok := req["value"]; ok {
		if fv, ok := v.(float64); ok {
			value = int(fv)
		} else if iv, ok := v.(int); ok {
			value = iv
		}
	}

	return &capacities.AccumulatorResetMessage{Value: value}, nil
}

func (rm *ResourceManager) buildCommandMessage(arg ControlArg) (capacities.Message, error) {
	req, ok := arg.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid argument type for CMD_COMMAND")
	}

	command := ""
	if cmd, ok := req["command"].(string); ok {
		command = cmd
	}

	value := 0
	if v, ok := req["value"]; ok {
		if fv, ok := v.(float64); ok {
			value = int(fv)
		} else if iv, ok := v.(int); ok {
			value = iv
		}
	}

	return &capacities.CommandMessage{
		Command: command,
		Value:   value,
	}, nil
}

func (rm *ResourceManager) buildRaiseLowerMessage(arg ControlArg) (capacities.Message, error) {
	req, ok := arg.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid argument type for CMD_RAISE_LOWER")
	}

	delta := 0.0
	if v, ok := req["delta"]; ok {
		if fv, ok := v.(float64); ok {
			delta = fv
		} else if iv, ok := v.(int); ok {
			delta = float64(iv)
		}
	}

	return &capacities.RaiseLowerCommandMessage{Delta: delta}, nil
}

func (rm *ResourceManager) buildSetPointMessage(arg ControlArg) (capacities.Message, error) {
	req, ok := arg.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid argument type for CMD_SET_POINT")
	}

	value := 0.0
	if v, ok := req["value"]; ok {
		if fv, ok := v.(float64); ok {
			value = fv
		} else if iv, ok := v.(int); ok {
			value = float64(iv)
		}
	}

	return &capacities.SetPointMessage{Value: value}, nil
}

