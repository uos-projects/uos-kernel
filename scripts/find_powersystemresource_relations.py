#!/usr/bin/env python3
"""
查找 CIM16 schema 中 PowerSystemResource 的所有关联关系
"""
import sys
from pathlib import Path

# 添加 cimpyorm 到 sys.path
project_root = Path(__file__).parent.parent
sys.path.insert(0, str(project_root / 'cimpyorm'))

from cimpyorm.Model.Schema import Schema
from cimpyorm.Model.Elements.Class import CIMClass
from cimpyorm.Model.Elements.Property import CIMProp
from cimpyorm.backends import InMemory

def find_powersystemresource_relations():
    """查找 PowerSystemResource 的所有关联关系"""
    # CIM16 schema 路径
    schema_path = project_root / 'cimpyorm' / 'cimpyorm' / 'res' / 'schemata' / 'CIM16'
    
    print("=" * 100)
    print("CIM16 Schema 中 PowerSystemResource 的所有关联关系")
    print("=" * 100)
    print(f"Schema 路径: {schema_path}\n")
    
    # 创建内存数据库并加载 schema
    backend = InMemory()
    backend.reset()
    session = backend.ORM
    
    schema = Schema(dataset=session, version="16", rdfs_path=str(schema_path))
    
    # 查找 PowerSystemResource 类
    psr_class = session.query(CIMClass).filter(
        CIMClass.name == 'PowerSystemResource'
    ).first()
    
    if not psr_class:
        print("未找到 PowerSystemResource 类")
        return
    
    print(f"【PowerSystemResource 基本信息】")
    print(f"  名称: {psr_class.name}")
    print(f"  命名空间: {psr_class.namespace_name}")
    if psr_class.parent_name:
        print(f"  父类: {psr_class.parent_name}")
    print()
    
    # 查找 PowerSystemResource 的所有属性
    print("=" * 100)
    print("【PowerSystemResource 的所有属性】")
    print("=" * 100)
    
    # 直接定义的属性
    direct_props = session.query(CIMProp).filter(
        CIMProp.cls_name == 'PowerSystemResource',
        CIMProp.cls_namespace == psr_class.namespace_name
    ).all()
    
    # 递归查找继承的属性
    def get_inherited_props(class_name, namespace_name, visited=None):
        """递归获取继承的属性"""
        if visited is None:
            visited = set()
        
        key = (class_name, namespace_name)
        if key in visited:
            return []
        
        visited.add(key)
        
        props = []
        # 当前类的属性
        current_props = session.query(CIMProp).filter(
            CIMProp.cls_name == class_name,
            CIMProp.cls_namespace == namespace_name
        ).all()
        props.extend(current_props)
        
        # 父类的属性
        parent = session.query(CIMClass).filter(
            CIMClass.name == class_name,
            CIMClass.namespace_name == namespace_name
        ).first()
        
        if parent and parent.parent_name:
            parent_props = get_inherited_props(
                parent.parent_name,
                parent.parent_namespace or namespace_name,
                visited
            )
            props.extend(parent_props)
        
        return props
    
    all_props = get_inherited_props('PowerSystemResource', psr_class.namespace_name)
    
    # 按类型分类
    object_properties = []
    datatype_properties = []
    enum_properties = []
    
    for prop in all_props:
        if prop.inverse_class_name:
            object_properties.append(prop)
        elif prop.datatype_name:
            if prop.type_ == 'CIMProp_Enumeration':
                enum_properties.append(prop)
            else:
                datatype_properties.append(prop)
    
    print(f"\n找到 {len(all_props)} 个属性（包括继承的）")
    print(f"  - ObjectProperty（对象属性）: {len(object_properties)} 个")
    print(f"  - DatatypeProperty（数据类型属性）: {len(datatype_properties)} 个")
    print(f"  - EnumerationProperty（枚举属性）: {len(enum_properties)} 个")
    
    # 显示 ObjectProperty（关联关系）
    print("\n" + "=" * 100)
    print("【ObjectProperty - 关联的其他类型】")
    print("=" * 100)
    
    if object_properties:
        for prop in sorted(object_properties, key=lambda p: p.name):
            print(f"\n- {prop.name}")
            print(f"  类型: {prop.type_}")
            print(f"  关联的类: {prop.inverse_class_name}")
            if prop.inverse_property_name:
                print(f"  反向属性: {prop.inverse_property_name}")
            if prop.multiplicity:
                print(f"  基数: {prop.multiplicity}")
            if prop.optional is not None:
                print(f"  可选: {prop.optional}")
            if prop.used is not None:
                print(f"  已使用: {prop.used}")
    else:
        print("无 ObjectProperty")
    
    # 显示 DatatypeProperty
    print("\n" + "=" * 100)
    print("【DatatypeProperty - 数据类型属性】")
    print("=" * 100)
    
    if datatype_properties:
        for prop in sorted(datatype_properties, key=lambda p: p.name):
            print(f"\n- {prop.name}")
            print(f"  类型: {prop.type_}")
            print(f"  数据类型: {prop.datatype_name}")
            if prop.multiplicity:
                print(f"  基数: {prop.multiplicity}")
            if prop.optional is not None:
                print(f"  可选: {prop.optional}")
    else:
        print("无 DatatypeProperty")
    
    # 显示 EnumerationProperty
    print("\n" + "=" * 100)
    print("【EnumerationProperty - 枚举属性】")
    print("=" * 100)
    
    if enum_properties:
        for prop in sorted(enum_properties, key=lambda p: p.name):
            print(f"\n- {prop.name}")
            print(f"  类型: {prop.type_}")
            if prop.multiplicity:
                print(f"  基数: {prop.multiplicity}")
            if prop.optional is not None:
                print(f"  可选: {prop.optional}")
    else:
        print("无 EnumerationProperty")
    
    # 统计关联的类
    print("\n" + "=" * 100)
    print("【关联的类统计】")
    print("=" * 100)
    
    associated_classes = {}
    for prop in object_properties:
        class_name = prop.inverse_class_name
        if class_name not in associated_classes:
            associated_classes[class_name] = []
        associated_classes[class_name].append({
            'property': prop.name,
            'inverse': prop.inverse_property_name,
            'multiplicity': prop.multiplicity,
        })
    
    print(f"\nPowerSystemResource 关联了 {len(associated_classes)} 个不同的类：\n")
    for class_name, props_list in sorted(associated_classes.items()):
        print(f"  {class_name}:")
        for prop_info in props_list:
            print(f"    - {prop_info['property']} (反向: {prop_info['inverse']}, 基数: {prop_info['multiplicity']})")
    
    print("\n" + "=" * 100)

if __name__ == '__main__':
    find_powersystemresource_relations()

