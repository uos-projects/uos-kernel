"""映射规则配置：定义语义模型到 Iceberg 表的映射规则"""

from __future__ import annotations

from typing import Dict, List, Optional, Callable


# 包 ID 到命名空间的映射
NAMESPACE_MAPPING: Dict[str, str] = {
    "core_topology": "grid",
    "wires_equipment": "grid",
    "asset_mgmt": "assets",
    "metering_market": "metering",
}

# CIM 属性名到 Iceberg 字段类型的映射
ATTRIBUTE_TYPE_MAPPING: Dict[str, str] = {
    # 标识符类型
    "mRID": "STRING",
    "entity_id": "STRING",
    "assetID": "STRING",
    "workID": "STRING",
    "taskID": "STRING",
    
    # 字符串类型
    "name": "STRING",
    "description": "STRING",
    "region": "STRING",
    "serialNumber": "STRING",
    "criticality": "STRING",
    "priority": "STRING",
    "status": "STRING",
    "type": "STRING",
    "crew": "STRING",
    "serviceCategory": "STRING",
    "readingType": "STRING",
    "constructionKind": "STRING",
    
    # 数值类型
    "nominalVoltage": "DOUBLE",
    "nominal_voltage_kv": "DOUBLE",
    "ratedCurrent": "DOUBLE",
    "length": "DOUBLE",
    "r": "DOUBLE",  # 电阻
    "x": "DOUBLE",  # 电抗
    "value": "DOUBLE",
    
    # 布尔类型
    "normalOpen": "BOOLEAN",
    
    # 时间类型
    "installDate": "TIMESTAMP",
    "readingTime": "TIMESTAMP",
    "sequenceNumber": "INT",
}

# 属性名到字段名的映射（CIM 属性名 -> Iceberg 字段名）
ATTRIBUTE_NAME_MAPPING: Dict[str, str] = {
    "mRID": "entity_id",
    "assetID": "asset_id",
    "workID": "work_id",
    "taskID": "task_id",
    "nominalVoltage": "nominal_voltage_kv",
    "ratedCurrent": "rated_current",
    "normalOpen": "normal_open",
    "installDate": "install_date",
    "serialNumber": "serial_number",
    "serviceCategory": "service_category",
    "readingType": "reading_type",
    "readingTime": "reading_time",
    "constructionKind": "construction_kind",
    "sequenceNumber": "sequence_number",
}

# 实体类型到分区字段的映射
PARTITION_STRATEGY: Dict[str, List[str]] = {
    "Substation": ["region"],
    "VoltageLevel": ["substation_id"],  # 注意：需要从关系或属性中推断
    "PowerTransformer": ["substation_id"],
    "Breaker": ["voltage_level_id"],
    "Asset": ["criticality"],
    "Work": ["status"],
    "Meter": ["usage_point_id"],
    "UsagePoint": ["service_category"],
}

# 关系类型到关系表名的映射规则
RELATION_TABLE_MAPPING: Dict[str, str] = {
    # 拓扑关系
    "contains": "substation_voltage_level",  # 默认，可根据源实体类型调整
    "partOf": "breaker_voltage_level",  # 默认
    "locatedIn": "transformer_substation",
    
    # 资产关系
    "realises": "asset_resource",
    "targets": "work_resource",
    
    # 计量关系
    "measuredBy": "meter_usage_point",
}

# 通用字段（所有表都需要）
COMMON_COLUMNS: List[Dict[str, str]] = [
    {"name": "valid_from", "type": "TIMESTAMP"},
    {"name": "valid_to", "type": "TIMESTAMP"},
    {"name": "op_type", "type": "STRING"},
    {"name": "ingestion_ts", "type": "TIMESTAMP"},
]


def get_namespace(package_id: str) -> str:
    """根据包 ID 获取命名空间"""
    return NAMESPACE_MAPPING.get(package_id, "default")


def get_field_type(attribute_name: str) -> str:
    """根据属性名获取字段类型"""
    # 先检查精确匹配
    if attribute_name in ATTRIBUTE_TYPE_MAPPING:
        return ATTRIBUTE_TYPE_MAPPING[attribute_name]
    
    # 根据属性名模式推断
    if attribute_name.endswith("ID") or attribute_name.endswith("Id"):
        return "STRING"
    if attribute_name.endswith("Date") or attribute_name.endswith("Time"):
        return "TIMESTAMP"
    if attribute_name in ["length", "r", "x", "value"]:
        return "DOUBLE"
    if attribute_name in ["normalOpen", "isActive"]:
        return "BOOLEAN"
    
    # 默认返回 STRING
    return "STRING"


def get_field_name(attribute_name: str) -> str:
    """将 CIM 属性名转换为 Iceberg 字段名"""
    # 先检查映射表
    if attribute_name in ATTRIBUTE_NAME_MAPPING:
        return ATTRIBUTE_NAME_MAPPING[attribute_name]
    
    # 驼峰转下划线
    import re
    # 在大写字母前插入下划线（除了第一个字符）
    s1 = re.sub("(.)([A-Z][a-z]+)", r"\1_\2", attribute_name)
    # 在连续大写字母前插入下划线
    s2 = re.sub("([a-z0-9])([A-Z])", r"\1_\2", s1)
    return s2.lower()


def get_partition_fields(entity_name: str) -> List[str]:
    """获取实体的分区字段"""
    return PARTITION_STRATEGY.get(entity_name, [])


def get_relation_table_name(relation_type: str, source_entity: str, target_entity: str) -> str:
    """根据关系类型和实体类型确定关系表名"""
    # 特殊规则：根据实体类型组合确定表名
    if relation_type == "contains":
        if "Substation" in source_entity:
            return "substation_voltage_level"
        # 可以添加更多规则
    elif relation_type == "partOf":
        if "Breaker" in source_entity:
            return "breaker_voltage_level"
    
    # 使用默认映射
    return RELATION_TABLE_MAPPING.get(relation_type, f"{relation_type}_relation")


def camel_to_snake(name: str) -> str:
    """将驼峰命名转换为下划线命名"""
    import re
    s1 = re.sub("(.)([A-Z][a-z]+)", r"\1_\2", name)
    s2 = re.sub("([a-z0-9])([A-Z])", r"\1_\2", s1)
    return s2.lower()


if __name__ == "__main__":
    # 测试映射规则
    print("命名空间映射:")
    for pkg_id, ns in NAMESPACE_MAPPING.items():
        print(f"  {pkg_id} -> {ns}")
    
    print("\n属性类型映射示例:")
    test_attrs = ["mRID", "nominalVoltage", "normalOpen", "installDate"]
    for attr in test_attrs:
        print(f"  {attr} -> {get_field_type(attr)}")
    
    print("\n字段名转换示例:")
    for attr in test_attrs:
        print(f"  {attr} -> {get_field_name(attr)}")
    
    print("\n分区策略示例:")
    test_entities = ["Substation", "VoltageLevel", "Breaker"]
    for entity in test_entities:
        print(f"  {entity} -> {get_partition_fields(entity)}")

