#!/usr/bin/env python3
"""
查找 CIM16 schema 中所有 Control 类及其关联类型
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

def find_control_classes():
    """查找所有 Control 相关的类"""
    # CIM16 schema 路径
    schema_path = project_root / 'cimpyorm' / 'cimpyorm' / 'res' / 'schemata' / 'CIM16'
    
    print("=" * 100)
    print("CIM16 Schema 中所有 Control 类及其关联类型")
    print("=" * 100)
    print(f"Schema 路径: {schema_path}\n")
    
    # 创建内存数据库并加载 schema
    backend = InMemory()
    backend.reset()
    session = backend.ORM
    
    schema = Schema(dataset=session, version="16", rdfs_path=str(schema_path))
    
    # 查找所有 Control 相关的类
    # 查找名称包含 "Control" 的类
    control_classes = session.query(CIMClass).filter(
        CIMClass.name.like('%Control%')
    ).all()
    
    # 也查找继承自 Control 的类
    all_control_classes = set(control_classes)
    
    # 递归查找所有子类
    def find_subclasses(parent_class):
        """递归查找所有子类"""
        subclasses = session.query(CIMClass).filter(
            CIMClass.parent_name == parent_class.name,
            CIMClass.parent_namespace == parent_class.namespace_name
        ).all()
        for subclass in subclasses:
            if subclass not in all_control_classes:
                all_control_classes.add(subclass)
                find_subclasses(subclass)
    
    # 查找 Control 基类
    control_base = session.query(CIMClass).filter(
        CIMClass.name == 'Control'
    ).first()
    
    if control_base:
        all_control_classes.add(control_base)
        find_subclasses(control_base)
    
    # 按名称排序
    control_classes_sorted = sorted(all_control_classes, key=lambda c: c.name)
    
    print(f"找到 {len(control_classes_sorted)} 个 Control 相关的类\n")
    
    # 遍历每个 Control 类
    for idx, cim_class in enumerate(control_classes_sorted, 1):
        print("=" * 100)
        print(f"[{idx}] {cim_class.name}")
        print("=" * 100)
        
        # 类基本信息
        print(f"\n【类基本信息】")
        print(f"  名称: {cim_class.name}")
        print(f"  命名空间: {cim_class.namespace_name}")
        if cim_class.parent_name:
            print(f"  父类: {cim_class.parent_name}")
        if cim_class.package_name:
            print(f"  包: {cim_class.package_name}")
        
        # 查找类的所有属性
        print(f"\n【属性 (Properties)】")
        props = session.query(CIMProp).filter(
            CIMProp.cls_name == cim_class.name,
            CIMProp.cls_namespace == cim_class.namespace_name
        ).all()
        
        if not props:
            print("  无属性")
        else:
            for prop in props:
                print(f"\n  - {prop.name}")
                print(f"    类型: {prop.type_}")
                print(f"    命名空间: {prop.namespace_name}")
                
                # 属性类型信息
                if prop.datatype_name:
                    print(f"    数据类型: {prop.datatype_name}")
                
                # 关联的类（如果是 ObjectProperty）
                if prop.inverse_class_name:
                    print(f"    关联的类: {prop.inverse_class_name}")
                    if prop.inverse_property_name:
                        print(f"    反向属性: {prop.inverse_property_name}")
                
                # 基数信息
                if prop.multiplicity:
                    print(f"    基数: {prop.multiplicity}")
                if prop.optional is not None:
                    print(f"    可选: {prop.optional}")
                if prop.used is not None:
                    print(f"    已使用: {prop.used}")
        
        # 查找所有关联的类型（通过属性）
        print(f"\n【关联的类型】")
        associated_types = set()
        
        for prop in props:
            # 通过 inverse_class_name 找到关联的类
            if prop.inverse_class_name:
                associated_types.add((prop.inverse_class_name, prop.inverse_property_name))
            
            # 通过 datatype 找到数据类型
            if prop.datatype_name:
                associated_types.add(("Datatype", prop.datatype_name))
        
        if not associated_types:
            print("  无关联类型")
        else:
            for assoc_type, prop_name in sorted(associated_types):
                if assoc_type == "Datatype":
                    print(f"  - {assoc_type}: {prop_name}")
                else:
                    print(f"  - 类: {assoc_type} (通过属性: {prop_name})")
        
        # 查找继承关系
        print(f"\n【继承关系】")
        if cim_class.parent_name:
            print(f"  父类: {cim_class.parent_name}")
        
        # 查找子类
        children = session.query(CIMClass).filter(
            CIMClass.parent_name == cim_class.name,
            CIMClass.parent_namespace == cim_class.namespace_name
        ).all()
        
        if children:
            print(f"  子类 ({len(children)} 个):")
            for child in children:
                print(f"    - {child.name}")
        else:
            print("  无子类")
        
        print()
    
    # 统计信息
    print("=" * 100)
    print("统计信息")
    print("=" * 100)
    
    # 统计属性类型
    prop_types = {}
    for cim_class in control_classes_sorted:
        props = session.query(CIMProp).filter(
            CIMProp.cls_name == cim_class.name,
            CIMProp.cls_namespace == cim_class.namespace_name
        ).all()
        for prop in props:
            prop_type = prop.type_ or "Unknown"
            prop_types[prop_type] = prop_types.get(prop_type, 0) + 1
    
    print("\n属性类型分布:")
    for prop_type, count in sorted(prop_types.items()):
        print(f"  {prop_type}: {count} 个")
    
    # 统计关联的类
    associated_classes = {}
    for cim_class in control_classes_sorted:
        props = session.query(CIMProp).filter(
            CIMProp.cls_name == cim_class.name,
            CIMProp.cls_namespace == cim_class.namespace_name
        ).all()
        for prop in props:
            if prop.inverse_class_name:
                class_name = prop.inverse_class_name
                if class_name not in associated_classes:
                    associated_classes[class_name] = []
                associated_classes[class_name].append(f"{cim_class.name}.{prop.name}")
    
    print("\n关联的类:")
    for class_name, props_list in sorted(associated_classes.items()):
        print(f"\n  {class_name}:")
        for prop_ref in props_list:
            print(f"    - {prop_ref}")
    
    print("\n" + "=" * 100)

if __name__ == '__main__':
    find_control_classes()

