package actors

import (
	"context"
	"time"

	"github.com/uos-projects/uos-kernel/actors/state"
)

// CIMResourceActor 代表一个 CIM 资源（数字孪生）
// 不体现 OWL 的层次关系，只维护 OWL 类的 URI reference
// 所有属性通过 map[string]interface{} 存储，与 OWL 定义一致（包括继承的属性）
type CIMResourceActor struct {
	*BaseResourceActor // 嵌入基础资源 Actor（能力管理）

	// OWL 类的完整 URI（用于语义解释）
	OWLClassURI string

	// 状态后端（用于持久化）
	stateBackend state.StateBackend

	// 快照存储（用于 Iceberg 持久化）
	snapshotStore state.SnapshotStore

	// 运行时上下文（供 Capacity 访问状态）
	runtimeContext *RuntimeContext

	// 属性映射（存储 OWL 定义的所有属性，包括继承的）
	// key: 属性名（CIM 属性名，如 "mRID", "nominalVoltage"）
	// value: 属性值
	properties map[string]interface{}
}

// NewCIMResourceActor 创建新的 CIM 资源 Actor
func NewCIMResourceActor(
	id string,
	owlClassURI string,
	behavior ActorBehavior,
) *CIMResourceActor {
	// 从 OWL URI 提取类名作为 resourceType
	localName := extractLocalName(owlClassURI)

	// 创建基础资源 Actor
	baseResourceActor := NewBaseResourceActor(id, localName, behavior)

	// 初始化状态后端
	stateBackend := state.NewMemoryStateBackend()
	runtimeContext := NewRuntimeContext(stateBackend)

	actor := &CIMResourceActor{
		BaseResourceActor: baseResourceActor,
		OWLClassURI:       owlClassURI,
		stateBackend:      stateBackend,
		snapshotStore:     state.NewMemorySnapshotStore(), // 默认使用内存快照存储
		runtimeContext:    runtimeContext,
		properties:        make(map[string]interface{}),
	}

	return actor
}

// GetOWLClassURI 获取 OWL 类的 URI
func (a *CIMResourceActor) GetOWLClassURI() string {
	return a.OWLClassURI
}

// SetProperty 设置属性值
func (a *CIMResourceActor) SetProperty(name string, value interface{}) {
	a.properties[name] = value
}

// GetProperty 获取属性值
func (a *CIMResourceActor) GetProperty(name string) (interface{}, bool) {
	value, exists := a.properties[name]
	return value, exists
}

// GetAllProperties 获取所有属性
func (a *CIMResourceActor) GetAllProperties() map[string]interface{} {
	// 返回副本，避免外部修改
	result := make(map[string]interface{})
	for k, v := range a.properties {
		result[k] = v
	}
	return result
}

// GetRuntimeContext 获取运行时上下文
func (a *CIMResourceActor) GetRuntimeContext() *RuntimeContext {
	return a.runtimeContext
}

// GetStateBackend 获取状态后端
func (a *CIMResourceActor) GetStateBackend() state.StateBackend {
	return a.stateBackend
}

// SetSnapshotStore 设置快照存储（例如设置为 IcebergBackend）
func (a *CIMResourceActor) SetSnapshotStore(store state.SnapshotStore) {
	a.snapshotStore = store
}

// SaveSnapshotAsync 异步保存快照到存储后端
func (a *CIMResourceActor) SaveSnapshotAsync(ctx context.Context) error {
	// 1. 创建内存快照
	// Sequence 暂时使用时间戳
	snapshot, err := a.CreateSnapshot(time.Now().UnixMilli())
	if err != nil {
		return err
	}

	// 2. 将 Snapshot 转换为 map（避免循环导入）
	snapshotMap := map[string]interface{}{
		"ActorID":     snapshot.ActorID,
		"OWLClassURI": snapshot.OWLClassURI,
		"Sequence":    snapshot.Sequence,
		"Timestamp":   snapshot.Timestamp,
		"Properties":  snapshot.Properties,
		"State":       snapshot.State,
	}

	// 3. 保存到后端
	// 注意：实际生产中这应该是一个异步任务，不阻塞 Actor 主线程
	// 但在这个简单的实现中，我们同步调用，或者启动 Goroutine
	// 为了演示清晰和确保完成，这里先用同步调用，或者 go func
	go func() {
		// 使用新的 Context 防止外部取消导致保存失败
		// 这里应该有更好的生命周期管理
		err := a.snapshotStore.Save(context.Background(), a.ResourceID(), snapshotMap)
		if err != nil {
			// Log error
			// fmt.Printf("Failed to save snapshot for %s: %v\n", a.ResourceID(), err)
		}
	}()
	return nil
}

// CreateSnapshot 创建状态快照
func (a *CIMResourceActor) CreateSnapshot(sequence int64) (*Snapshot, error) {
	stateSnapshot, err := a.stateBackend.Snapshot()
	if err != nil {
		return nil, err
	}

	return &Snapshot{
		ActorID:     a.ResourceID(),
		OWLClassURI: a.OWLClassURI,
		Sequence:    sequence,
		Timestamp:   time.Now(),
		Properties:  a.GetAllProperties(),
		State:       stateSnapshot,
	}, nil
}

// RestoreFromSnapshot 从快照恢复状态
func (a *CIMResourceActor) RestoreFromSnapshot(snapshot *Snapshot) error {
	// 恢复属性
	a.properties = snapshot.Properties

	// 恢复状态后端
	return a.stateBackend.Restore(snapshot.State)
}

// extractLocalName 从 OWL URI 提取本地名称
func extractLocalName(uri string) string {
	// 查找最后一个 '#' 或 '/'
	for i := len(uri) - 1; i >= 0; i-- {
		if uri[i] == '#' || uri[i] == '/' {
			return uri[i+1:]
		}
	}
	return uri
}
