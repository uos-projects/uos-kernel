#!/usr/bin/env python3
"""测试语义模型 API 功能"""

import sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).parent.parent))

from ontology.parser import get_semantic_model
from mapping.mapper import EntityMapper
from table_gen.generator import SQLGenerator


def test_semantic_model():
    """测试语义模型解析"""
    print("=== 测试语义模型解析 ===\n")
    
    model = get_semantic_model()
    print(f"版本: {model.version}")
    print(f"包数量: {len(model.packages)}")
    print(f"实体总数: {len(model.get_all_entities())}\n")
    
    # 测试获取实体
    entity = model.get_entity("Substation")
    if entity:
        print(f"实体: {entity.name}")
        print(f"  属性: {entity.key_attributes}")
        print(f"  关系: {[r.type for r in entity.relationships]}\n")


def test_mapping():
    """测试映射功能"""
    print("=== 测试实体到表映射 ===\n")
    
    model = get_semantic_model()
    mapper = EntityMapper(model)
    
    test_entities = ["Substation", "VoltageLevel", "Breaker"]
    
    for entity_name in test_entities:
        schema = mapper.map_entity_by_name(entity_name)
        if schema:
            print(f"实体: {entity_name}")
            print(f"  表名: {schema.get_full_table_name()}")
            print(f"  列数: {len(schema.columns)}")
            print(f"  分区: {schema.partition_fields}\n")


def test_sql_generation():
    """测试 SQL 生成"""
    print("=== 测试 SQL 生成 ===\n")
    
    model = get_semantic_model()
    mapper = EntityMapper(model)
    schema = mapper.map_entity_by_name("Substation")
    
    if schema:
        sqls = SQLGenerator.generate_table_sqls(schema)
        print("生成的 SQL:")
        for sql in sqls:
            print(sql)
            print()


if __name__ == "__main__":
    test_semantic_model()
    test_mapping()
    test_sql_generation()

