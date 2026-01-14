package actors

import (
	"time"
)

// ResourceSnapshot 表示资源 Actor 状态的快照（业务中立）
type ResourceSnapshot struct {
	// ActorID Actor 的唯一标识符
	ActorID string

	// Sequence 快照序列号（用于排序和恢复）
	Sequence int64

	// Timestamp 快照创建时间
	Timestamp time.Time

	// Properties Actor 的属性映射
	Properties map[string]interface{}

	// State 状态后端的状态快照
	State interface{}
}

// Snapshot 表示 Actor 状态的快照（已废弃，使用 ResourceSnapshot）
// 保留此类型以保持向后兼容，新代码应使用 ResourceSnapshot
// Deprecated: 使用 ResourceSnapshot 代替
type Snapshot = ResourceSnapshot
