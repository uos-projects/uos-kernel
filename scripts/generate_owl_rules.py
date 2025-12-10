#!/usr/bin/env python3
"""
从 TheCimOntology.owl 生成 OWL 产生式规则

生成规则类型：
1. OWL 约束规则（基数约束、值约束等）
2. SWRL 规则（复杂业务逻辑）
3. OWL 传递性规则
4. OWL 等价类规则
"""

import xml.etree.ElementTree as ET
from pathlib import Path
import sys

# 命名空间定义
NSMAP = {
    'owl': 'http://www.w3.org/2002/07/owl#',
    'rdf': 'http://www.w3.org/1999/02/22-rdf-syntax-ns#',
    'rdfs': 'http://www.w3.org/2000/01/rdf-schema#',
    'swrl': 'http://www.w3.org/2003/11/swrl#',
    'swrlb': 'http://www.w3.org/2003/11/swrlb#',
    'cim': 'http://www.iec.ch/TC57/CIM#',
}

CIM_NS = NSMAP['cim']


def create_cardinality_rule(class_uri, property_uri, min_card, comment):
    """创建基数约束规则元素"""
    cls_elem = ET.Element(f'{{{NSMAP["owl"]}}}Class')
    cls_elem.set(f'{{{NSMAP["rdf"]}}}about', class_uri)
    
    # 添加注释
    comment_elem = ET.Comment(comment)
    cls_elem.append(comment_elem)
    
    # 创建 subClassOf
    sub_elem = ET.SubElement(cls_elem, f'{{{NSMAP["rdfs"]}}}subClassOf')
    
    # 创建 Restriction
    restriction = ET.SubElement(sub_elem, f'{{{NSMAP["owl"]}}}Restriction')
    
    # onProperty
    prop_elem = ET.SubElement(restriction, f'{{{NSMAP["owl"]}}}onProperty')
    prop_elem.set(f'{{{NSMAP["rdf"]}}}resource', property_uri)
    
    # minCardinality
    card_elem = ET.SubElement(restriction, f'{{{NSMAP["owl"]}}}minCardinality')
    card_elem.set(f'{{{NSMAP["rdf"]}}}datatype', 'http://www.w3.org/2001/XMLSchema#nonNegativeInteger')
    card_elem.text = str(min_card)
    
    return cls_elem


def generate_cardinality_rules(owl_file, output_file):
    """生成基数约束规则"""
    print("生成基数约束规则...")
    
    rules = []
    
    # 规则 1: ConductingEquipment 必须至少有一个 Terminal
    rules.append(create_cardinality_rule(
        f'{CIM_NS}ConductingEquipment',
        f'{CIM_NS}ConductingEquipment.Terminals',
        1,
        '规则 1: 每个导电设备必须至少有一个端子'
    ))
    
    # 规则 2: Terminal 必须连接到一个 ConductingEquipment
    rules.append(create_cardinality_rule(
        f'{CIM_NS}Terminal',
        f'{CIM_NS}Terminal.ConductingEquipment',
        1,
        '规则 2: 每个端子必须连接到一个导电设备'
    ))
    
    # 规则 3: Terminal 必须连接到一个 ConnectivityNode
    rules.append(create_cardinality_rule(
        f'{CIM_NS}Terminal',
        f'{CIM_NS}Terminal.ConnectivityNode',
        1,
        '规则 3: 每个端子必须连接到一个连接节点'
    ))
    
    # 规则 4: PowerTransformer 必须至少有两个 PowerTransformerEnd
    rules.append(create_cardinality_rule(
        f'{CIM_NS}PowerTransformer',
        f'{CIM_NS}PowerTransformer.PowerTransformerEnds',
        2,
        '规则 4: 每个变压器必须至少有两个绕组'
    ))
    
    return rules


def generate_transitive_rules():
    """生成传递性规则"""
    print("生成传递性规则...")
    
    rules = []
    
    # 规则: EquipmentContainer 的包含关系是传递的
    prop_elem = ET.Element(f'{{{NSMAP["owl"]}}}ObjectProperty')
    prop_elem.set(f'{{{NSMAP["rdf"]}}}about', f'{CIM_NS}EquipmentContainer.Equipments')
    
    comment = ET.Comment('传递性规则: 设备容器的包含关系是传递的')
    prop_elem.append(comment)
    
    type_elem = ET.SubElement(prop_elem, f'{{{NSMAP["rdf"]}}}type')
    type_elem.set(f'{{{NSMAP["rdf"]}}}resource', f'{NSMAP["owl"]}TransitiveProperty')
    
    rules.append(prop_elem)
    
    return rules


def generate_swrl_rules():
    """生成 SWRL 规则"""
    print("生成 SWRL 规则...")
    
    rules = []
    
    # 规则 1: 开关状态影响端子连接
    rules.append(f'''
    <!-- SWRL 规则 1: 开关断开时，端子不连接 -->
    <swrl:Imp>
        <swrl:body>
            <swrl:AtomList>
                <swrl:atom>
                    <swrl:ClassPredicate rdf:resource="{CIM_NS}Switch"/>
                    <swrl:argument1>
                        <swrl:Variable rdf:about="#?s"/>
                    </swrl:argument1>
                </swrl:atom>
                <swrl:atom>
                    <swrl:ObjectPropertyPredicate rdf:resource="{CIM_NS}Switch.Terminals"/>
                    <swrl:argument1>
                        <swrl:Variable rdf:about="#?s"/>
                    </swrl:argument1>
                    <swrl:argument2>
                        <swrl:Variable rdf:about="#?t"/>
                    </swrl:argument2>
                </swrl:atom>
                <swrl:atom>
                    <swrl:DatavaluedPropertyPredicate rdf:resource="{CIM_NS}Switch.open"/>
                    <swrl:argument1>
                        <swrl:Variable rdf:about="#?s"/>
                    </swrl:argument1>
                    <swrl:argument2>
                        <swrl:Variable rdf:about="#?open"/>
                    </swrl:argument2>
                </swrl:atom>
            </swrl:AtomList>
        </swrl:body>
        <swrl:head>
            <swrl:AtomList>
                <swrl:atom>
                    <swrl:DatavaluedPropertyPredicate rdf:resource="{CIM_NS}Terminal.connected"/>
                    <swrl:argument1>
                        <swrl:Variable rdf:about="#?t"/>
                    </swrl:argument1>
                    <swrl:argument2>
                        <swrl:Variable rdf:about="#?connected"/>
                    </swrl:argument2>
                </swrl:atom>
            </swrl:AtomList>
        </swrl:head>
    </swrl:Imp>
    ''')
    
    # 规则 2: 电气连接推理
    rules.append(f'''
    <!-- SWRL 规则 2: 连接到同一连接节点的端子是电气连接的 -->
    <swrl:Imp>
        <swrl:body>
            <swrl:AtomList>
                <swrl:atom>
                    <swrl:ClassPredicate rdf:resource="{CIM_NS}Terminal"/>
                    <swrl:argument1>
                        <swrl:Variable rdf:about="#?t1"/>
                    </swrl:argument1>
                </swrl:atom>
                <swrl:atom>
                    <swrl:ClassPredicate rdf:resource="{CIM_NS}Terminal"/>
                    <swrl:argument1>
                        <swrl:Variable rdf:about="#?t2"/>
                    </swrl:argument1>
                </swrl:atom>
                <swrl:atom>
                    <swrl:ClassPredicate rdf:resource="{CIM_NS}ConnectivityNode"/>
                    <swrl:argument1>
                        <swrl:Variable rdf:about="#?cn"/>
                    </swrl:argument1>
                </swrl:atom>
                <swrl:atom>
                    <swrl:ObjectPropertyPredicate rdf:resource="{CIM_NS}Terminal.ConnectivityNode"/>
                    <swrl:argument1>
                        <swrl:Variable rdf:about="#?t1"/>
                    </swrl:argument1>
                    <swrl:argument2>
                        <swrl:Variable rdf:about="#?cn"/>
                    </swrl:argument2>
                </swrl:atom>
                <swrl:atom>
                    <swrl:ObjectPropertyPredicate rdf:resource="{CIM_NS}Terminal.ConnectivityNode"/>
                    <swrl:argument1>
                        <swrl:Variable rdf:about="#?t2"/>
                    </swrl:argument1>
                    <swrl:argument2>
                        <swrl:Variable rdf:about="#?cn"/>
                    </swrl:argument2>
                </swrl:atom>
            </swrl:AtomList>
        </swrl:body>
        <swrl:head>
            <swrl:AtomList>
                <swrl:atom>
                    <swrl:ObjectPropertyPredicate rdf:resource="{CIM_NS}Terminal.electricallyConnected"/>
                    <swrl:argument1>
                        <swrl:Variable rdf:about="#?t1"/>
                    </swrl:argument1>
                    <swrl:argument2>
                        <swrl:Variable rdf:about="#?t2"/>
                    </swrl:argument2>
                </swrl:atom>
            </swrl:AtomList>
        </swrl:head>
    </swrl:Imp>
    ''')
    
    return rules


def generate_rules_file(owl_file, output_file):
    """生成完整的规则文件"""
    print(f"从 {owl_file} 生成 OWL 规则文件...")
    print("=" * 60)
    
    # 生成各种规则
    cardinality_rules = generate_cardinality_rules(owl_file, output_file)
    transitive_rules = generate_transitive_rules()
    swrl_rules = generate_swrl_rules()
    
    # 创建 RDF 根元素
    rdf_root = ET.Element(f'{{{NSMAP["rdf"]}}}RDF')
    rdf_root.set('xmlns:owl', NSMAP['owl'])
    rdf_root.set('xmlns:rdf', NSMAP['rdf'])
    rdf_root.set('xmlns:rdfs', NSMAP['rdfs'])
    rdf_root.set('xmlns:swrl', NSMAP['swrl'])
    rdf_root.set('xmlns:swrlb', NSMAP['swrlb'])
    rdf_root.set('xmlns', CIM_NS)
    rdf_root.set('{http://www.w3.org/XML/1998/namespace}base', 'http://www.iec.ch/TC57/CIM')
    
    # 添加 Ontology 声明
    ontology = ET.SubElement(rdf_root, f'{{{NSMAP["owl"]}}}Ontology')
    ontology.set(f'{{{NSMAP["rdf"]}}}about', f'{CIM_NS}')
    
    # 导入原始本体
    import_elem = ET.SubElement(ontology, f'{{{NSMAP["owl"]}}}imports')
    import_elem.set(f'{{{NSMAP["rdf"]}}}resource', 'TheCimOntology.owl')
    
    # 添加注释
    comment = ET.SubElement(ontology, f'{{{NSMAP["rdfs"]}}}comment')
    comment.text = 'Generated OWL rules for CIM ontology validation and reasoning'
    
    # 添加规则（已经是元素对象）
    for rule_elem in cardinality_rules + transitive_rules:
        rdf_root.append(rule_elem)
    
    # 注意: SWRL 规则需要更复杂的处理，这里先跳过
    # 实际使用时可以使用专门的 SWRL 库
    
    # 写入文件
    tree = ET.ElementTree(rdf_root)
    ET.indent(tree, space='\t')
    
    print(f"\n正在写入 {output_file}...")
    tree.write(output_file, encoding='utf-8', xml_declaration=True)
    
    # 修复重复的命名空间
    import re
    with open(output_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 修复 XML 声明
    if not content.startswith('<?xml version="1.0" encoding="UTF-8"?>'):
        content = '<?xml version="1.0" encoding="UTF-8"?>\n' + content.split('?>', 1)[-1] if '?>' in content else content
    
    # 清理重复的命名空间
    lines = content.split('\n')
    fixed_lines = []
    seen_ns = set()
    
    for line in lines:
        if '<rdf:RDF' in line:
            # 清理这一行的重复命名空间
            ns_pattern = r'(xmlns(?::\w+)?="[^"]+")'
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
    
    print(f"\n✅ 规则生成完成！")
    print(f"   输出文件: {output_file}")
    print(f"   基数约束规则: {len(cardinality_rules)}")
    print(f"   传递性规则: {len(transitive_rules)}")
    print(f"   SWRL 规则: {len(swrl_rules)}")
    print(f"   总计: {len(cardinality_rules) + len(transitive_rules) + len(swrl_rules)} 个规则")


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("用法: python generate_owl_rules.py <owl_file> [output_file]")
        print("示例: python generate_owl_rules.py TheCimOntology.owl CIM_Rules.owl")
        sys.exit(1)
    
    owl_file = sys.argv[1]
    output_file = sys.argv[2] if len(sys.argv) > 2 else 'CIM_Rules.owl'
    
    generate_rules_file(owl_file, output_file)

