package kernel

import (
	"context"
	"fmt"

	"github.com/uos-projects/uos-kernel/actors"
)

// ControlCommand 控制命令（类似 ioctl 的命令）
type ControlCommand int32

const (
	// General Commands
	CMD_GET_RESOURCE_INFO ControlCommand = iota + 0x1000
	CMD_LIST_CAPABILITIES                // 0x1001
	CMD_LIST_EVENTS                      // 0x1002

	// Control Commands (通用控制命令，通过 Capacity 系统处理)
	CMD_EXECUTE_CAPACITY // 0x1003 - 执行指定 Capacity 的命令
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
func (rm *Manager) RCtl(ctx context.Context, fd ResourceDescriptor, cmd ControlCommand, arg ControlArg) (interface{}, error) {
	handle, err := rm.GetHandle(fd)
	if err != nil {
		return nil, err
	}

	// 通过Resource访问Actor
	resource := handle.resource
	resource.mu.RLock()
	actor := resource.actor
	resource.mu.RUnlock()

	// 根据命令处理
	switch cmd {
	case CMD_GET_RESOURCE_INFO:
		// 查询类命令，直接返回结果
		return &ResourceInfo{
			ResourceID:   actor.ResourceID(),
			ResourceType: actor.ResourceType(),
		}, nil

	case CMD_LIST_CAPABILITIES:
		return actor.ListCapabilities(), nil

	case CMD_LIST_EVENTS:
		return actor.ListEvents(), nil

	case CMD_EXECUTE_CAPACITY:
		// 通用控制命令：通过 Capacity 名称和消息执行
		// arg 应该是 map[string]interface{}，包含：
		//   - "capacity": Capacity 名称
		//   - "message": 要发送的消息（actors.Message）
		req, ok := arg.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid argument type for CMD_EXECUTE_CAPACITY, expected map[string]interface{}")
		}

		// 获取消息
		msg, ok := req["message"].(actors.Message)
		if !ok {
			return nil, fmt.Errorf("message field is required and must implement actors.Message")
		}

		// 可选：验证 Capacity 是否存在
		if capacityName, ok := req["capacity"].(string); ok {
			if !actor.HasCapacity(capacityName) {
				return nil, fmt.Errorf("capacity %s not found", capacityName)
			}
		}

		// 通过 Actor.Send() 发送消息
		// Actor 的 Receive() 会自动路由到对应的 Capacity
		if !actor.Send(msg) {
			return nil, fmt.Errorf("failed to send control message: mailbox full")
		}
		return nil, nil

	case CMD_SYNC:
		// 强制同步状态到持久化层
		// 使用 PropertyHolder 接口（业务中立）
		if propActor, ok := actor.(interface {
			SaveSnapshotAsync(context.Context) error
		}); ok {
			if err := propActor.SaveSnapshotAsync(ctx); err != nil {
				return nil, fmt.Errorf("failed to sync snapshot: %w", err)
			}
			return nil, nil // API 返回成功，异步保存
		}
		return nil, fmt.Errorf("resource does not support sync")

	default:
		return nil, fmt.Errorf("unknown control command: %d", cmd)
	}
}
