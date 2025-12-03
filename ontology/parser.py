"""解析 cim_scope.yaml 语义模型文件"""

from __future__ import annotations

import os
from pathlib import Path
from typing import Dict, Any

import yaml

from ontology.model import (
    SemanticModel,
    Package,
    Entity,
    Relationship,
    RelationTemplate,
)


def load_cim_scope() -> Path:
    """获取 cim_scope.yaml 文件路径"""
    # 获取项目根目录
    current_dir = Path(__file__).parent
    yaml_path = current_dir / "cim_scope.yaml"
    
    if not yaml_path.exists():
        raise FileNotFoundError(f"语义模型文件不存在: {yaml_path}")
    
    return yaml_path


def parse_relationship(rel_data: Dict[str, Any]) -> Relationship:
    """解析关系定义"""
    return Relationship(
        type=rel_data.get("type", ""),
        target=rel_data.get("target", ""),
        cardinality=rel_data.get("cardinality"),
    )


def parse_entity(entity_data: Dict[str, Any], package_id: str, package_name: str) -> Entity:
    """解析实体定义"""
    relationships = [
        parse_relationship(rel) for rel in entity_data.get("relationships", [])
    ]
    
    return Entity(
        name=entity_data.get("name", ""),
        package_id=package_id,
        package_name=package_name,
        inherits=entity_data.get("inherits"),
        key_attributes=entity_data.get("key_attributes", []),
        relationships=relationships,
    )


def parse_package(package_data: Dict[str, Any]) -> Package:
    """解析包定义"""
    package_id = package_data.get("id", "")
    package_name = package_data.get("name", "")
    
    entities = [
        parse_entity(entity_data, package_id, package_name)
        for entity_data in package_data.get("entities", [])
    ]
    
    return Package(
        id=package_id,
        name=package_name,
        summary=package_data.get("summary", ""),
        entities=entities,
    )


def parse_relation_template(template_data: Dict[str, Any]) -> RelationTemplate:
    """解析关系模板"""
    schema = template_data.get("schema", {})
    return RelationTemplate(
        id=template_data.get("id", ""),
        source_role=schema.get("source_role", ""),
        target_role=schema.get("target_role", ""),
        mandatory=schema.get("mandatory", False),
        cardinality=schema.get("cardinality", ""),
    )


def parse_semantic_model(yaml_path: Path | None = None) -> SemanticModel:
    """解析语义模型 YAML 文件"""
    if yaml_path is None:
        yaml_path = load_cim_scope()
    
    with open(yaml_path, "r", encoding="utf-8") as f:
        data = yaml.safe_load(f)
    
    packages = [parse_package(pkg) for pkg in data.get("packages", [])]
    
    relation_templates = [
        parse_relation_template(template)
        for template in data.get("relation_templates", [])
    ]
    
    return SemanticModel(
        version=data.get("version", ""),
        description=data.get("description", ""),
        cim_release=data.get("cim_release", ""),
        last_reviewed=data.get("last_reviewed", ""),
        packages=packages,
        relation_templates=relation_templates,
    )


def get_semantic_model() -> SemanticModel:
    """获取解析后的语义模型（单例模式，可缓存）"""
    return parse_semantic_model()


if __name__ == "__main__":
    # 测试解析
    model = parse_semantic_model()
    print(f"语义模型版本: {model.version}")
    print(f"包数量: {len(model.packages)}")
    print(f"实体总数: {len(model.get_all_entities())}")
    
    print("\n所有实体:")
    for entity in model.get_all_entities():
        print(f"  - {entity.name} (包: {entity.package_name})")
        print(f"    属性: {entity.key_attributes}")
        print(f"    关系: {[r.type for r in entity.relationships]}")

