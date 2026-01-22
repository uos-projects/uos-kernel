package kernel

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/uos-projects/uos-kernel/actor"
)

// ResourceDescriptor 资源描述符（类似文件描述符）
type ResourceDescriptor int32

const (
	InvalidDescriptor ResourceDescriptor = -1
	MinDescriptor     ResourceDescriptor = 0
	MaxDescriptor     ResourceDescriptor = 0x7FFFFFFF
)

// Resource 资源（内部使用）
// 每个resourceID对应一个Resource，管理资源的引用计数和排他性
type Resource struct {
	resourceID string
	actor      actor.ResourceActor // 使用接口，支持 BaseResourceActor 及其子类
	refCount   int32                // 引用计数（有多少个描述符打开了这个资源）
	exclusive  bool                 // 是否排他性资源（true表示只能被一个描述符打开）
	mu         sync.RWMutex
}

// ResourceHandle 资源句柄（内部使用）
// 每个描述符对应一个Handle，指向一个Resource
type ResourceHandle struct {
	descriptor ResourceDescriptor
	resource   *Resource // 指向资源
}

// Manager 资源管理器（类似文件系统）
type Manager struct {
	system    *actor.System
	handles   map[ResourceDescriptor]*ResourceHandle // 描述符 -> Handle
	resources map[string]*Resource                   // resourceID -> Resource
	nextFD    int32                                  // 下一个可用的描述符
	mu        sync.RWMutex
}

// NewManager 创建资源管理器
func NewManager(system *actor.System) *Manager {
	return &Manager{
		system:    system,
		handles:   make(map[ResourceDescriptor]*ResourceHandle),
		resources: make(map[string]*Resource),
		nextFD:    int32(MinDescriptor),
	}
}

// Open 打开资源，返回资源描述符
// 每次Open都会分配新的描述符，但资源本身会维护引用计数
// 如果资源是排他性的且已被打开，则返回错误
func (rm *Manager) Open(resourceID string) (ResourceDescriptor, error) {
	return rm.OpenWithExclusive(resourceID, false)
}

// OpenWithExclusive 打开资源，可指定是否为排他性资源
func (rm *Manager) OpenWithExclusive(resourceID string, exclusive bool) (ResourceDescriptor, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// 检查资源是否存在
	a, exists := rm.system.Get(resourceID)
	if !exists {
		return InvalidDescriptor, fmt.Errorf("resource %s not found", resourceID)
	}

	// 尝试转换为 ResourceActor 接口
	resourceActor, ok := a.(actor.ResourceActor)
	if !ok {
		return InvalidDescriptor, fmt.Errorf("resource %s is not a ResourceActor", resourceID)
	}

	// 获取或创建Resource
	resource, exists := rm.resources[resourceID]
	if !exists {
		// 创建新资源
		resource = &Resource{
			resourceID: resourceID,
			actor:      resourceActor,
			refCount:   0,
			exclusive:  exclusive,
		}
		rm.resources[resourceID] = resource
	} else {
		// 资源已存在，检查排他性
		resource.mu.RLock()
		isExclusive := resource.exclusive || exclusive
		currentRefCount := resource.refCount
		resource.mu.RUnlock()

		if isExclusive && currentRefCount > 0 {
			return InvalidDescriptor, fmt.Errorf("resource %s is exclusive and already opened", resourceID)
		}

		// 更新排他性标志（如果新请求是排他性的）
		if exclusive {
			resource.mu.Lock()
			resource.exclusive = true
			resource.mu.Unlock()
		}
	}

	// 分配新的描述符（每次Open都分配新的）
	fd := ResourceDescriptor(atomic.AddInt32(&rm.nextFD, 1))
	if fd > MaxDescriptor {
		return InvalidDescriptor, fmt.Errorf("descriptor pool exhausted")
	}

	// 创建Handle
	handle := &ResourceHandle{
		descriptor: fd,
		resource:   resource,
	}

	rm.handles[fd] = handle

	// 增加Resource的引用计数
	atomic.AddInt32(&resource.refCount, 1)

	return fd, nil
}

// Close 关闭资源描述符
// 删除Handle，减少Resource的引用计数
// 当引用计数变为0时，排他性资源可以再次打开
// 类似 POSIX close()
func (rm *Manager) Close(fd ResourceDescriptor) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	handle, exists := rm.handles[fd]
	if !exists {
		return fmt.Errorf("invalid descriptor: %d", fd)
	}

	// 减少Resource的引用计数
	resource := handle.resource
	refCount := atomic.AddInt32(&resource.refCount, -1)

	// 删除Handle
	delete(rm.handles, fd)

	// 检查引用计数异常（不应该小于0）
	if refCount < 0 {
		// 引用计数异常，重置为0并记录警告
		atomic.StoreInt32(&resource.refCount, 0)
		// TODO: 添加日志记录异常情况
	}

	// 注意：当引用计数变为0时，Resource仍然保留在resources map中
	// 这样可以：
	// 1. 快速重新打开资源（无需重新查找Actor）
	// 2. 排他性资源在refCount=0时可以再次打开（Open检查refCount > 0）
	// 3. 如果后续需要清理，可以在资源管理器中添加清理策略

	return nil
}

// GetHandle 获取资源句柄（内部使用）
func (rm *Manager) GetHandle(fd ResourceDescriptor) (*ResourceHandle, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	handle, exists := rm.handles[fd]
	if !exists {
		return nil, fmt.Errorf("invalid descriptor: %d", fd)
	}

	return handle, nil
}

// GetResource 获取资源（内部使用）
func (rm *Manager) GetResource(resourceID string) (*Resource, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	resource, exists := rm.resources[resourceID]
	if !exists {
		return nil, fmt.Errorf("resource %s not found", resourceID)
	}

	return resource, nil
}
