package meta

import (
	"fmt"
	"reflect"
)

// AttributeType 属性类型
type AttributeType string

const (
	AttrTypeString AttributeType = "string"
	AttrTypeInt    AttributeType = "int"
	AttrTypeFloat  AttributeType = "float"
	AttrTypeBool   AttributeType = "bool"
	AttrTypeTime   AttributeType = "time"
	AttrTypeObject AttributeType = "object"
	AttrTypeArray  AttributeType = "array"
)

// AttributeDescriptor 属性描述符
type AttributeDescriptor struct {
	Name        string
	Type        AttributeType
	Required    bool
	Key         bool
	Description string
	Default     interface{}
}

// RelationshipDescriptor 关系描述符
type RelationshipDescriptor struct {
	Name        string
	TargetType  string
	Cardinality string // one_to_one, one_to_many, many_to_one, many_to_many
	Inverse     string
	Description string
}

// OperationDescriptor 操作描述符
type OperationDescriptor struct {
	Name        string
	InputType   reflect.Type
	OutputType  reflect.Type
	Description string
	Validator   func(interface{}, interface{}) error                // ResourceHandle, operation args
	Handler     func(interface{}, interface{}) (interface{}, error) // ResourceHandle, operation args
}

// CapabilityDescriptor 能力描述符
type CapabilityDescriptor struct {
	Name        string
	Operations  []string
	Description string
	Validators  map[string]func(interface{}, interface{}) error                // ResourceHandle, operation args
	Handlers    map[string]func(interface{}, interface{}) (interface{}, error) // ResourceHandle, operation args
}

// LifecycleState 生命周期状态
type LifecycleState string

const (
	StateCreated  LifecycleState = "created"
	StateActive   LifecycleState = "active"
	StateInactive LifecycleState = "inactive"
	StateDeleted  LifecycleState = "deleted"
)

// LifecycleTransition 生命周期转换
type LifecycleTransition struct {
	From    LifecycleState
	To      LifecycleState
	Trigger string
}

// LifecycleDescriptor 生命周期描述符
type LifecycleDescriptor struct {
	States      []LifecycleState
	Transitions []LifecycleTransition
}

// AccessControlDescriptor 访问控制描述符
type AccessControlDescriptor struct {
	DefaultPermissions []string
	RequiredRoles      []string
}

// TypeDescriptor 类型描述符（类似内核的 OBJECT_TYPE）
type TypeDescriptor struct {
	Name          string
	BaseType      *TypeDescriptor
	Attributes    []AttributeDescriptor
	Relationships []RelationshipDescriptor
	Capabilities  []CapabilityDescriptor
	Lifecycle     LifecycleDescriptor
	AccessControl AccessControlDescriptor
	Operations    map[string]OperationDescriptor
}

// GetAttribute 获取属性描述符
func (td *TypeDescriptor) GetAttribute(name string) (*AttributeDescriptor, error) {
	// 先查找当前类型的属性
	for i := range td.Attributes {
		if td.Attributes[i].Name == name {
			return &td.Attributes[i], nil
		}
	}

	// 如果没找到，查找基类的属性
	if td.BaseType != nil {
		return td.BaseType.GetAttribute(name)
	}

	return nil, fmt.Errorf("attribute %s not found in type %s", name, td.Name)
}

// GetAllAttributes 获取所有属性（包括继承的属性）
func (td *TypeDescriptor) GetAllAttributes() []AttributeDescriptor {
	attrs := make(map[string]AttributeDescriptor)

	// 收集基类属性
	if td.BaseType != nil {
		baseAttrs := td.BaseType.GetAllAttributes()
		for _, attr := range baseAttrs {
			attrs[attr.Name] = attr
		}
	}

	// 覆盖当前类型的属性
	for _, attr := range td.Attributes {
		attrs[attr.Name] = attr
	}

	// 转换为切片
	result := make([]AttributeDescriptor, 0, len(attrs))
	for _, attr := range attrs {
		result = append(result, attr)
	}

	return result
}

// HasCapability 检查是否具有某种能力
func (td *TypeDescriptor) HasCapability(name string) bool {
	for _, cap := range td.Capabilities {
		if cap.Name == name {
			return true
		}
	}

	// 检查基类
	if td.BaseType != nil {
		return td.BaseType.HasCapability(name)
	}

	return false
}

// GetCapability 获取能力描述符
func (td *TypeDescriptor) GetCapability(name string) (*CapabilityDescriptor, error) {
	for i := range td.Capabilities {
		if td.Capabilities[i].Name == name {
			return &td.Capabilities[i], nil
		}
	}

	// 检查基类
	if td.BaseType != nil {
		return td.BaseType.GetCapability(name)
	}

	return nil, fmt.Errorf("capability %s not found in type %s", name, td.Name)
}

// GetAllCapabilities 获取所有能力（包括继承的能力）
func (td *TypeDescriptor) GetAllCapabilities() []CapabilityDescriptor {
	caps := make(map[string]CapabilityDescriptor)

	// 收集基类能力
	if td.BaseType != nil {
		baseCaps := td.BaseType.GetAllCapabilities()
		for _, cap := range baseCaps {
			caps[cap.Name] = cap
		}
	}

	// 覆盖当前类型的能力
	for _, cap := range td.Capabilities {
		caps[cap.Name] = cap
	}

	// 转换为切片
	result := make([]CapabilityDescriptor, 0, len(caps))
	for _, cap := range caps {
		result = append(result, cap)
	}

	return result
}

// IsSubtypeOf 检查是否是某个类型的子类型
func (td *TypeDescriptor) IsSubtypeOf(typeName string) bool {
	if td.Name == typeName {
		return true
	}

	if td.BaseType != nil {
		return td.BaseType.IsSubtypeOf(typeName)
	}

	return false
}

// ValidateOperation 验证操作是否支持
func (td *TypeDescriptor) ValidateOperation(capabilityName, operationName string) error {
	cap, err := td.GetCapability(capabilityName)
	if err != nil {
		return fmt.Errorf("capability %s not found: %w", capabilityName, err)
	}

	for _, op := range cap.Operations {
		if op == operationName {
			return nil
		}
	}

	return fmt.Errorf("operation %s not found in capability %s", operationName, capabilityName)
}
