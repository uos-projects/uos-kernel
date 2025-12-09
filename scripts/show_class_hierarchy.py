#!/usr/bin/env python3
"""
展示 FullGrid 数据集对应的 Python 类的完整层次结构
"""
import sys
from pathlib import Path
from collections import defaultdict

# 添加 cimpyorm 到 sys.path
project_root = Path(__file__).parent.parent
sys.path.insert(0, str(project_root / 'cimpyorm'))

from cimpyorm import load
from cimpyorm.Model.Elements.Class import CIMClass

def build_hierarchy_tree(session):
    """构建类层次结构树"""
    # 获取所有类
    all_classes = session.query(CIMClass).all()
    
    # 创建类字典，key 是 (namespace, name)
    classes_dict = {}
    for cls in all_classes:
        key = (cls.namespace.short if cls.namespace else 'cim', cls.name)
        classes_dict[key] = cls
    
    # 构建父子关系树
    tree = defaultdict(list)
    root_classes = []
    
    for cls in all_classes:
        if cls.parent is None:
            root_classes.append(cls)
        else:
            parent_key = (cls.parent.namespace.short if cls.parent.namespace else 'cim', cls.parent.name)
            tree[parent_key].append(cls)
    
    return tree, root_classes, classes_dict

def print_tree(cls, tree, classes_dict, prefix="", is_last=True, max_depth=None, current_depth=0):
    """递归打印树结构"""
    if max_depth and current_depth >= max_depth:
        return
    
    namespace = cls.namespace.short if cls.namespace else 'cim'
    class_name = f"{namespace}_{cls.name}"
    
    # 打印当前节点
    connector = "└── " if is_last else "├── "
    print(f"{prefix}{connector}{class_name}")
    
    # 更新前缀
    new_prefix = prefix + ("    " if is_last else "│   ")
    
    # 获取子类
    key = (namespace, cls.name)
    children = tree.get(key, [])
    
    # 按名称排序子类
    children.sort(key=lambda x: x.name)
    
    # 递归打印子类
    for i, child in enumerate(children):
        is_last_child = (i == len(children) - 1)
        print_tree(child, tree, classes_dict, new_prefix, is_last_child, max_depth, current_depth + 1)

def main():
    # 数据库路径
    db_path = project_root / 'cimpyorm' / 'cimpyorm' / 'res' / 'datasets' / 'FullGrid' / 'out.db'
    
    print("=" * 100)
    print("FullGrid 数据集对应的 Python 类层次结构")
    print("=" * 100)
    
    # 加载数据库
    session, model = load(str(db_path))
    
    # 获取 Schema
    schema = session.schema
    
    # 构建层次结构树
    tree, root_classes, classes_dict = build_hierarchy_tree(session)
    
    # 按名称排序根类
    root_classes.sort(key=lambda x: x.name)
    
    print(f"\n总共 {len(classes_dict)} 个类")
    print(f"根类数量: {len(root_classes)}\n")
    
    # 打印完整的类层次结构
    print("=" * 100)
    print("完整类层次结构:")
    print("=" * 100)
    
    for i, root_cls in enumerate(root_classes):
        is_last_root = (i == len(root_classes) - 1)
        print_tree(root_cls, tree, classes_dict, "", is_last_root)
    
    # 统计信息
    print("\n" + "=" * 100)
    print("统计信息:")
    print("=" * 100)
    
    # 统计每个根类的子类数量
    def count_children(cls_key, tree):
        count = len(tree.get(cls_key, []))
        for child_key in tree.get(cls_key, []):
            child_namespace = child_key.namespace.short if child_key.namespace else 'cim'
            child_key_tuple = (child_namespace, child_key.name)
            count += count_children(child_key_tuple, tree)
        return count
    
    print("\n根类及其子类数量:")
    for root_cls in root_classes:
        namespace = root_cls.namespace.short if root_cls.namespace else 'cim'
        key = (namespace, root_cls.name)
        child_count = count_children(key, tree)
        print(f"  {namespace}_{root_cls.name}: {child_count} 个子类")
    
    # 查找一些重要的类
    print("\n" + "=" * 100)
    print("重要类的层次结构:")
    print("=" * 100)
    
    important_classes = [
        ('cim', 'IdentifiedObject'),
        ('cim', 'PowerSystemResource'),
        ('cim', 'Measurement'),
        ('cim', 'Control'),
        ('cim', 'Equipment'),
        ('cim', 'ConductingEquipment'),
    ]
    
    for ns, name in important_classes:
        key = (ns, name)
        if key in classes_dict:
            cls = classes_dict[key]
            print(f"\n{ns}_{name} 的子类层次:")
            print_tree(cls, tree, classes_dict, "", True, max_depth=3)
    
    print("\n" + "=" * 100)
    print("完成!")
    print("=" * 100)

if __name__ == '__main__':
    main()

