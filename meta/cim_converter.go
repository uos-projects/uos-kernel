package meta

import (
	"fmt"
	"strings"
)

// CIMToTypeSystemConverter CIM到类型系统的转换器
type CIMToTypeSystemConverter struct {
	registry *TypeRegistry
}

// NewCIMToTypeSystemConverter 创建转换器
func NewCIMToTypeSystemConverter(registry *TypeRegistry) *CIMToTypeSystemConverter {
	return &CIMToTypeSystemConverter{
		registry: registry,
	}
}

// ConvertEntityToType 将CIM Entity转换为类型定义
// 这个方法需要从Python的Entity对象获取信息，这里提供一个接口
// 实际调用时，需要从ontology包获取Entity信息
func (c *CIMToTypeSystemConverter) ConvertEntityToType(entityName string, entityInfo EntityInfo) (*TypeDescriptor, error) {
	// 检查基类是否存在
	var baseType *TypeDescriptor
	if entityInfo.Inherits != "" {
		var err error
		baseType, err = c.registry.Get(entityInfo.Inherits)
		if err != nil {
			return nil, fmt.Errorf("base type %s not found: %w", entityInfo.Inherits, err)
		}
	}

	// 转换属性
	attrs := make([]AttributeDescriptor, len(entityInfo.Attributes))
	for i, attr := range entityInfo.Attributes {
		attrType, err := parseAttributeType(attr.Type)
		if err != nil {
			return nil, fmt.Errorf("invalid attribute type %s: %w", attr.Type, err)
		}
		attrs[i] = AttributeDescriptor{
			Name:        attr.Name,
			Type:        attrType,
			Required:    attr.Required,
			Key:         attr.Key,
			Description: attr.Description,
		}
	}

	// 转换关系
	relationships := make([]RelationshipDescriptor, len(entityInfo.Relationships))
	for i, rel := range entityInfo.Relationships {
		relationships[i] = RelationshipDescriptor{
			Name:        rel.Type,
			TargetType:  rel.Target,
			Cardinality: inferCardinality(rel),
			Inverse:     inferInverse(rel),
			Description: rel.Description,
		}
	}

	// 推断能力（从关系推断）
	capabilities := c.inferCapabilities(entityInfo)

	// 创建类型描述符
	descriptor := &TypeDescriptor{
		Name:          entityName,
		BaseType:      baseType,
		Attributes:    attrs,
		Relationships: relationships,
		Capabilities:  capabilities,
		Lifecycle: LifecycleDescriptor{
			States: []LifecycleState{StateCreated, StateActive, StateInactive, StateDeleted},
			Transitions: []LifecycleTransition{
				{From: StateCreated, To: StateActive, Trigger: "activate"},
				{From: StateActive, To: StateInactive, Trigger: "deactivate"},
				{From: StateInactive, To: StateActive, Trigger: "reactivate"},
				{From: StateActive, To: StateDeleted, Trigger: "delete"},
			},
		},
		AccessControl: AccessControlDescriptor{
			DefaultPermissions: []string{"read", "execute"},
			RequiredRoles:      []string{"operator", "engineer"},
		},
		Operations: make(map[string]OperationDescriptor),
	}

	return descriptor, nil
}

// EntityInfo CIM实体信息（从Python转换过来的结构）
type EntityInfo struct {
	Name          string
	Inherits      string
	Attributes    []AttributeInfo
	Relationships []RelationshipInfo
}

// AttributeInfo 属性信息
type AttributeInfo struct {
	Name        string
	Type        string
	Required    bool
	Key         bool
	Description string
}

// RelationshipInfo 关系信息
type RelationshipInfo struct {
	Type        string
	Target      string
	Description string
}

// inferCapabilities 从CIM关系推断能力
func (c *CIMToTypeSystemConverter) inferCapabilities(entityInfo EntityInfo) []CapabilityDescriptor {
	capabilities := make(map[string]CapabilityDescriptor)

	// 检查是否有Control关系
	if c.hasControlRelationship(entityInfo) {
		capabilities["Control"] = CapabilityDescriptor{
			Name:        "Control",
			Operations:  []string{"execute", "query", "cancel"},
			Description: "控制能力",
			Validators:  make(map[string]func(interface{}, interface{}) error),
			Handlers:    make(map[string]func(interface{}, interface{}) (interface{}, error)),
		}
	}

	// 检查是否有Measurement关系
	if c.hasMeasurementRelationship(entityInfo) {
		capabilities["Measurement"] = CapabilityDescriptor{
			Name:        "Measurement",
			Operations:  []string{"subscribe", "read", "unsubscribe"},
			Description: "测量能力",
			Validators:  make(map[string]func(interface{}, interface{}) error),
			Handlers:    make(map[string]func(interface{}, interface{}) (interface{}, error)),
		}
	}

	// 检查是否有Switch相关的关系（推断SwitchControl能力）
	if c.hasSwitchRelationship(entityInfo) {
		capabilities["SwitchControl"] = CapabilityDescriptor{
			Name:        "SwitchControl",
			Operations:  []string{"open", "close", "queryState"},
			Description: "开关控制能力",
			Validators:  make(map[string]func(interface{}, interface{}) error),
			Handlers:    make(map[string]func(interface{}, interface{}) (interface{}, error)),
		}
	}

	// 转换为切片
	result := make([]CapabilityDescriptor, 0, len(capabilities))
	for _, cap := range capabilities {
		result = append(result, cap)
	}

	return result
}

// hasControlRelationship 检查是否有Control关系
func (c *CIMToTypeSystemConverter) hasControlRelationship(entityInfo EntityInfo) bool {
	// 检查实体名称是否包含Control相关关键词
	// 或者检查关系是否指向Control类型
	for _, rel := range entityInfo.Relationships {
		if strings.Contains(strings.ToLower(rel.Target), "control") {
			return true
		}
	}
	return false
}

// hasMeasurementRelationship 检查是否有Measurement关系
func (c *CIMToTypeSystemConverter) hasMeasurementRelationship(entityInfo EntityInfo) bool {
	// 检查实体名称是否包含Measurement相关关键词
	// 或者检查关系是否指向Measurement类型
	for _, rel := range entityInfo.Relationships {
		if strings.Contains(strings.ToLower(rel.Target), "measurement") {
			return true
		}
	}
	return false
}

// hasSwitchRelationship 检查是否有Switch相关关系
func (c *CIMToTypeSystemConverter) hasSwitchRelationship(entityInfo EntityInfo) bool {
	// 检查实体名称
	if strings.Contains(strings.ToLower(entityInfo.Name), "switch") ||
		strings.Contains(strings.ToLower(entityInfo.Name), "breaker") {
		return true
	}

	// 检查关系
	for _, rel := range entityInfo.Relationships {
		if strings.Contains(strings.ToLower(rel.Target), "switch") ||
			strings.Contains(strings.ToLower(rel.Target), "breaker") {
			return true
		}
	}
	return false
}

// inferCardinality 推断关系的基数
func inferCardinality(rel RelationshipInfo) string {
	// 根据关系类型推断基数
	// 这里简化处理，实际应该从CIM模型获取
	switch strings.ToLower(rel.Type) {
	case "contains", "has":
		return "one_to_many"
	case "partof", "belongsto":
		return "many_to_one"
	case "connects", "links":
		return "many_to_many"
	default:
		return "one_to_one"
	}
}

// inferInverse 推断关系的反向关系
func inferInverse(rel RelationshipInfo) string {
	// 根据关系类型推断反向关系
	switch strings.ToLower(rel.Type) {
	case "contains":
		return "partOf"
	case "partof":
		return "contains"
	case "has":
		return "belongsTo"
	case "belongsto":
		return "has"
	default:
		return ""
	}
}
