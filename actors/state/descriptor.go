package state

// StateDescriptor 状态描述符（定义状态的名称和类型）
type StateDescriptor interface {
	// Name 返回状态名称
	Name() string
}

// ValueStateDescriptor 值状态描述符
type ValueStateDescriptor[T any] struct {
	name string
}

// NewValueStateDescriptor 创建值状态描述符
func NewValueStateDescriptor[T any](name string) *ValueStateDescriptor[T] {
	return &ValueStateDescriptor[T]{name: name}
}

// Name 返回状态名称
func (d *ValueStateDescriptor[T]) Name() string {
	return d.name
}

// ListStateDescriptor 列表状态描述符
type ListStateDescriptor[T any] struct {
	name string
}

// NewListStateDescriptor 创建列表状态描述符
func NewListStateDescriptor[T any](name string) *ListStateDescriptor[T] {
	return &ListStateDescriptor[T]{name: name}
}

// Name 返回状态名称
func (d *ListStateDescriptor[T]) Name() string {
	return d.name
}

// MapStateDescriptor 映射状态描述符
type MapStateDescriptor[K comparable, V any] struct {
	name string
}

// NewMapStateDescriptor 创建映射状态描述符
func NewMapStateDescriptor[K comparable, V any](name string) *MapStateDescriptor[K, V] {
	return &MapStateDescriptor[K, V]{name: name}
}

// Name 返回状态名称
func (d *MapStateDescriptor[K, V]) Name() string {
	return d.name
}

