package state

// ValueState 值状态接口（单个可变值）
type ValueState[T any] interface {
	// Value 获取当前值
	Value() (T, error)

	// Update 更新值
	Update(value T) error

	// Clear 清空值
	Clear() error
}

// ListState 列表状态接口（可变列表）
type ListState[T any] interface {
	// Get 获取所有元素
	Get() ([]T, error)

	// Add 添加元素
	Add(value T) error

	// AddAll 添加所有元素
	AddAll(values []T) error

	// Update 更新整个列表
	Update(values []T) error

	// Clear 清空列表
	Clear() error
}

// MapState 映射状态接口（键值映射）
type MapState[K comparable, V any] interface {
	// Get 获取指定键的值
	Get(key K) (V, error)

	// Put 设置键值对
	Put(key K, value V) error

	// PutAll 设置所有键值对
	PutAll(entries map[K]V) error

	// Remove 移除指定键
	Remove(key K) error

	// Contains 检查是否包含指定键
	Contains(key K) (bool, error)

	// Keys 获取所有键
	Keys() ([]K, error)

	// Values 获取所有值
	Values() ([]V, error)

	// Entries 获取所有键值对
	Entries() (map[K]V, error)

	// Clear 清空映射
	Clear() error
}

