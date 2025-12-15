package state

import (
	"fmt"
	"sync"
)

// StateBackend 状态后端接口（定义状态如何存储）
// 注意：Go 接口方法不能有类型参数，所以使用类型擦除的方式
type StateBackend interface {
	// GetValueState 获取值状态（通过类型断言获取具体类型）
	GetValueState(descriptor StateDescriptor) (interface{}, error)

	// GetListState 获取列表状态（通过类型断言获取具体类型）
	GetListState(descriptor StateDescriptor) (interface{}, error)

	// GetMapState 获取映射状态（通过类型断言获取具体类型）
	GetMapState(descriptor StateDescriptor) (interface{}, error)

	// Snapshot 创建状态快照
	Snapshot() (interface{}, error)

	// Restore 从快照恢复状态
	Restore(snapshot interface{}) error
}

// MemoryStateBackend 内存状态后端（用于测试和开发）
type MemoryStateBackend struct {
	mu        sync.RWMutex
	states    map[string]interface{} // 状态名称 -> 状态实现
	snapshots map[string]interface{} // 快照数据
}

// NewMemoryStateBackend 创建内存状态后端
func NewMemoryStateBackend() *MemoryStateBackend {
	return &MemoryStateBackend{
		states:    make(map[string]interface{}),
		snapshots: make(map[string]interface{}),
	}
}

// GetValueState 获取值状态（返回 interface{}，需要通过类型断言获取具体类型）
func (b *MemoryStateBackend) GetValueState(descriptor StateDescriptor) (interface{}, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	name := descriptor.Name()
	if state, exists := b.states[name]; exists {
		return state, nil
	}

	// 无法从 descriptor 推断类型，返回错误
	// 实际使用时应该通过类型特定的辅助函数创建
	return nil, fmt.Errorf("state %s not found, use typed helper functions", name)
}

// GetListState 获取列表状态（返回 interface{}，需要通过类型断言获取具体类型）
func (b *MemoryStateBackend) GetListState(descriptor StateDescriptor) (interface{}, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	name := descriptor.Name()
	if state, exists := b.states[name]; exists {
		return state, nil
	}

	return nil, fmt.Errorf("state %s not found, use typed helper functions", name)
}

// GetMapState 获取映射状态（返回 interface{}，需要通过类型断言获取具体类型）
func (b *MemoryStateBackend) GetMapState(descriptor StateDescriptor) (interface{}, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	name := descriptor.Name()
	if state, exists := b.states[name]; exists {
		return state, nil
	}

	return nil, fmt.Errorf("state %s not found, use typed helper functions", name)
}

// GetValueStateTyped 类型安全的获取值状态辅助函数（包级别函数）
func GetValueStateTyped[T any](backend *MemoryStateBackend, descriptor *ValueStateDescriptor[T]) (ValueState[T], error) {
	backend.mu.Lock()
	defer backend.mu.Unlock()

	name := descriptor.Name()
	if state, exists := backend.states[name]; exists {
		if valueState, ok := state.(ValueState[T]); ok {
			return valueState, nil
		}
		return nil, fmt.Errorf("state %s exists but has wrong type", name)
	}

	// 创建新的值状态
	valueState := &valueStateImpl[T]{
		name:  name,
		value: new(T),
		mu:    sync.RWMutex{},
	}
	backend.states[name] = valueState
	return valueState, nil
}

// GetListStateTyped 类型安全的获取列表状态辅助函数（包级别函数）
func GetListStateTyped[T any](backend *MemoryStateBackend, descriptor *ListStateDescriptor[T]) (ListState[T], error) {
	backend.mu.Lock()
	defer backend.mu.Unlock()

	name := descriptor.Name()
	if state, exists := backend.states[name]; exists {
		if listState, ok := state.(ListState[T]); ok {
			return listState, nil
		}
		return nil, fmt.Errorf("state %s exists but has wrong type", name)
	}

	// 创建新的列表状态
	listState := &listStateImpl[T]{
		name:  name,
		items: make([]T, 0),
		mu:    sync.RWMutex{},
	}
	backend.states[name] = listState
	return listState, nil
}

// GetMapStateTyped 类型安全的获取映射状态辅助函数（包级别函数）
func GetMapStateTyped[K comparable, V any](backend *MemoryStateBackend, descriptor *MapStateDescriptor[K, V]) (MapState[K, V], error) {
	backend.mu.Lock()
	defer backend.mu.Unlock()

	name := descriptor.Name()
	if state, exists := backend.states[name]; exists {
		if mapState, ok := state.(MapState[K, V]); ok {
			return mapState, nil
		}
		return nil, fmt.Errorf("state %s exists but has wrong type", name)
	}

	// 创建新的映射状态
	mapState := &mapStateImpl[K, V]{
		name: name,
		data: make(map[K]V),
		mu:   sync.RWMutex{},
	}
	backend.states[name] = mapState
	return mapState, nil
}

// Snapshot 创建状态快照
func (b *MemoryStateBackend) Snapshot() (interface{}, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// 序列化所有状态
	snapshot := make(map[string]interface{})
	for name, state := range b.states {
		// 这里简化处理，实际应该序列化状态数据
		snapshot[name] = state
	}

	return snapshot, nil
}

// Restore 从快照恢复状态
func (b *MemoryStateBackend) Restore(snapshot interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if snapshotMap, ok := snapshot.(map[string]interface{}); ok {
		b.states = snapshotMap
		return nil
	}

	return fmt.Errorf("invalid snapshot format")
}

// valueStateImpl 值状态实现
type valueStateImpl[T any] struct {
	name  string
	value *T
	mu    sync.RWMutex
}

func (s *valueStateImpl[T]) Value() (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return *s.value, nil
}

func (s *valueStateImpl[T]) Update(value T) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.value = &value
	return nil
}

func (s *valueStateImpl[T]) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.value = new(T)
	return nil
}

// listStateImpl 列表状态实现
type listStateImpl[T any] struct {
	name  string
	items []T
	mu    sync.RWMutex
}

func (s *listStateImpl[T]) Get() ([]T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]T, len(s.items))
	copy(result, s.items)
	return result, nil
}

func (s *listStateImpl[T]) Add(value T) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, value)
	return nil
}

func (s *listStateImpl[T]) AddAll(values []T) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, values...)
	return nil
}

func (s *listStateImpl[T]) Update(values []T) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make([]T, len(values))
	copy(s.items, values)
	return nil
}

func (s *listStateImpl[T]) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make([]T, 0)
	return nil
}

// mapStateImpl 映射状态实现
type mapStateImpl[K comparable, V any] struct {
	name string
	data map[K]V
	mu   sync.RWMutex
}

func (s *mapStateImpl[K, V]) Get(key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists := s.data[key]
	if !exists {
		var zero V
		return zero, fmt.Errorf("key not found")
	}
	return value, nil
}

func (s *mapStateImpl[K, V]) Put(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	return nil
}

func (s *mapStateImpl[K, V]) PutAll(entries map[K]V) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range entries {
		s.data[k] = v
	}
	return nil
}

func (s *mapStateImpl[K, V]) Remove(key K) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	return nil
}

func (s *mapStateImpl[K, V]) Contains(key K) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.data[key]
	return exists, nil
}

func (s *mapStateImpl[K, V]) Keys() ([]K, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]K, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func (s *mapStateImpl[K, V]) Values() ([]V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	values := make([]V, 0, len(s.data))
	for _, v := range s.data {
		values = append(values, v)
	}
	return values, nil
}

func (s *mapStateImpl[K, V]) Entries() (map[K]V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[K]V)
	for k, v := range s.data {
		result[k] = v
	}
	return result, nil
}

func (s *mapStateImpl[K, V]) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[K]V)
	return nil
}

