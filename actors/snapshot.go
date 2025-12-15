package actors

import (
	"time"
)

// Snapshot 表示 Actor 状态的快照
type Snapshot struct {
	// ActorID Actor 的唯一标识符
	ActorID string

	// OWLClassURI OWL 类的完整 URI（用于语义解释）
	OWLClassURI string

	// Sequence 快照序列号（用于排序和恢复）
	Sequence int64

	// Timestamp 快照创建时间
	Timestamp time.Time

	// Properties Actor 的属性映射（OWL 定义的属性）
	Properties map[string]interface{}

	// State 状态后端的状态快照（ValueState/ListState/MapState）
	State interface{}
}

