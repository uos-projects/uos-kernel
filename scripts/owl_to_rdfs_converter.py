#!/usr/bin/env python3
"""
将 TheCimOntology.owl (OWL 格式) 转换为 RDFS 格式（与 schemata/CIM16 格式一致）

转换规则：
1. owl:Class → rdf:Description + rdf:type rdfs:Class
2. owl:ObjectProperty → rdf:Description + rdf:type rdf:Property
3. owl:DatatypeProperty → rdf:Description + rdf:type rdf:Property
4. 命名空间转换：http://www.iec.ch/TC57/CIM# → http://iec.ch/TC57/2013/CIM-schema-cim16#
5. owl:inverseOf → cims:inverseRoleName（保留反向关系信息）
"""

import xml.etree.ElementTree as ET
from pathlib import Path
from collections import defaultdict
import sys

# 命名空间定义
NSMAP = {
    'owl': 'http://www.w3.org/2002/07/owl#',
    'rdf': 'http://www.w3.org/1999/02/22-rdf-syntax-ns#',
    'rdfs': 'http://www.w3.org/2000/01/rdf-schema#',
    'xsd': 'http://www.w3.org/2001/XMLSchema#',
    'cims': 'http://iec.ch/TC57/1999/rdf-schema-extensions-19990926#',
    'cim': 'http://iec.ch/TC57/2013/CIM-schema-cim16#',
}

# 源和目标命名空间
SOURCE_NS = 'http://www.iec.ch/TC57/CIM#'
TARGET_NS = 'http://iec.ch/TC57/2013/CIM-schema-cim16#'
TARGET_BASE = 'http://iec.ch/TC57/2013/CIM-schema-cim16'


def convert_uri(uri):
    """转换 URI 从源命名空间到目标命名空间"""
    if uri.startswith(SOURCE_NS):
        return uri.replace(SOURCE_NS, TARGET_NS)
    return uri


def get_local_name(uri):
    """从完整 URI 中提取本地名称"""
    if '#' in uri:
        return uri.split('#')[-1]
    elif '/' in uri:
        return uri.split('/')[-1]
    return uri


def convert_class(owl_class_elem):
    """转换 owl:Class 为 rdf:Description"""
    about = owl_class_elem.get(f'{{{NSMAP["rdf"]}}}about')
    if not about:
        return None
    
    # 转换 URI
    local_name = get_local_name(about)
    new_about = f"#{local_name}"
    
    # 创建 rdf:Description
    desc = ET.Element(f'{{{NSMAP["rdf"]}}}Description')
    desc.set(f'{{{NSMAP["rdf"]}}}about', new_about)
    
    # 添加 rdf:type
    type_elem = ET.SubElement(desc, f'{{{NSMAP["rdf"]}}}type')
    type_elem.set(f'{{{NSMAP["rdf"]}}}resource', f'{NSMAP["rdfs"]}Class')
    
    # 转换子元素
    for child in owl_class_elem:
        tag = child.tag
        if tag == f'{{{NSMAP["rdfs"]}}}subClassOf':
            # 保留 subClassOf
            sub_elem = ET.SubElement(desc, tag)
            resource = child.get(f'{{{NSMAP["rdf"]}}}resource')
            if resource:
                local = get_local_name(resource)
                sub_elem.set(f'{{{NSMAP["rdf"]}}}resource', f'#{local}')
        elif tag == f'{{{NSMAP["rdfs"]}}}comment':
            # 保留 comment
            comment_elem = ET.SubElement(desc, tag)
            comment_elem.set(f'{{{NSMAP["rdf"]}}}parseType', 'Literal')
            if child.text:
                comment_elem.text = child.text
        elif tag == f'{{{NSMAP["rdfs"]}}}label':
            # 保留 label
            label_elem = ET.SubElement(desc, tag)
            xml_lang = child.get('{http://www.w3.org/XML/1998/namespace}lang')
            if xml_lang:
                label_elem.set('{http://www.w3.org/XML/1998/namespace}lang', xml_lang)
            if child.text:
                label_elem.text = child.text
    
    return desc


def convert_object_property(prop_elem):
    """转换 owl:ObjectProperty 为 rdf:Description"""
    about = prop_elem.get(f'{{{NSMAP["rdf"]}}}about')
    if not about:
        return None
    
    local_name = get_local_name(about)
    new_about = f"#{local_name}"
    
    desc = ET.Element(f'{{{NSMAP["rdf"]}}}Description')
    desc.set(f'{{{NSMAP["rdf"]}}}about', new_about)
    
    # 添加 rdf:type
    type_elem = ET.SubElement(desc, f'{{{NSMAP["rdf"]}}}type')
    type_elem.set(f'{{{NSMAP["rdf"]}}}resource', f'{NSMAP["rdf"]}Property')
    
    # 转换子元素
    for child in prop_elem:
        tag = child.tag
        if tag == f'{{{NSMAP["owl"]}}}inverseOf':
            # 转换为 cims:inverseRoleName
            inv_elem = ET.SubElement(desc, f'{{{NSMAP["cims"]}}}inverseRoleName')
            resource = child.get(f'{{{NSMAP["rdf"]}}}resource')
            if resource:
                local = get_local_name(resource)
                inv_elem.set(f'{{{NSMAP["rdf"]}}}resource', f'#{local}')
        elif tag == f'{{{NSMAP["rdfs"]}}}domain':
            domain_elem = ET.SubElement(desc, tag)
            resource = child.get(f'{{{NSMAP["rdf"]}}}resource')
            if resource:
                local = get_local_name(resource)
                domain_elem.set(f'{{{NSMAP["rdf"]}}}resource', f'#{local}')
        elif tag == f'{{{NSMAP["rdfs"]}}}range':
            range_elem = ET.SubElement(desc, tag)
            resource = child.get(f'{{{NSMAP["rdf"]}}}resource')
            if resource:
                local = get_local_name(resource)
                range_elem.set(f'{{{NSMAP["rdf"]}}}resource', f'#{local}')
        elif tag == f'{{{NSMAP["rdfs"]}}}comment':
            comment_elem = ET.SubElement(desc, tag)
            comment_elem.set(f'{{{NSMAP["rdf"]}}}parseType', 'Literal')
            if child.text:
                comment_elem.text = child.text
        elif tag == f'{{{NSMAP["rdfs"]}}}label':
            label_elem = ET.SubElement(desc, tag)
            xml_lang = child.get('{http://www.w3.org/XML/1998/namespace}lang')
            if xml_lang:
                label_elem.set('{http://www.w3.org/XML/1998/namespace}lang', xml_lang)
            if child.text:
                label_elem.text = child.text
    
    return desc


def convert_datatype_property(prop_elem):
    """转换 owl:DatatypeProperty 为 rdf:Description"""
    about = prop_elem.get(f'{{{NSMAP["rdf"]}}}about')
    if not about:
        return None
    
    local_name = get_local_name(about)
    new_about = f"#{local_name}"
    
    desc = ET.Element(f'{{{NSMAP["rdf"]}}}Description')
    desc.set(f'{{{NSMAP["rdf"]}}}about', new_about)
    
    # 添加 rdf:type
    type_elem = ET.SubElement(desc, f'{{{NSMAP["rdf"]}}}type')
    type_elem.set(f'{{{NSMAP["rdf"]}}}resource', f'{NSMAP["rdf"]}Property')
    
    # 添加 cims:stereotype
    stereotype_elem = ET.SubElement(desc, f'{{{NSMAP["cims"]}}}stereotype')
    stereotype_elem.set(f'{{{NSMAP["rdf"]}}}resource', 'http://iec.ch/TC57/NonStandard/UML#attribute')
    
    # 转换子元素
    for child in prop_elem:
        tag = child.tag
        if tag == f'{{{NSMAP["rdfs"]}}}domain':
            domain_elem = ET.SubElement(desc, tag)
            resource = child.get(f'{{{NSMAP["rdf"]}}}resource')
            if resource:
                local = get_local_name(resource)
                domain_elem.set(f'{{{NSMAP["rdf"]}}}resource', f'#{local}')
        elif tag == f'{{{NSMAP["rdfs"]}}}range':
            range_elem = ET.SubElement(desc, tag)
            resource = child.get(f'{{{NSMAP["rdf"]}}}resource')
            if resource:
                local = get_local_name(resource)
                range_elem.set(f'{{{NSMAP["rdf"]}}}resource', f'#{local}')
        elif tag == f'{{{NSMAP["rdfs"]}}}comment':
            comment_elem = ET.SubElement(desc, tag)
            comment_elem.set(f'{{{NSMAP["rdf"]}}}parseType', 'Literal')
            if child.text:
                comment_elem.text = child.text
        elif tag == f'{{{NSMAP["rdfs"]}}}label':
            label_elem = ET.SubElement(desc, tag)
            xml_lang = child.get('{http://www.w3.org/XML/1998/namespace}lang')
            if xml_lang:
                label_elem.set('{http://www.w3.org/XML/1998/namespace}lang', xml_lang)
            if child.text:
                label_elem.text = child.text
    
    return desc


def convert_owl_to_rdfs(owl_file, output_file):
    """主转换函数"""
    print(f"正在解析 {owl_file}...")
    tree = ET.parse(owl_file)
    root = tree.getroot()
    
    # 创建新的 RDF 根元素
    rdf_root = ET.Element(f'{{{NSMAP["rdf"]}}}RDF')
    # 设置命名空间（避免重复）
    rdf_root.set('xmlns:cims', NSMAP['cims'])
    rdf_root.set('xmlns:rdf', NSMAP['rdf'])
    rdf_root.set('xmlns:xsd', NSMAP['xsd'])
    rdf_root.set('xmlns:cim', NSMAP['cim'])
    rdf_root.set('xmlns:rdfs', NSMAP['rdfs'])
    rdf_root.set('{http://www.w3.org/XML/1998/namespace}base', TARGET_BASE)
    
    # 转换类
    print("转换类定义...")
    classes = root.findall(f'.//{{{NSMAP["owl"]}}}Class', NSMAP)
    print(f"  找到 {len(classes)} 个类")
    for owl_class in classes:
        rdfs_class = convert_class(owl_class)
        if rdfs_class:
            rdf_root.append(rdfs_class)
    
    # 转换 ObjectProperty
    print("转换 ObjectProperty...")
    obj_props = root.findall(f'.//{{{NSMAP["owl"]}}}ObjectProperty', NSMAP)
    print(f"  找到 {len(obj_props)} 个 ObjectProperty")
    for prop in obj_props:
        rdfs_prop = convert_object_property(prop)
        if rdfs_prop:
            rdf_root.append(rdfs_prop)
    
    # 转换 DatatypeProperty
    print("转换 DatatypeProperty...")
    datatype_props = root.findall(f'.//{{{NSMAP["owl"]}}}DatatypeProperty', NSMAP)
    print(f"  找到 {len(datatype_props)} 个 DatatypeProperty")
    for prop in datatype_props:
        rdfs_prop = convert_datatype_property(prop)
        if rdfs_prop:
            rdf_root.append(rdfs_prop)
    
    # 创建新的树并写入
    new_tree = ET.ElementTree(rdf_root)
    ET.indent(new_tree, space='\t')
    
    print(f"正在写入 {output_file}...")
    # 直接写入，然后手动修复命名空间
    new_tree.write(output_file, encoding='utf-8', xml_declaration=True)
    
    # 读取并修复重复的命名空间
    import re
    with open(output_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 修复 XML 声明
    if not content.startswith('<?xml version="1.0" encoding="UTF-8"?>'):
        content = '<?xml version="1.0" encoding="UTF-8"?>\n' + content.split('?>', 1)[-1] if '?>' in content else content
    
    # 找到 RDF 根元素行并清理重复的命名空间
    lines = content.split('\n')
    fixed_lines = []
    seen_ns = set()
    
    for i, line in enumerate(lines):
        if '<rdf:RDF' in line:
            # 清理这一行的重复命名空间
            # 提取所有命名空间声明
            ns_pattern = r'(xmlns:\w+="[^"]+")'
            ns_declarations = re.findall(ns_pattern, line)
            xml_base_pattern = r'(xml:base="[^"]+")'
            xml_base = re.findall(xml_base_pattern, line)
            
            # 去重命名空间
            unique_ns = []
            for ns in ns_declarations:
                ns_name = ns.split('=')[0]
                if ns_name not in seen_ns:
                    unique_ns.append(ns)
                    seen_ns.add(ns_name)
            
            # 重建 RDF 元素行
            fixed_line = '<rdf:RDF'
            for ns in unique_ns:
                fixed_line += ' ' + ns
            for xb in xml_base:
                fixed_line += ' ' + xb
            fixed_line += '>'
            fixed_lines.append(fixed_line)
        else:
            fixed_lines.append(line)
    
    content = '\n'.join(fixed_lines)
    
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f"✅ 转换完成！")
    print(f"   输出文件: {output_file}")
    print(f"   类: {len(classes)}")
    print(f"   ObjectProperty: {len(obj_props)}")
    print(f"   DatatypeProperty: {len(datatype_props)}")


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("用法: python owl_to_rdfs_converter.py <owl_file> [output_file]")
        print("示例: python owl_to_rdfs_converter.py TheCimOntology.owl TheCimOntology_RDFS.rdf")
        sys.exit(1)
    
    owl_file = sys.argv[1]
    output_file = sys.argv[2] if len(sys.argv) > 2 else 'TheCimOntology_RDFS.rdf'
    
    convert_owl_to_rdfs(owl_file, output_file)

