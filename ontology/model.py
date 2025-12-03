"""语义模型数据结构和类型定义"""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import List, Optional, Dict


@dataclass
class Relationship:
    """关系定义"""
    type: str
    target: str
    cardinality: Optional[str] = None


@dataclass
class Entity:
    """实体类型定义"""
    name: str
    package_id: str
    package_name: str
    inherits: Optional[str] = None
    key_attributes: List[str] = field(default_factory=list)
    relationships: List[Relationship] = field(default_factory=list)
    
    def get_all_attributes(self, entity_registry: Dict[str, 'Entity']) -> List[str]:
        """获取所有属性（包括继承的属性）"""
        attrs = set(self.key_attributes)
        
        # 递归获取父类属性
        if self.inherits:
            parent = entity_registry.get(self.inherits)
            if parent:
                attrs.update(parent.get_all_attributes(entity_registry))
        
        return sorted(list(attrs))


@dataclass
class Package:
    """包定义"""
    id: str
    name: str
    summary: str
    entities: List[Entity] = field(default_factory=list)


@dataclass
class RelationTemplate:
    """关系模板定义"""
    id: str
    source_role: str
    target_role: str
    mandatory: bool
    cardinality: str


@dataclass
class SemanticModel:
    """完整的语义模型"""
    version: str
    description: str
    cim_release: str
    last_reviewed: str
    packages: List[Package] = field(default_factory=list)
    relation_templates: List[RelationTemplate] = field(default_factory=list)
    
    def get_entity(self, name: str) -> Optional[Entity]:
        """根据名称获取实体"""
        for package in self.packages:
            for entity in package.entities:
                if entity.name == name:
                    return entity
        return None
    
    def get_all_entities(self) -> List[Entity]:
        """获取所有实体"""
        entities = []
        for package in self.packages:
            entities.extend(package.entities)
        return entities
    
    def get_entity_registry(self) -> Dict[str, Entity]:
        """获取实体注册表（用于查找父类）"""
        return {entity.name: entity for entity in self.get_all_entities()}

