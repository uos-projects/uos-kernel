package kernel

import (
	"fmt"
	"sync"
)

// TypeRegistry 类型注册表（类似内核的类型管理器）
type TypeRegistry struct {
	types map[string]*TypeDescriptor
	mu    sync.RWMutex
}

// NewTypeRegistry 创建类型注册表
func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{
		types: make(map[string]*TypeDescriptor),
	}
}

// Register 注册类型
func (tr *TypeRegistry) Register(descriptor *TypeDescriptor) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	// 检查是否已存在
	if _, exists := tr.types[descriptor.Name]; exists {
		return fmt.Errorf("type %s already registered", descriptor.Name)
	}

	// 如果存在基类，验证基类已注册
	if descriptor.BaseType != nil {
		baseName := descriptor.BaseType.Name
		if _, exists := tr.types[baseName]; !exists {
			return fmt.Errorf("base type %s not registered", baseName)
		}
		// 设置基类引用
		descriptor.BaseType = tr.types[baseName]
	}

	tr.types[descriptor.Name] = descriptor
	return nil
}

// Get 获取类型描述符
func (tr *TypeRegistry) Get(typeName string) (*TypeDescriptor, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	descriptor, exists := tr.types[typeName]
	if !exists {
		return nil, fmt.Errorf("type %s not found", typeName)
	}

	return descriptor, nil
}

// MustGet 获取类型描述符（不存在则panic）
func (tr *TypeRegistry) MustGet(typeName string) *TypeDescriptor {
	descriptor, err := tr.Get(typeName)
	if err != nil {
		panic(err)
	}
	return descriptor
}

// List 列出所有已注册的类型
func (tr *TypeRegistry) List() []string {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	names := make([]string, 0, len(tr.types))
	for name := range tr.types {
		names = append(names, name)
	}

	return names
}

// Exists 检查类型是否存在
func (tr *TypeRegistry) Exists(typeName string) bool {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	_, exists := tr.types[typeName]
	return exists
}

// Unregister 注销类型（谨慎使用）
func (tr *TypeRegistry) Unregister(typeName string) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	// 检查是否有其他类型继承此类型
	for name, desc := range tr.types {
		if desc.BaseType != nil && desc.BaseType.Name == typeName {
			return fmt.Errorf("cannot unregister type %s: type %s inherits from it", typeName, name)
		}
	}

	delete(tr.types, typeName)
	return nil
}
