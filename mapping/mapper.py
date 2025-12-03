"""自动映射引擎：从语义模型生成 Iceberg 表结构"""

from __future__ import annotations

from dataclasses import dataclass
from typing import List, Dict, Optional

from ontology.model import Entity, SemanticModel
from ontology.parser import get_semantic_model
from mapping.rules import (
    get_namespace,
    get_field_type,
    get_field_name,
    get_partition_fields,
    COMMON_COLUMNS,
    camel_to_snake,
)


@dataclass
class Column:
    """表列定义"""
    name: str
    type: str
    nullable: bool = True
    comment: Optional[str] = None


@dataclass
class TableSchema:
    """Iceberg 表结构定义"""
    namespace: str
    table_name: str
    columns: List[Column]
    partition_fields: List[str]
    comment: Optional[str] = None
    
    def get_full_table_name(self) -> str:
        """获取完整表名"""
        return f"ontology.{self.namespace}.{self.table_name}"


class EntityMapper:
    """实体到表的映射器"""
    
    def __init__(self, semantic_model: Optional[SemanticModel] = None):
        """初始化映射器"""
        if semantic_model is None:
            semantic_model = get_semantic_model()
        self.semantic_model = semantic_model
        self.entity_registry = semantic_model.get_entity_registry()
    
    def map_entity_to_table(self, entity: Entity) -> TableSchema:
        """将实体映射到表结构"""
        # 1. 确定命名空间
        namespace = get_namespace(entity.package_id)
        
        # 2. 生成表名（驼峰转下划线）
        table_name = camel_to_snake(entity.name)
        
        # 3. 收集所有属性（包括继承的）
        all_attributes = entity.get_all_attributes(self.entity_registry)
        
        # 4. 映射属性到列
        columns = []
        for attr in all_attributes:
            column = Column(
                name=get_field_name(attr),
                type=get_field_type(attr),
                comment=f"CIM attribute: {attr}",
            )
            columns.append(column)
        
        # 4.5. 根据关系添加外键字段（用于分区）
        # 例如：VoltageLevel 有 partOf -> Substation 关系，需要 substation_id
        for rel in entity.relationships:
            if rel.type == "partOf":
                # 生成外键字段名：Substation -> substation_id
                fk_field = f"{camel_to_snake(rel.target)}_id"
                # 检查是否已存在
                if not any(c.name == fk_field for c in columns):
                    columns.append(Column(
                        name=fk_field,
                        type="STRING",
                        comment=f"Foreign key from relationship: {rel.type} -> {rel.target}",
                    ))
        
        # 5. 添加通用字段
        for common_col in COMMON_COLUMNS:
            # 检查是否已存在（避免重复）
            if not any(c.name == common_col["name"] for c in columns):
                columns.append(Column(
                    name=common_col["name"],
                    type=common_col["type"],
                    comment="Common time travel field",
                ))
        
        # 6. 自动确定分区字段
        partition_fields = self._auto_determine_partition_fields(entity, columns)
        
        return TableSchema(
            namespace=namespace,
            table_name=table_name,
            columns=columns,
            partition_fields=partition_fields,
            comment=f"Table for CIM entity: {entity.name}",
        )
    
    def map_all_entities(self) -> Dict[str, TableSchema]:
        """映射所有实体到表结构"""
        schemas = {}
        for entity in self.semantic_model.get_all_entities():
            schema = self.map_entity_to_table(entity)
            schemas[entity.name] = schema
        return schemas
    
    def map_entity_by_name(self, entity_name: str) -> Optional[TableSchema]:
        """根据实体名称映射到表结构"""
        entity = self.semantic_model.get_entity(entity_name)
        if entity is None:
            return None
        return self.map_entity_to_table(entity)
    
    def _auto_determine_partition_fields(self, entity: Entity, columns: List[Column]) -> List[str]:
        """自动确定分区字段的智能策略"""
        partition_fields = []
        
        # 策略 1: 使用配置的显式分区策略
        configured_fields = get_partition_fields(entity.name)
        if configured_fields:
            # 确保配置的字段在列中存在
            partition_fields = [f for f in configured_fields if any(c.name == f for c in columns)]
            if partition_fields:
                return partition_fields
        
        # 策略 2: 根据关系自动推断（partOf 关系通常表示层级，适合分区）
        for rel in entity.relationships:
            if rel.type == "partOf":
                # 如果实体属于另一个实体，使用父实体的 ID 作为分区字段
                fk_field = f"{camel_to_snake(rel.target)}_id"
                if any(c.name == fk_field for c in columns):
                    partition_fields = [fk_field]
                    return partition_fields
        
        # 策略 3: 根据属性模式匹配（常见分区字段模式）
        # 优先级：region > *_id (外键) > status/category 等枚举字段
        priority_patterns = [
            "region",                    # 地理分区
            "substation_id",             # 变电站分区
            "voltage_level_id",          # 电压等级分区
            "status",                    # 状态分区
            "criticality",               # 重要性分区
            "service_category",          # 服务类别分区
        ]
        
        for pattern in priority_patterns:
            for col in columns:
                if col.name == pattern:
                    partition_fields = [col.name]
                    return partition_fields
        
        # 策略 4: 查找外键字段（以 _id 结尾的字段）
        fk_fields = [col.name for col in columns if col.name.endswith("_id") and col.name != "entity_id"]
        if fk_fields:
            # 选择第一个外键字段
            partition_fields = [fk_fields[0]]
            return partition_fields
        
        # 策略 5: 查找枚举类型字段（STRING 类型，可能的值有限）
        # 这里简化处理，选择第一个非时间、非ID的字符串字段
        candidate_fields = [
            col.name for col in columns
            if col.type == "STRING"
            and col.name not in ["entity_id", "name", "description"]
            and not col.name.endswith("_id")
            and col.name not in ["valid_from", "valid_to", "op_type", "ingestion_ts"]
        ]
        if candidate_fields:
            partition_fields = [candidate_fields[0]]
            return partition_fields
        
        # 策略 6: 如果没有合适的字段，返回空列表（不分区）
        return []


if __name__ == "__main__":
    # 测试映射
    mapper = EntityMapper()
    
    print("=== 实体到表映射测试 ===\n")
    
    # 测试单个实体
    test_entities = ["Substation", "VoltageLevel", "PowerTransformer", "Breaker"]
    
    for entity_name in test_entities:
        schema = mapper.map_entity_by_name(entity_name)
        if schema:
            print(f"实体: {entity_name}")
            print(f"  表名: {schema.get_full_table_name()}")
            print(f"  列数: {len(schema.columns)}")
            print(f"  分区字段: {schema.partition_fields}")
            print(f"  列列表:")
            for col in schema.columns[:5]:  # 只显示前5列
                print(f"    - {col.name}: {col.type}")
            if len(schema.columns) > 5:
                print(f"    ... 还有 {len(schema.columns) - 5} 列")
            print()

