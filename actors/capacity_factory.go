package actors

import (
	"fmt"

	"github.com/uos-projects/uos-kernel/actors/capacities"
)

// CapacityFactory 根据 Control 类型创建对应的 Capacity
type CapacityFactory struct{}

// NewCapacityFactory 创建新的能力工厂
func NewCapacityFactory() *CapacityFactory {
	return &CapacityFactory{}
}

// CreateCapacity 根据 Control 类型创建对应的 Capacity
func (f *CapacityFactory) CreateCapacity(controlType string, controlID string) (capacities.Capacity, error) {
	switch controlType {
	case "AccumulatorReset":
		return capacities.NewAccumulatorResetCapacity(controlID), nil
	case "Command":
		return capacities.NewCommandCapacity(controlID), nil
	case "RaiseLowerCommand":
		return capacities.NewRaiseLowerCommandCapacity(controlID), nil
	case "SetPoint":
		return capacities.NewSetPointCapacity(controlID), nil
	default:
		return nil, fmt.Errorf("unknown control type: %s", controlType)
	}
}

// RegisterCapacity 注册自定义的 Capacity 创建函数（用于扩展）
type CapacityCreator func(controlID string) capacities.Capacity

var capacityRegistry = make(map[string]CapacityCreator)

// RegisterCapacityType 注册自定义的 Capacity 类型
func RegisterCapacityType(controlType string, creator CapacityCreator) {
	capacityRegistry[controlType] = creator
}

// CreateCapacityWithRegistry 使用注册表创建 Capacity（支持扩展）
func (f *CapacityFactory) CreateCapacityWithRegistry(controlType string, controlID string) (capacities.Capacity, error) {
	// 先检查注册表
	if creator, exists := capacityRegistry[controlType]; exists {
		return creator(controlID), nil
	}

	// 回退到默认实现
	return f.CreateCapacity(controlType, controlID)
}
