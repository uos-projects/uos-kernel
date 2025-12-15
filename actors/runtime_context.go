package actors

import (
	"context"
	"fmt"

	"github.com/uos-projects/uos-kernel/actors/state"
)

// RuntimeContext 运行时上下文，提供状态访问 API
// 类似于 Flink 的 RuntimeContext，允许 Capacity 访问 Actor 管理的状态
// 注意：Go 方法不能有类型参数，所以使用包级别的类型安全函数
type RuntimeContext struct {
	stateBackend state.StateBackend
}

// NewRuntimeContext 创建新的运行时上下文
func NewRuntimeContext(stateBackend state.StateBackend) *RuntimeContext {
	return &RuntimeContext{
		stateBackend: stateBackend,
	}
}

// GetValueState 获取值状态（类型安全的包级别函数）
func GetValueState[T any](runtimeCtx *RuntimeContext, name string) (state.ValueState[T], error) {
	descriptor := state.NewValueStateDescriptor[T](name)

	// 使用类型安全的辅助函数
	if backend, ok := runtimeCtx.stateBackend.(*state.MemoryStateBackend); ok {
		return state.GetValueStateTyped(backend, descriptor)
	}

	// 回退到接口方法（需要类型断言）
	stateInterface, err := runtimeCtx.stateBackend.GetValueState(descriptor)
	if err != nil {
		var zero state.ValueState[T]
		return zero, err
	}

	valueState, ok := stateInterface.(state.ValueState[T])
	if !ok {
		var zero state.ValueState[T]
		return zero, fmt.Errorf("state %s has wrong type", name)
	}

	return valueState, nil
}

// GetListState 获取列表状态（类型安全的包级别函数）
func GetListState[T any](runtimeCtx *RuntimeContext, name string) (state.ListState[T], error) {
	descriptor := state.NewListStateDescriptor[T](name)

	// 使用类型安全的辅助函数
	if backend, ok := runtimeCtx.stateBackend.(*state.MemoryStateBackend); ok {
		return state.GetListStateTyped(backend, descriptor)
	}

	// 回退到接口方法（需要类型断言）
	stateInterface, err := runtimeCtx.stateBackend.GetListState(descriptor)
	if err != nil {
		var zero state.ListState[T]
		return zero, err
	}

	listState, ok := stateInterface.(state.ListState[T])
	if !ok {
		var zero state.ListState[T]
		return zero, fmt.Errorf("state %s has wrong type", name)
	}

	return listState, nil
}

// GetMapState 获取映射状态（类型安全的包级别函数）
func GetMapState[K comparable, V any](runtimeCtx *RuntimeContext, name string) (state.MapState[K, V], error) {
	descriptor := state.NewMapStateDescriptor[K, V](name)

	// 使用类型安全的辅助函数
	if backend, ok := runtimeCtx.stateBackend.(*state.MemoryStateBackend); ok {
		return state.GetMapStateTyped(backend, descriptor)
	}

	// 回退到接口方法（需要类型断言）
	stateInterface, err := runtimeCtx.stateBackend.GetMapState(descriptor)
	if err != nil {
		var zero state.MapState[K, V]
		return zero, err
	}

	mapState, ok := stateInterface.(state.MapState[K, V])
	if !ok {
		var zero state.MapState[K, V]
		return zero, fmt.Errorf("state %s has wrong type", name)
	}

	return mapState, nil
}

// contextKey 用于在 context.Context 中存储 RuntimeContext 的键
type contextKey struct{}

// WithRuntimeContext 将 RuntimeContext 注入到 context.Context 中
func WithRuntimeContext(ctx context.Context, runtimeCtx *RuntimeContext) context.Context {
	return context.WithValue(ctx, contextKey{}, runtimeCtx)
}

// GetRuntimeContextFromContext 从 context.Context 中获取 RuntimeContext
func GetRuntimeContextFromContext(ctx context.Context) *RuntimeContext {
	runtimeCtx, ok := ctx.Value(contextKey{}).(*RuntimeContext)
	if !ok {
		return nil
	}
	return runtimeCtx
}

// MustGetRuntimeContext 从 context.Context 中获取 RuntimeContext，如果不存在则 panic
func MustGetRuntimeContext(ctx context.Context) *RuntimeContext {
	runtimeCtx := GetRuntimeContextFromContext(ctx)
	if runtimeCtx == nil {
		panic(fmt.Errorf("runtime context not found in context"))
	}
	return runtimeCtx
}

