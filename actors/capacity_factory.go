package actors

import (
	"fmt"

	"github.com/uos-projects/uos-kernel/actors/capacities"
	"github.com/uos-projects/uos-kernel/actors/mq"
)

// CapacityFactory 根据 Control 或 Measurement 类型创建对应的 Capacity
type CapacityFactory struct {
	mqConsumer mq.MQConsumer // MQ 消费者（用于创建 Measurement Capacity）
}

// NewCapacityFactory 创建新的能力工厂
func NewCapacityFactory() *CapacityFactory {
	return &CapacityFactory{
		mqConsumer: nil, // 默认不使用 MQ
	}
}

// NewCapacityFactoryWithMQ 创建带 MQ 消费者的能力工厂
func NewCapacityFactoryWithMQ(consumer mq.MQConsumer) *CapacityFactory {
	return &CapacityFactory{
		mqConsumer: consumer,
	}
}

// SetMQConsumer 设置 MQ 消费者
func (f *CapacityFactory) SetMQConsumer(consumer mq.MQConsumer) {
	f.mqConsumer = consumer
}

// CreateCapacity 根据类型创建对应的 Capacity
// controlType 可以是 Control 类型（如 "AccumulatorReset"）或 Measurement 类型（如 "AnalogMeasurement"）
// resourceID 可以是 ControlID 或 MeasurementID
func (f *CapacityFactory) CreateCapacity(controlType string, resourceID string) (capacities.Capacity, error) {
	switch controlType {
	case "AccumulatorReset":
		return capacities.NewAccumulatorResetCapacity(resourceID), nil
	case "Command":
		return capacities.NewCommandCapacity(resourceID), nil
	case "RaiseLowerCommand":
		return capacities.NewRaiseLowerCommandCapacity(resourceID), nil
	case "SetPoint":
		return capacities.NewSetPointCapacity(resourceID), nil
	case "AnalogMeasurement":
		if f.mqConsumer == nil {
			return nil, fmt.Errorf("MQ consumer is required for AnalogMeasurement capacity")
		}
		// 默认测量类型，可以通过额外参数指定
		return capacities.NewAnalogMeasurementCapacity(resourceID, "AnalogMeasurement", f.mqConsumer), nil
	case "AccumulatorMeasurement":
		if f.mqConsumer == nil {
			return nil, fmt.Errorf("MQ consumer is required for AccumulatorMeasurement capacity")
		}
		return capacities.NewAccumulatorMeasurementCapacity(resourceID, "AccumulatorMeasurement", f.mqConsumer), nil
	case "DiscreteMeasurement":
		if f.mqConsumer == nil {
			return nil, fmt.Errorf("MQ consumer is required for DiscreteMeasurement capacity")
		}
		return capacities.NewDiscreteMeasurementCapacity(resourceID, "DiscreteMeasurement", f.mqConsumer), nil
	case "StringMeasurement":
		if f.mqConsumer == nil {
			return nil, fmt.Errorf("MQ consumer is required for StringMeasurement capacity")
		}
		return capacities.NewStringMeasurementCapacity(resourceID, "StringMeasurement", f.mqConsumer), nil
	default:
		return nil, fmt.Errorf("unknown capacity type: %s", controlType)
	}
}

// CreateMeasurementCapacity 创建测量 Capacity（便捷方法）
func (f *CapacityFactory) CreateMeasurementCapacity(
	measurementType string,
	measurementID string,
	measurementSubType string, // 如 "ThreePhaseActivePower"
) (capacities.Capacity, error) {
	if f.mqConsumer == nil {
		return nil, fmt.Errorf("MQ consumer is required for Measurement capacity")
	}

	switch measurementType {
	case "AnalogMeasurement":
		return capacities.NewAnalogMeasurementCapacity(measurementID, measurementSubType, f.mqConsumer), nil
	case "AccumulatorMeasurement":
		return capacities.NewAccumulatorMeasurementCapacity(measurementID, measurementSubType, f.mqConsumer), nil
	case "DiscreteMeasurement":
		return capacities.NewDiscreteMeasurementCapacity(measurementID, measurementSubType, f.mqConsumer), nil
	case "StringMeasurement":
		return capacities.NewStringMeasurementCapacity(measurementID, measurementSubType, f.mqConsumer), nil
	default:
		return nil, fmt.Errorf("unknown measurement type: %s", measurementType)
	}
}

// RegisterCapacity 注册自定义的 Capacity 创建函数（用于扩展）
type CapacityCreator func(resourceID string) capacities.Capacity

var capacityRegistry = make(map[string]CapacityCreator)

// RegisterCapacityType 注册自定义的 Capacity 类型
func RegisterCapacityType(capacityType string, creator CapacityCreator) {
	capacityRegistry[capacityType] = creator
}

// CreateCapacityWithRegistry 使用注册表创建 Capacity（支持扩展）
func (f *CapacityFactory) CreateCapacityWithRegistry(capacityType string, resourceID string) (capacities.Capacity, error) {
	// 先检查注册表
	if creator, exists := capacityRegistry[capacityType]; exists {
		return creator(resourceID), nil
	}

	// 回退到默认实现
	return f.CreateCapacity(capacityType, resourceID)
}
