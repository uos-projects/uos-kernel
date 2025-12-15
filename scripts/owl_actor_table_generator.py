#!/usr/bin/env python3
"""
从 OWL 模型生成：
1. 通用的 CIMResourceActor（不体现层次关系，只维护 OWL URI reference）
2. 每个 OWL 类的独立 Iceberg 表（包含所有继承的属性）
"""

import xml.etree.ElementTree as ET
from pathlib import Path
from typing import List, Dict, Optional, Set
from dataclasses import dataclass, field
import re

NSMAP = {
    'owl': 'http://www.w3.org/2002/07/owl#',
    'rdf': 'http://www.w3.org/1999/02/22-rdf-syntax-ns#',
    'rdfs': 'http://www.w3.org/2000/01/rdf-schema#',
    'xsd': 'http://www.w3.org/2001/XMLSchema#',
    'cim': 'http://www.iec.ch/TC57/CIM#',
}

CIM_NS = 'http://www.iec.ch/TC57/CIM#'

@dataclass
class Property:
    """属性定义"""
    uri: str
    local_name: str
    type: str  # ObjectProperty 或 DatatypeProperty
    range_uri: Optional[str]  # 值域 URI
    domain_uri: Optional[str]  # 定义域 URI
    comment: str
    xsd_type: Optional[str] = None  # 如果是 DatatypeProperty，XSD 类型

@dataclass
class OWLClass:
    """OWL 类定义"""
    uri: str
    local_name: str
    comment: str
    super_classes: List[str] = field(default_factory=list)  # 父类 URI 列表
    all_properties: List[Property] = field(default_factory=list)  # 所有属性（包括继承的）
    is_power_system_resource: bool = False
    is_control: bool = False
    is_measurement: bool = False

class OWLParser:
    """OWL 解析器 - 收集所有属性（包括继承的）"""
    
    def __init__(self, owl_file: str):
        self.owl_file = owl_file
        self.tree = ET.parse(owl_file)
        self.root = self.tree.getroot()
        self.classes: Dict[str, OWLClass] = {}
        self.properties: Dict[str, Property] = {}
        
    def parse(self):
        """解析 OWL 文件"""
        # 1. 解析所有类
        self._parse_classes()
        
        # 2. 解析所有属性
        self._parse_properties()
        
        # 3. 构建继承关系并收集所有属性
        self._build_inheritance_and_properties()
        
    def _parse_classes(self):
        """解析类定义"""
        classes = self.root.findall('.//{http://www.w3.org/2002/07/owl#}Class', NSMAP)
        
        for cls_elem in classes:
            about = cls_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}about')
            if not about:
                continue
                
            local_name = self._get_local_name(about)
            
            # 获取注释
            comment_elem = cls_elem.find('.//{http://www.w3.org/2000/01/rdf-schema#}comment')
            comment = comment_elem.text if comment_elem is not None else ""
            
            # 获取父类
            super_classes = []
            for sub_elem in cls_elem.findall('.//{http://www.w3.org/2000/01/rdf-schema#}subClassOf'):
                resource = sub_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}resource')
                if resource:
                    super_classes.append(resource)
            
            self.classes[about] = OWLClass(
                uri=about,
                local_name=local_name,
                comment=comment,
                super_classes=super_classes,
            )
    
    def _parse_properties(self):
        """解析属性定义"""
        # 解析 ObjectProperty
        obj_props = self.root.findall('.//{http://www.w3.org/2002/07/owl#}ObjectProperty', NSMAP)
        for prop_elem in obj_props:
            about = prop_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}about')
            if not about:
                continue
                
            local_name = self._get_local_name(about)
            
            # 获取 domain 和 range
            domain_elem = prop_elem.find('.//{http://www.w3.org/2000/01/rdf-schema#}domain')
            range_elem = prop_elem.find('.//{http://www.w3.org/2000/01/rdf-schema#}range')
            
            domain = domain_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}resource') if domain_elem is not None else None
            range_uri = range_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}resource') if range_elem is not None else None
            
            comment_elem = prop_elem.find('.//{http://www.w3.org/2000/01/rdf-schema#}comment')
            comment = comment_elem.text if comment_elem is not None else ""
            
            self.properties[about] = Property(
                uri=about,
                local_name=local_name,
                type="ObjectProperty",
                range_uri=range_uri,
                domain_uri=domain,
                comment=comment,
            )
        
        # 解析 DatatypeProperty
        datatype_props = self.root.findall('.//{http://www.w3.org/2002/07/owl#}DatatypeProperty', NSMAP)
        for prop_elem in datatype_props:
            about = prop_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}about')
            if not about:
                continue
                
            local_name = self._get_local_name(about)
            
            domain_elem = prop_elem.find('.//{http://www.w3.org/2000/01/rdf-schema#}domain')
            range_elem = prop_elem.find('.//{http://www.w3.org/2000/01/rdf-schema#}range')
            
            domain = domain_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}resource') if domain_elem is not None else None
            range_uri = range_elem.get('{http://www.w3.org/1999/02/22/rdf-syntax-ns#}resource') if range_elem is not None else None
            
            # 提取 XSD 类型
            xsd_type = None
            if range_uri and 'XMLSchema#' in range_uri:
                xsd_type = range_uri.split('#')[-1]
            
            comment_elem = prop_elem.find('.//{http://www.w3.org/2000/01/rdf-schema#}comment')
            comment = comment_elem.text if comment_elem is not None else ""
            
            self.properties[about] = Property(
                uri=about,
                local_name=local_name,
                type="DatatypeProperty",
                range_uri=range_uri,
                domain_uri=domain,
                comment=comment,
                xsd_type=xsd_type,
            )
    
    def _build_inheritance_and_properties(self):
        """构建继承关系并收集所有属性（包括继承的）"""
        # 为每个类收集所有属性（包括继承的）
        for cls_uri, cls_def in self.classes.items():
            # 收集所有属性（递归收集父类的属性）
            all_props = self._collect_all_properties(cls_uri, set())
            cls_def.all_properties = all_props
            
            # 判断是否继承自 PowerSystemResource/Control/Measurement
            cls_def.is_power_system_resource = self._is_power_system_resource(cls_uri)
            cls_def.is_control = self._is_control(cls_uri)
            cls_def.is_measurement = self._is_measurement(cls_uri)
    
    def _collect_all_properties(self, cls_uri: str, visited: Set[str]) -> List[Property]:
        """递归收集类的所有属性（包括继承的）"""
        if cls_uri in visited:
            return []
        visited.add(cls_uri)
        
        properties = []
        
        # 收集当前类的属性
        for prop_uri, prop in self.properties.items():
            if prop.domain_uri == cls_uri:
                properties.append(prop)
        
        # 递归收集父类的属性
        if cls_uri in self.classes:
            for super_uri in self.classes[cls_uri].super_classes:
                super_props = self._collect_all_properties(super_uri, visited)
                properties.extend(super_props)
        
        return properties
    
    def _is_power_system_resource(self, uri: str) -> bool:
        """判断是否继承自 PowerSystemResource"""
        power_system_resource_uri = f"{CIM_NS}PowerSystemResource"
        
        if uri == power_system_resource_uri:
            return True
        
        if uri not in self.classes:
            return False
        
        # 检查父类
        for super_uri in self.classes[uri].super_classes:
            if self._is_power_system_resource(super_uri):
                return True
        
        return False
    
    def _is_control(self, uri: str) -> bool:
        """判断是否继承自 Control"""
        control_uri = f"{CIM_NS}Control"
        
        if uri == control_uri:
            return True
        
        if uri not in self.classes:
            return False
        
        for super_uri in self.classes[uri].super_classes:
            if self._is_control(super_uri):
                return True
        
        return False
    
    def _is_measurement(self, uri: str) -> bool:
        """判断是否继承自 Measurement"""
        measurement_uri = f"{CIM_NS}Measurement"
        
        if uri == measurement_uri:
            return True
        
        if uri not in self.classes:
            return False
        
        for super_uri in self.classes[uri].super_classes:
            if self._is_measurement(super_uri):
                return True
        
        return False
    
    def _get_local_name(self, uri: str) -> str:
        """从 URI 提取本地名称"""
        if '#' in uri:
            return uri.split('#')[-1]
        elif '/' in uri:
            return uri.split('/')[-1]
        return uri
    
    def get_power_system_resources(self) -> List[OWLClass]:
        """获取所有 PowerSystemResource 子类"""
        return [cls for cls in self.classes.values() if cls.is_power_system_resource]
    
    def get_all_classes(self) -> List[OWLClass]:
        """获取所有类（用于生成表）"""
        return list(self.classes.values())


class ActorCodeGenerator:
    """Actor 代码生成器 - 生成通用的 CIMResourceActor"""
    
    def __init__(self, output_file: str = "actors/cim_resource_actor.go"):
        self.output_file = Path(output_file)
        self.output_file.parent.mkdir(parents=True, exist_ok=True)
    
    def generate(self):
        """生成通用的 CIMResourceActor"""
        code = """package actors

import (
	"context"
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
	runtimeContext RuntimeContext

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
func (a *CIMResourceActor) GetRuntimeContext() RuntimeContext {
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
"""
        self.output_file.write_text(code, encoding='utf-8')
        print(f"Generated: {self.output_file}")


class IcebergTableGenerator:
    """Iceberg 表生成器 - 为每个 OWL 类生成独立的表"""
    
    def __init__(self, output_file: str = "scripts/generated_iceberg_tables.sql"):
        self.output_file = Path(output_file)
        self.output_file.parent.mkdir(parents=True, exist_ok=True)
    
    def generate_table_sql(self, cls: OWLClass) -> str:
        """为单个类生成 Snapshot 表 SQL"""
        table_name = self._to_snake_case(cls.local_name)
        namespace = self._get_namespace(cls)
        
        # 基础字段（所有 Actor 快照表都有）
        base_columns = [
            "actor_id STRING NOT NULL",
            "owl_class_uri STRING NOT NULL",
            "sequence BIGINT NOT NULL",
            "timestamp TIMESTAMP NOT NULL",
            "snapshot_id STRING",
        ]
        
        # 状态字段（从类的所有属性生成，包括继承的）
        state_columns = []
        seen_props = set()  # 避免重复属性
        
        for prop in cls.all_properties:
            if prop.type == 'DatatypeProperty':
                # 提取纯属性名（去掉类名前缀，如 Equipment.aggregate -> aggregate）
                prop_name = self._extract_property_name(prop.local_name)
                
                # 避免重复（如果父类已有同名属性）
                if prop_name in seen_props:
                    continue
                seen_props.add(prop_name)
                
                iceberg_type = self._map_to_iceberg_type(prop.xsd_type)
                col_name = self._to_snake_case(prop_name)
                comment = prop.comment.replace("'", "''") if prop.comment else ""
                state_columns.append(f"    {col_name} {iceberg_type} COMMENT '{comment}'")
        
        # 时间旅行字段
        time_travel_columns = [
            "valid_from TIMESTAMP",
            "valid_to TIMESTAMP",
            "op_type STRING",
            "ingestion_ts TIMESTAMP",
        ]
        
        all_columns = base_columns + state_columns + time_travel_columns
        
        # 生成父类列表（用于注释）
        parent_names = [self._get_local_name(uri) for uri in cls.super_classes[:3]]
        parent_list = ", ".join(parent_names) if parent_names else "none"
        
        sql = f"""CREATE TABLE IF NOT EXISTS ontology.{namespace}.{table_name}_snapshots (
{',\n'.join(all_columns)}
)
USING ICEBERG
PARTITIONED BY (actor_id)
COMMENT 'Snapshot table for {cls.local_name} Actor state persistence (includes inherited properties from: {parent_list})'
"""
        return sql
    
    def generate_all(self, parser: OWLParser):
        """生成所有表的 SQL"""
        sqls = []
        
        # 1. 创建命名空间
        sqls.append("CREATE NAMESPACE IF NOT EXISTS ontology.grid")
        sqls.append("CREATE NAMESPACE IF NOT EXISTS ontology.assets")
        sqls.append("CREATE NAMESPACE IF NOT EXISTS ontology.metering")
        sqls.append("CREATE NAMESPACE IF NOT EXISTS ontology.actors")
        
        # 2. 为每个 PowerSystemResource 生成快照表
        power_system_resources = parser.get_power_system_resources()
        print(f"  Generating tables for {len(power_system_resources)} PowerSystemResource classes...")
        
        for cls in power_system_resources:
            sqls.append(self.generate_table_sql(cls))
        
        # 写入文件
        self.output_file.write_text("\n\n".join(sqls), encoding='utf-8')
        print(f"Generated Iceberg table SQLs: {self.output_file}")
    
    def _get_namespace(self, cls: OWLClass) -> str:
        """根据类名推断命名空间"""
        if "Substation" in cls.local_name or "VoltageLevel" in cls.local_name or "Breaker" in cls.local_name:
            return "grid"
        elif "Asset" in cls.local_name or "Work" in cls.local_name:
            return "assets"
        elif "Meter" in cls.local_name or "UsagePoint" in cls.local_name:
            return "metering"
        return "grid"  # 默认
    
    def _map_to_iceberg_type(self, xsd_type: Optional[str]) -> str:
        """映射 XSD 类型到 Iceberg 类型"""
        if not xsd_type:
            return "STRING"
        
        type_mapping = {
            "string": "STRING",
            "int": "INT",
            "integer": "BIGINT",
            "float": "FLOAT",
            "double": "DOUBLE",
            "boolean": "BOOLEAN",
            "dateTime": "TIMESTAMP",
            "date": "DATE",
        }
        
        return type_mapping.get(xsd_type.lower(), "STRING")
    
    def _to_snake_case(self, name: str) -> str:
        """转换为 snake_case"""
        s1 = re.sub('(.)([A-Z][a-z]+)', r'\1_\2', name)
        return re.sub('([a-z0-9])([A-Z])', r'\1_\2', s1).lower()
    
    def _get_local_name(self, uri: str) -> str:
        """从 URI 提取本地名称"""
        if '#' in uri:
            return uri.split('#')[-1]
        elif '/' in uri:
            return uri.split('/')[-1]
        return uri
    
    def _extract_property_name(self, local_name: str) -> str:
        """从本地名称提取纯属性名（去掉类名前缀）
        例如：Equipment.aggregate -> aggregate
        """
        if '.' in local_name:
            return local_name.split('.')[-1]
        return local_name


def main():
    owl_file = "TheCimOntology.owl"
    
    if not Path(owl_file).exists():
        print(f"Error: {owl_file} not found")
        return
    
    print("=" * 60)
    print("OWL → Actor + Iceberg Table Generator")
    print("=" * 60)
    
    # 1. 解析 OWL
    print("\n[1/3] Parsing OWL file...")
    parser = OWLParser(owl_file)
    parser.parse()
    
    print(f"  Found {len(parser.classes)} classes")
    print(f"  Found {len(parser.properties)} properties")
    print(f"  Found {len(parser.get_power_system_resources())} PowerSystemResource classes")
    
    # 2. 生成通用的 CIMResourceActor
    print("\n[2/3] Generating CIMResourceActor...")
    actor_gen = ActorCodeGenerator()
    actor_gen.generate()
    
    # 3. 生成 Iceberg 表结构（每个类独立的表）
    print("\n[3/3] Generating Iceberg table SQLs...")
    table_gen = IcebergTableGenerator()
    table_gen.generate_all(parser)
    
    print("\n" + "=" * 60)
    print("Generation completed!")
    print("=" * 60)
    print("\nGenerated files:")
    print("  - actors/cim_resource_actor.go")
    print("  - scripts/generated_iceberg_tables.sql")

if __name__ == "__main__":
    main()

