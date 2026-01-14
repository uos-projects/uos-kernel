package capacities

import (
	"reflect"
)

// CanHandleByCommandShape 辅助函数：使用 CommandShape 检查消息是否可处理
// 如果 capacity 实现了 CommandShape，则使用 CommandShape 进行类型检查
func CanHandleByCommandShape(capacity Capacity, msg Message) bool {
	if shape, ok := capacity.(CommandShape); ok {
		msgType := reflect.TypeOf(msg)
		for _, cmdType := range shape.AcceptableCommands() {
			if cmdType.MessageType == msgType {
				return true
			}
		}
		return false
	}
	return false
}
