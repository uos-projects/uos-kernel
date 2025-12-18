package resource

import (
	"context"
	"fmt"
	"github.com/uos-projects/uos-kernel/actors"
	"github.com/uos-projects/uos-kernel/kernel"
)

// Open flags
const (
	O_RDONLY = 0x0000
	O_WRONLY = 0x0001
	O_RDWR   = 0x0002
	O_CREAT  = 0x0200
	O_EXCL   = 0x0800
)

// ResourceKernel 资源内核（类似操作系统内核）
// 面向用户的高级接口，提供类型验证和POSIX风格的系统调用
type ResourceKernel struct {
	system       *actors.System // 保存 System 引用以支持创建 Actor
	typeRegistry *kernel.TypeRegistry
	resourceMgr  *ResourceManager
	ioctlMapping map[int]kernel.IoctlCommandDef
}

// NewResourceKernel 创建资源内核
func NewResourceKernel(system *actors.System) *ResourceKernel {
	return &ResourceKernel{
		system:       system,
		typeRegistry: kernel.NewTypeRegistry(),
		resourceMgr:  NewResourceManager(system),
		ioctlMapping: make(map[int]kernel.IoctlCommandDef),
	}
}

// LoadTypeSystem 加载类型系统定义
func (k *ResourceKernel) LoadTypeSystem(filePath string) error {
	registry, err := kernel.LoadTypeSystem(filePath)
	if err != nil {
		return fmt.Errorf("failed to load type system: %w", err)
	}
	k.typeRegistry = registry

	// 加载ioctl命令映射
	mapping, err := kernel.GetIoctlCommandMapping(filePath)
	if err != nil {
		return fmt.Errorf("failed to load ioctl mapping: %w", err)
	}
	k.ioctlMapping = mapping

	return nil
}

// Open 打开资源（类似POSIX open）
func (k *ResourceKernel) Open(resourceType string, resourceID string, flags int) (ResourceDescriptor, error) {
	// 验证资源类型是否存在
	var typeDesc *kernel.TypeDescriptor
	if k.typeRegistry.Exists(resourceType) {
		typeDesc, _ = k.typeRegistry.Get(resourceType)
	} else {
		return InvalidDescriptor, fmt.Errorf("resource type %s not found", resourceType)
	}

	// 尝试通过ResourceManager打开资源
	fd, err := k.resourceMgr.Open(resourceID)
	if err == nil {
		// 资源已存在
		if (flags & O_CREAT) != 0 && (flags & O_EXCL) != 0 {
			k.Close(fd) // 关闭刚才打开的
			return InvalidDescriptor, fmt.Errorf("resource %s already exists", resourceID)
		}

		// [安全增强] 检查已存在的资源类型是否匹配
		state, err := k.resourceMgr.Read(context.Background(), fd)
		if err != nil {
			k.Close(fd)
			return InvalidDescriptor, fmt.Errorf("failed to read resource state: %w", err)
		}

		if state.ResourceType != resourceType {
			k.Close(fd)
			return InvalidDescriptor, fmt.Errorf("resource type mismatch: expected %s, got %s", resourceType, state.ResourceType)
		}

		return fd, nil
	}

	// 资源不存在
	// 如果设置了 O_CREAT，则尝试创建资源
	// 为了响应用户的特定请求，如果 flags == 0 (默认)，我们也尝试创建（或者我们可以要求用户使用 O_CREAT）
	// 这里我们遵循 POSIX 惯例，要求 O_CREAT。但为了方便演示用户请求的场景，我们可以暂时放宽，
	// 或者最好是引导用户加上 O_CREAT。这里我实现 O_CREAT 逻辑。
	if (flags & O_CREAT) != 0 {
		if err := k.createResource(resourceType, resourceID, typeDesc); err != nil {
			return InvalidDescriptor, fmt.Errorf("failed to create resource: %w", err)
		}

		// 创建成功后再次尝试打开
		return k.resourceMgr.Open(resourceID)
	}

	return InvalidDescriptor, err
}

// createResource 创建新资源
func (k *ResourceKernel) createResource(resourceType, resourceID string, typeDesc *kernel.TypeDescriptor) error {
	// 1. 构造 OWL Class URI
	// 假设 URI 格式为 CIM 命名空间 + 类型名
	owlClassURI := "http://www.iec.ch/TC57/CIM#" + resourceType

	// 2. 创建 Actor
	// 使用 CIMResourceActor
	actor := actors.NewCIMResourceActor(resourceID, owlClassURI, nil)

	// 3. 根据类型定义添加 Capabilities
	factory := actors.NewCapacityFactory()
	
	// 获取所有能力（包括继承的）
	capabilities := typeDesc.GetAllCapabilities()
	
	for _, capDesc := range capabilities {
		if capDesc.Name == "Control" || capDesc.Name == "SwitchControl" {
			// 默认添加 CommandCapacity，这是最通用的
			// ID 命名习惯: resourceID + "-command"
			cmdCap, _ := factory.CreateCapacity("Command", resourceID+"-command")
			if cmdCap != nil {
				actor.AddCapacity(cmdCap)
			}
		}
	}

	// 4. 注册到 System
	if err := k.system.Register(actor); err != nil {
		return err
	}

	return nil
}

// Close 关闭资源（类似POSIX close）
func (k *ResourceKernel) Close(fd ResourceDescriptor) error {
	return k.resourceMgr.Close(fd)
}

// Read 读取资源状态（类似POSIX read）
func (k *ResourceKernel) Read(ctx context.Context, fd ResourceDescriptor) (*ActorState, error) {
	return k.resourceMgr.Read(ctx, fd)
}

// Write 写入资源状态（类似POSIX write）
func (k *ResourceKernel) Write(ctx context.Context, fd ResourceDescriptor, req *WriteRequest) error {
	return k.resourceMgr.Write(ctx, fd, req)
}

// Stat 查询资源信息（类似POSIX stat）
func (k *ResourceKernel) Stat(ctx context.Context, fd ResourceDescriptor) (*ResourceStat, error) {
	// 通过Read获取资源信息
	state, err := k.resourceMgr.Read(ctx, fd)
	if err != nil {
		return nil, err
	}

	// 获取资源类型
	resourceType := state.ResourceType

	// 从类型注册表获取类型描述符
	var typeDesc *kernel.TypeDescriptor
	if k.typeRegistry.Exists(resourceType) {
		var err error
		typeDesc, err = k.typeRegistry.Get(resourceType)
		if err != nil {
			// 如果获取失败，typeDesc保持为nil
			typeDesc = nil
		}
	}

	// 获取能力列表
	actorCapabilities := state.Capabilities

	// 构建能力描述符列表
	capDescs := make([]CapabilityInfo, len(actorCapabilities))
	for i, capName := range actorCapabilities {
		var capDesc *kernel.CapabilityDescriptor
		if typeDesc != nil {
			cap, _ := typeDesc.GetCapability(capName)
			if cap != nil {
				capDesc = cap
			}
		}

		capDescs[i] = CapabilityInfo{
			Name:        capName,
			Operations:  getOperations(capDesc),
			Description: getDescription(capDesc),
		}
	}

	return &ResourceStat{
		ResourceID:     state.ResourceID,
		ResourceType:   resourceType,
		TypeDescriptor: typeDesc,
		Capabilities:   capDescs,
	}, nil
}

// Ioctl 控制操作（类似POSIX ioctl）
func (k *ResourceKernel) Ioctl(ctx context.Context, fd ResourceDescriptor, request int, argp interface{}) (interface{}, error) {
	// 系统级命令处理（绕过 Capability 类型检查）
	if ControlCommand(request) == CMD_SYNC {
		result, err := k.resourceMgr.RCtl(ctx, fd, ControlCommand(request), argp)
		if err != nil {
			return nil, fmt.Errorf("ioctl sync failed: %w", err)
		}
		return result, nil
	}

	// 查找ioctl命令映射
	cmdDef, exists := k.ioctlMapping[request]
	if !exists {
		return nil, fmt.Errorf("unknown ioctl command: 0x%x", request)
	}

	// 通过Read获取资源信息
	state, err := k.resourceMgr.Read(ctx, fd)
	if err != nil {
		return nil, err
	}

	// 获取资源类型
	resourceType := state.ResourceType

	// 从类型注册表获取类型描述符
	var typeDesc *kernel.TypeDescriptor
	if k.typeRegistry.Exists(resourceType) {
		typeDesc, _ = k.typeRegistry.Get(resourceType)
	}

	// 验证操作
	if typeDesc != nil {
		if err := typeDesc.ValidateOperation(cmdDef.Capability, cmdDef.Operation); err != nil {
			return nil, fmt.Errorf("operation validation failed: %w", err)
		}
	}

	// 通过ResourceManager的RCtl执行操作
	// 需要将ioctl命令转换为ControlCommand
	controlCmd := convertIoctlToControlCommand(request)

	result, err := k.resourceMgr.RCtl(ctx, fd, controlCmd, argp)
	if err != nil {
		// 错误处理优化: 转换为更友好的错误信息
		return nil, fmt.Errorf("ioctl operation failed: %w", err)
	}

	return result, nil
}

// Find 查找资源（类似POSIX find）
func (k *ResourceKernel) Find(resourceType string, filter Filter) ([]ResourceDescriptor, error) {
	// 验证资源类型
	if !k.typeRegistry.Exists(resourceType) {
		return nil, fmt.Errorf("resource type %s not found", resourceType)
	}

	// 从Actor系统查找资源
	// 这里需要扩展ResourceManager或System来支持按类型查找
	// 暂时返回空列表
	return []ResourceDescriptor{}, nil
}

// Watch 监听资源变化（类似inotify）
func (k *ResourceKernel) Watch(fd ResourceDescriptor, events []EventType) (<-chan Event, error) {
	// TODO: 实现资源变化监听
	return nil, fmt.Errorf("watch not implemented yet")
}

// ResourceStat 资源统计信息
type ResourceStat struct {
	ResourceID     string
	ResourceType   string
	TypeDescriptor *kernel.TypeDescriptor
	Capabilities   []CapabilityInfo
}

// CapabilityInfo 能力信息
type CapabilityInfo struct {
	Name        string
	Operations  []string
	Description string
}

// Filter 资源过滤器
type Filter struct {
	Type      string
	Attribute string
	Value     interface{}
}

// EventType 事件类型
type EventType string

const (
	EventStateChange      EventType = "state_change"
	EventAttributeChange  EventType = "attribute_change"
	EventCapabilityChange EventType = "capability_change"
)

// Event 资源事件
type Event struct {
	Type       EventType
	ResourceID string
	Data       interface{}
}

// 辅助函数
func getOperations(cap *kernel.CapabilityDescriptor) []string {
	if cap == nil {
		return []string{}
	}
	return cap.Operations
}

func getDescription(cap *kernel.CapabilityDescriptor) string {
	if cap == nil {
		return ""
	}
	return cap.Description
}

// convertIoctlToControlCommand 将ioctl命令转换为ControlCommand
func convertIoctlToControlCommand(ioctlCmd int) ControlCommand {
	// 映射ioctl命令到ControlCommand
	// 这里需要根据实际的ioctl命令定义进行映射
	switch ioctlCmd {
	case 0x1001:
		return CMD_ACCUMULATOR_RESET
	case 0x1002:
		return CMD_COMMAND
	case 0x1003:
		return CMD_RAISE_LOWER
	case 0x1004:
		return CMD_SET_POINT
	default:
		return ControlCommand(ioctlCmd)
	}
}


