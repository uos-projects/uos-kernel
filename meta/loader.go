package meta

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// TypeSystemDefinition 类型系统定义（从YAML解析）
type TypeSystemDefinition struct {
	Version          string                 `yaml:"version"`
	Description      string                 `yaml:"description"`
	BaseTypes        []TypeDefinition       `yaml:"base_types"`
	ResourceTypes    []TypeDefinition       `yaml:"resource_types"`
	CapabilityTypes  []CapabilityDefinition `yaml:"capability_types"`
	OperationMapping OperationMappingDef    `yaml:"operation_mapping"`
}

// TypeDefinition 类型定义（YAML结构）
type TypeDefinition struct {
	Name          string                   `yaml:"name"`
	BaseType      string                   `yaml:"base_type"`
	Attributes    []AttributeDefinition    `yaml:"attributes"`
	Relationships []RelationshipDefinition `yaml:"relationships"`
	Capabilities  []CapabilityRef          `yaml:"capabilities"`
	Lifecycle     LifecycleDefinition      `yaml:"lifecycle"`
	AccessControl AccessControlDefinition  `yaml:"access_control"`
}

// AttributeDefinition 属性定义（YAML结构）
type AttributeDefinition struct {
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"`
	Required    bool        `yaml:"required"`
	Key         bool        `yaml:"key"`
	Description string      `yaml:"description"`
	Default     interface{} `yaml:"default"`
}

// RelationshipDefinition 关系定义（YAML结构）
type RelationshipDefinition struct {
	Name        string `yaml:"name"`
	TargetType  string `yaml:"target_type"`
	Cardinality string `yaml:"cardinality"`
	Inverse     string `yaml:"inverse"`
	Description string `yaml:"description"`
}

// CapabilityRef 能力引用（YAML结构）
type CapabilityRef struct {
	Name        string   `yaml:"name"`
	Operations  []string `yaml:"operations"`
	Description string   `yaml:"description"`
}

// CapabilityDefinition 能力定义（YAML结构）
type CapabilityDefinition struct {
	Name       string                `yaml:"name"`
	Operations []OperationDefinition `yaml:"operations"`
}

// OperationDefinition 操作定义（YAML结构）
type OperationDefinition struct {
	Name        string `yaml:"name"`
	InputType   string `yaml:"input_type"`
	OutputType  string `yaml:"output_type"`
	Description string `yaml:"description"`
}

// LifecycleDefinition 生命周期定义（YAML结构）
type LifecycleDefinition struct {
	States      []string               `yaml:"states"`
	Transitions []TransitionDefinition `yaml:"transitions"`
}

// TransitionDefinition 转换定义（YAML结构）
type TransitionDefinition struct {
	From    string `yaml:"from"`
	To      string `yaml:"to"`
	Trigger string `yaml:"trigger"`
}

// AccessControlDefinition 访问控制定义（YAML结构）
type AccessControlDefinition struct {
	DefaultPermissions []string `yaml:"default_permissions"`
	RequiredRoles      []string `yaml:"required_roles"`
}

// OperationMappingDef 操作映射定义（YAML结构）
type OperationMappingDef struct {
	IoctlCommands []IoctlCommandDef `yaml:"ioctl_commands"`
}

// IoctlCommandDef ioctl命令定义（YAML结构）
type IoctlCommandDef struct {
	IoctlCmd    int    `yaml:"ioctl_cmd"`
	Capability  string `yaml:"capability"`
	Operation   string `yaml:"operation"`
	MessageType string `yaml:"message_type"`
}

// LoadTypeSystem 从YAML文件加载类型系统
func LoadTypeSystem(filePath string) (*TypeRegistry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read type system file: %w", err)
	}

	var def TypeSystemDefinition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("failed to parse type system YAML: %w", err)
	}

	registry := NewTypeRegistry()

	// 先注册基础类型
	for _, baseType := range def.BaseTypes {
		descriptor, err := convertTypeDefinition(baseType, nil, registry)
		if err != nil {
			return nil, fmt.Errorf("failed to convert base type %s: %w", baseType.Name, err)
		}
		if err := registry.Register(descriptor); err != nil {
			return nil, fmt.Errorf("failed to register base type %s: %w", baseType.Name, err)
		}
	}

	// 注册资源类型（可能需要多轮，因为可能有继承关系）
	// 简单实现：假设YAML中已经按依赖顺序排列
	for _, resType := range def.ResourceTypes {
		var baseTypeDesc *TypeDescriptor
		if resType.BaseType != "" {
			baseTypeDesc, err = registry.Get(resType.BaseType)
			if err != nil {
				return nil, fmt.Errorf("base type %s not found for %s: %w", resType.BaseType, resType.Name, err)
			}
		}

		descriptor, err := convertTypeDefinition(resType, baseTypeDesc, registry)
		if err != nil {
			return nil, fmt.Errorf("failed to convert resource type %s: %w", resType.Name, err)
		}
		if err := registry.Register(descriptor); err != nil {
			return nil, fmt.Errorf("failed to register resource type %s: %w", resType.Name, err)
		}
	}

	return registry, nil
}

// convertTypeDefinition 将YAML定义转换为TypeDescriptor
func convertTypeDefinition(def TypeDefinition, baseType *TypeDescriptor, registry *TypeRegistry) (*TypeDescriptor, error) {
	// 转换属性
	attrs := make([]AttributeDescriptor, len(def.Attributes))
	for i, attrDef := range def.Attributes {
		attrType, err := parseAttributeType(attrDef.Type)
		if err != nil {
			return nil, fmt.Errorf("invalid attribute type %s: %w", attrDef.Type, err)
		}
		attrs[i] = AttributeDescriptor{
			Name:        attrDef.Name,
			Type:        attrType,
			Required:    attrDef.Required,
			Key:         attrDef.Key,
			Description: attrDef.Description,
			Default:     attrDef.Default,
		}
	}

	// 转换关系
	relationships := make([]RelationshipDescriptor, len(def.Relationships))
	for i := range def.Relationships {
		relDef := def.Relationships[i]
		relationships[i] = RelationshipDescriptor{
			Name:        relDef.Name,
			TargetType:  relDef.TargetType,
			Cardinality: relDef.Cardinality,
			Inverse:     relDef.Inverse,
			Description: relDef.Description,
		}
	}

	// 转换能力
	capabilities := make([]CapabilityDescriptor, len(def.Capabilities))
	for i, capRef := range def.Capabilities {
		capabilities[i] = CapabilityDescriptor{
			Name:        capRef.Name,
			Operations:  capRef.Operations,
			Description: capRef.Description,
			Validators:  make(map[string]func(interface{}, interface{}) error),
			Handlers:    make(map[string]func(interface{}, interface{}) (interface{}, error)),
		}
	}

	// 转换生命周期
	lifecycle := LifecycleDescriptor{
		States:      make([]LifecycleState, len(def.Lifecycle.States)),
		Transitions: make([]LifecycleTransition, len(def.Lifecycle.Transitions)),
	}
	for i, state := range def.Lifecycle.States {
		lifecycle.States[i] = LifecycleState(state)
	}
	for i, trans := range def.Lifecycle.Transitions {
		lifecycle.Transitions[i] = LifecycleTransition{
			From:    LifecycleState(trans.From),
			To:      LifecycleState(trans.To),
			Trigger: trans.Trigger,
		}
	}

	// 转换访问控制
	accessControl := AccessControlDescriptor{
		DefaultPermissions: def.AccessControl.DefaultPermissions,
		RequiredRoles:      def.AccessControl.RequiredRoles,
	}

	return &TypeDescriptor{
		Name:          def.Name,
		BaseType:      baseType,
		Attributes:    attrs,
		Relationships: relationships,
		Capabilities:  capabilities,
		Lifecycle:     lifecycle,
		AccessControl: accessControl,
		Operations:    make(map[string]OperationDescriptor),
	}, nil
}

// parseAttributeType 解析属性类型字符串
func parseAttributeType(typeStr string) (AttributeType, error) {
	switch typeStr {
	case "string":
		return AttrTypeString, nil
	case "int":
		return AttrTypeInt, nil
	case "float":
		return AttrTypeFloat, nil
	case "bool":
		return AttrTypeBool, nil
	case "time":
		return AttrTypeTime, nil
	case "object":
		return AttrTypeObject, nil
	case "array":
		return AttrTypeArray, nil
	default:
		return "", fmt.Errorf("unknown attribute type: %s", typeStr)
	}
}

// GetIoctlCommandMapping 获取ioctl命令映射（从类型系统定义）
func GetIoctlCommandMapping(filePath string) (map[int]IoctlCommandDef, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read type system file: %w", err)
	}

	var def TypeSystemDefinition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("failed to parse type system YAML: %w", err)
	}

	mapping := make(map[int]IoctlCommandDef)
	for _, cmd := range def.OperationMapping.IoctlCommands {
		mapping[cmd.IoctlCmd] = cmd
	}

	return mapping, nil
}
