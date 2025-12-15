package actors

import (
	"time"

	"github.com/uos-projects/uos-kernel/actors/state"
)

// CIMResourceActor 代表一个 CIM 资源（数字孪生）
// 不体现 OWL 的层次关系，只维护 OWL 类的 URI reference
// 所有属性通过 map[string]interface{} 存储，与 OWL 定义一致（包括继承的属性）
type CIMResourceActor struct {
	*PowerSystemResourceActor

	// OWL 类的完整 URI（用于语义解释）
	OWLClassURI string

	// 状态后端（用于持久化）
	stateBackend state.StateBackend

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
	resourceActor := NewPowerSystemResourceActor(id, localName, behavior)

	// 初始化状态后端
	stateBackend := state.NewMemoryStateBackend()
	runtimeContext := NewRuntimeContext(stateBackend)

	actor := &CIMResourceActor{
		PowerSystemResourceActor: resourceActor,
		OWLClassURI:              owlClassURI,
		stateBackend:             stateBackend,
		runtimeContext:           runtimeContext,
		properties:               make(map[string]interface{}),
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
