#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
解析 CIM XML 文件，提取设备连接关系并生成可视化 HTML
"""

import xml.etree.ElementTree as ET
import json
import re
from collections import defaultdict
from pathlib import Path

def extract_id_from_resource(resource_str):
    """从 rdf:resource 或 rdf:ID 中提取 ID"""
    if not resource_str:
        return None
    # 处理 #_xxx 格式
    if resource_str.startswith('#'):
        return resource_str[1:]
    # 处理完整 URI
    if '#' in resource_str:
        return resource_str.split('#')[-1]
    return resource_str

def parse_cim_xml(xml_file_path):
    """解析 CIM XML 文件，提取设备连接关系"""
    
    # 定义命名空间
    namespaces = {
        'cim': 'http://iec.ch/TC57/2013/CIM-schema-cim16#',
        'rdf': 'http://www.w3.org/1999/02/22-rdf-syntax-ns#',
        'md': 'http://iec.ch/TC57/61970-552/ModelDescription/1#'
    }
    
    tree = ET.parse(xml_file_path)
    root = tree.getroot()
    
    # 存储数据结构
    equipment = {}  # {equipment_id: {type, name, ...}}
    terminals = {}  # {terminal_id: {equipment_id, connectivity_node_id, name, sequence}}
    connectivity_nodes = {}  # {node_id: {name, terminals: [terminal_ids]}}
    
    # 提取所有设备（ConductingEquipment 及其子类）
    for elem in root.findall('.//cim:*', namespaces):
        tag_name = elem.tag.split('}')[-1] if '}' in elem.tag else elem.tag
        
        # 检查是否是设备类型（ConductingEquipment 的子类）
        # 使用更通用的方法：检查标签名是否以常见设备类型开头
        # 或者直接检查是否是 ConductingEquipment 的子类
        equipment_types = [
            'ACLineSegment', 'SynchronousMachine', 'AsynchronousMachine',
            'TransformerWinding', 'PowerTransformer', 'Breaker', 'Load',
            'ConductingEquipment', 'BusbarSection', 'Disconnector',
            'EnergyConsumer', 'GeneratingUnit', 'ThermalGeneratingUnit',
            'ExternalNetworkInjection', 'ShuntCompensator', 'SeriesCompensator',
            'StaticVarCompensator', 'ReactiveCapabilityCurve'
        ]
        
        # 检查是否是设备类型，或者标签名包含 Equipment
        is_equipment = (tag_name in equipment_types or 
                       'Equipment' in tag_name or
                       tag_name.endswith('Machine') or
                       tag_name.endswith('Transformer') or
                       tag_name.endswith('Compensator'))
        
        if is_equipment:
            equipment_id = elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}ID')
            if equipment_id:
                name_elem = elem.find('cim:IdentifiedObject.name', namespaces)
                name = name_elem.text if name_elem is not None else equipment_id
                
                equipment[equipment_id] = {
                    'id': equipment_id,
                    'type': tag_name,
                    'name': name,
                    'terminals': []
                }
    
    # 提取所有 Terminal
    for terminal_elem in root.findall('.//cim:Terminal', namespaces):
        terminal_id = terminal_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}ID')
        if not terminal_id:
            # 尝试 rdf:about
            terminal_id = terminal_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}about')
            if terminal_id:
                terminal_id = extract_id_from_resource(terminal_id)
        
        if not terminal_id:
            continue
        
        # 获取名称
        name_elem = terminal_elem.find('cim:IdentifiedObject.name', namespaces)
        name = name_elem.text if name_elem is not None else terminal_id
        
        # 获取序列号
        seq_elem = terminal_elem.find('cim:ACDCTerminal.sequenceNumber', namespaces)
        sequence = int(seq_elem.text) if seq_elem is not None else 0
        
        # 获取关联的设备
        equipment_elem = terminal_elem.find('cim:Terminal.ConductingEquipment', namespaces)
        equipment_id = None
        if equipment_elem is not None:
            equipment_resource = equipment_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}resource')
            if equipment_resource:
                equipment_id = extract_id_from_resource(equipment_resource)
        
        # 获取关联的连接节点
        node_elem = terminal_elem.find('cim:Terminal.ConnectivityNode', namespaces)
        connectivity_node_id = None
        if node_elem is not None:
            node_resource = node_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}resource')
            if node_resource:
                connectivity_node_id = extract_id_from_resource(node_resource)
        
        if equipment_id:
            terminals[terminal_id] = {
                'id': terminal_id,
                'name': name,
                'sequence': sequence,
                'equipment_id': equipment_id,
                'connectivity_node_id': connectivity_node_id
            }
            
            # 添加到设备的端子列表
            if equipment_id in equipment:
                equipment[equipment_id]['terminals'].append(terminal_id)
    
    # 提取 ConnectivityNode
    for node_elem in root.findall('.//cim:ConnectivityNode', namespaces):
        node_id = node_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}ID')
        if not node_id:
            node_id = node_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}about')
            if node_id:
                node_id = extract_id_from_resource(node_id)
        
        if not node_id:
            continue
        
        name_elem = node_elem.find('cim:IdentifiedObject.name', namespaces)
        name = name_elem.text if name_elem is not None else node_id
        
        connectivity_nodes[node_id] = {
            'id': node_id,
            'name': name,
            'terminals': []
        }
    
    # 构建连接节点的端子列表
    for terminal_id, terminal_info in terminals.items():
        node_id = terminal_info.get('connectivity_node_id')
        if node_id and node_id in connectivity_nodes:
            connectivity_nodes[node_id]['terminals'].append(terminal_id)
    
    # 构建设备之间的连接关系
    connections = []  # [{from_equipment, to_equipment, via_node}]
    
    for node_id, node_info in connectivity_nodes.items():
        node_terminals = node_info['terminals']
        if len(node_terminals) < 2:
            continue
        
        # 找到连接到这个节点的所有设备
        connected_equipment = set()
        for term_id in node_terminals:
            if term_id in terminals:
                eq_id = terminals[term_id]['equipment_id']
                if eq_id:
                    connected_equipment.add(eq_id)
        
        # 为每对设备创建连接
        equipment_list = list(connected_equipment)
        for i in range(len(equipment_list)):
            for j in range(i + 1, len(equipment_list)):
                connections.append({
                    'from': equipment_list[i],
                    'to': equipment_list[j],
                    'via_node': node_id
                })
    
    return {
        'equipment': equipment,
        'terminals': terminals,
        'connectivity_nodes': connectivity_nodes,
        'connections': connections
    }

def parse_topology_xml(tp_xml_path, eq_data):
    """解析 TP XML 文件，提取拓扑节点和映射关系"""
    
    namespaces = {
        'cim': 'http://iec.ch/TC57/2013/CIM-schema-cim16#',
        'rdf': 'http://www.w3.org/1999/02/22-rdf-syntax-ns#',
        'md': 'http://iec.ch/TC57/61970-552/ModelDescription/1#'
    }
    
    tree = ET.parse(tp_xml_path)
    root = tree.getroot()
    
    topological_nodes = {}  # {topo_node_id: {name, base_voltage, connectivity_nodes: []}}
    connectivity_to_topological = {}  # {connectivity_node_id: topological_node_id}
    terminal_to_topological = {}  # {terminal_id: topological_node_id}
    
    # 提取 TopologicalNode
    for topo_node_elem in root.findall('.//cim:TopologicalNode', namespaces):
        topo_node_id = topo_node_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}ID')
        if not topo_node_id:
            continue
        
        name_elem = topo_node_elem.find('cim:IdentifiedObject.name', namespaces)
        name = name_elem.text if name_elem is not None else topo_node_id
        
        base_voltage_elem = topo_node_elem.find('cim:TopologicalNode.BaseVoltage', namespaces)
        base_voltage_id = None
        if base_voltage_elem is not None:
            base_voltage_resource = base_voltage_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}resource')
            if base_voltage_resource:
                base_voltage_id = extract_id_from_resource(base_voltage_resource)
        
        topological_nodes[topo_node_id] = {
            'id': topo_node_id,
            'name': name,
            'base_voltage_id': base_voltage_id,
            'connectivity_nodes': [],
            'terminals': []
        }
    
    # 提取 ConnectivityNode -> TopologicalNode 映射
    for conn_node_elem in root.findall('.//cim:ConnectivityNode', namespaces):
        conn_node_id = conn_node_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}about')
        if not conn_node_id:
            conn_node_id = conn_node_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}ID')
        if conn_node_id:
            conn_node_id = extract_id_from_resource(conn_node_id)
        
        topo_node_elem = conn_node_elem.find('cim:ConnectivityNode.TopologicalNode', namespaces)
        if topo_node_elem is not None:
            topo_node_resource = topo_node_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}resource')
            if topo_node_resource:
                topo_node_id = extract_id_from_resource(topo_node_resource)
                if conn_node_id and topo_node_id:
                    connectivity_to_topological[conn_node_id] = topo_node_id
                    if topo_node_id in topological_nodes:
                        topological_nodes[topo_node_id]['connectivity_nodes'].append(conn_node_id)
    
    # 提取 Terminal -> TopologicalNode 映射
    for terminal_elem in root.findall('.//cim:Terminal', namespaces):
        terminal_id = terminal_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}about')
        if not terminal_id:
            terminal_id = terminal_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}ID')
        if terminal_id:
            terminal_id = extract_id_from_resource(terminal_id)
        
        topo_node_elem = terminal_elem.find('cim:Terminal.TopologicalNode', namespaces)
        if topo_node_elem is not None:
            topo_node_resource = topo_node_elem.get('{http://www.w3.org/1999/02/22-rdf-syntax-ns#}resource')
            if topo_node_resource:
                topo_node_id = extract_id_from_resource(topo_node_resource)
                if terminal_id and topo_node_id:
                    terminal_to_topological[terminal_id] = topo_node_id
                    if topo_node_id in topological_nodes:
                        topological_nodes[topo_node_id]['terminals'].append(terminal_id)
    
    # 构建拓扑节点之间的连接（通过设备）
    topology_connections = []
    seen_connections = set()
    
    for topo_node_id, topo_info in topological_nodes.items():
        # 找到连接到这个拓扑节点的所有设备
        for term_id in topo_info['terminals']:
            if term_id in eq_data['terminals']:
                eq_id = eq_data['terminals'][term_id].get('equipment_id')
                if eq_id and eq_id in eq_data['equipment']:
                    # 找到这个设备的其他端子
                    eq_terminals = eq_data['equipment'][eq_id].get('terminals', [])
                    for other_term_id in eq_terminals:
                        if other_term_id != term_id and other_term_id in terminal_to_topological:
                            other_topo_id = terminal_to_topological[other_term_id]
                            if other_topo_id != topo_node_id:
                                # 避免重复连接
                                conn_key = tuple(sorted([topo_node_id, other_topo_id]))
                                if conn_key not in seen_connections:
                                    seen_connections.add(conn_key)
                                    topology_connections.append({
                                        'from': topo_node_id,
                                        'to': other_topo_id,
                                        'via_equipment': eq_id
                                    })
    
    return {
        'topological_nodes': topological_nodes,
        'connectivity_to_topological': connectivity_to_topological,
        'terminal_to_topological': terminal_to_topological,
        'topology_connections': topology_connections
    }

def generate_html(data, output_file, topology_data=None):
    """生成可视化 HTML 页面"""
    
    # 准备 vis.js 数据
    nodes = []
    edges = []
    
    # 设备类型颜色映射
    type_colors = {
        'ACLineSegment': '#4CAF50',
        'SynchronousMachine': '#2196F3',
        'AsynchronousMachine': '#FF9800',
        'PowerTransformer': '#9C27B0',
        'TransformerWinding': '#9C27B0',
        'Breaker': '#F44336',
        'Load': '#FFC107',
        'EnergyConsumer': '#FFC107',
        'BusbarSection': '#607D8B',
        'Disconnector': '#E91E63',
        'GeneratingUnit': '#00BCD4',
        'ThermalGeneratingUnit': '#00BCD4',
        'ExternalNetworkInjection': '#795548',
        'ConductingEquipment': '#757575',
        'default': '#757575'
    }
    
    # 添加设备节点
    equipment_map = {}
    for eq_id, eq_info in data['equipment'].items():
        eq_type = eq_info['type']
        color = type_colors.get(eq_type, type_colors['default'])
        
        node_id = f"eq_{eq_id}"
        equipment_map[eq_id] = node_id
        
        # 显示设备名称、类型和ID
        # 简化ID显示（只显示前8个字符和后8个字符，中间用...）
        eq_id_short = eq_id
        if len(eq_id) > 20:
            eq_id_short = f"{eq_id[:8]}...{eq_id[-8:]}"
        
        nodes.append({
            'id': node_id,
            'label': f"{eq_info['name']}\n({eq_type})\nID: {eq_id_short}",
            'group': eq_type,
            'color': color,
            'shape': 'box',
            'title': f"ID: {eq_id}\nType: {eq_type}\nName: {eq_info['name']}"
        })
    
    # 添加连接节点（可选，用于显示连接点）
    show_nodes = False  # 设置为 True 可以显示连接节点
    if show_nodes:
        for node_id, node_info in data['connectivity_nodes'].items():
            if len(node_info['terminals']) >= 2:  # 只显示有多个连接的节点
                nodes.append({
                    'id': f"node_{node_id}",
                    'label': node_info['name'] or node_id[:8],
                    'group': 'ConnectivityNode',
                    'color': '#E0E0E0',
                    'shape': 'dot',
                    'size': 10
                })
    
    # 添加连接边
    for conn in data['connections']:
        from_id = equipment_map.get(conn['from'])
        to_id = equipment_map.get(conn['to'])
        
        if from_id and to_id:
            edges.append({
                'from': from_id,
                'to': to_id,
                'arrows': 'to',
                'color': {'color': '#999'},
                'title': f"Via ConnectivityNode: {conn['via_node'][:8]}..."
            })
    
    # 准备拓扑网络数据（如果有）
    topology_nodes = []
    topology_edges = []
    mapping_edges = []  # ConnectivityNode 到 TopologicalNode 的映射边
    
    if topology_data:
        # 添加拓扑节点
        for topo_id, topo_info in topology_data['topological_nodes'].items():
            topology_nodes.append({
                'id': f"topo_{topo_id}",
                'label': f"{topo_info['name']}\n(Topo)",
                'group': 'TopologicalNode',
                'color': '#FF6B6B',
                'shape': 'ellipse',
                'title': f"TopologicalNode: {topo_info['name']}\nID: {topo_id}\nConnectivityNodes: {len(topo_info['connectivity_nodes'])}"
            })
        
        # 添加拓扑节点之间的连接
        for conn in topology_data['topology_connections']:
            topology_edges.append({
                'from': f"topo_{conn['from']}",
                'to': f"topo_{conn['to']}",
                'arrows': 'to',
                'color': {'color': '#4ECDC4'},
                'width': 3,
                'title': f"Via Equipment"
            })
        
        # 添加 ConnectivityNode 到 TopologicalNode 的映射边（虚线）
        # 根据连接的设备类型给 ConnectivityNode 着色
        conn_node_set = set()
        conn_node_equipment_types = {}  # {conn_node_id: [equipment_types]}
        
        # 统计每个 ConnectivityNode 连接的设备类型
        for conn_node_id, topo_node_id in topology_data['connectivity_to_topological'].items():
            conn_node_set.add(conn_node_id)
            equipment_types = []
            
            # 找到连接到这个 ConnectivityNode 的所有设备
            if conn_node_id in data['connectivity_nodes']:
                terminals = data['connectivity_nodes'][conn_node_id].get('terminals', [])
                for term_id in terminals:
                    if term_id in data['terminals']:
                        eq_id = data['terminals'][term_id].get('equipment_id')
                        if eq_id and eq_id in data['equipment']:
                            eq_type = data['equipment'][eq_id]['type']
                            if eq_type not in equipment_types:
                                equipment_types.append(eq_type)
            
            conn_node_equipment_types[conn_node_id] = equipment_types
        
        # 在拓扑网络中显示 ConnectivityNode，根据设备类型着色
        for conn_node_id, topo_node_id in topology_data['connectivity_to_topological'].items():
            # 在拓扑网络中显示 ConnectivityNode（小节点）
            if conn_node_id not in [n['id'].replace('conn_', '') for n in topology_nodes]:
                equipment_types = conn_node_equipment_types.get(conn_node_id, [])
                
                # 确定颜色：如果有多个设备类型，使用第一个；如果没有，使用默认灰色
                if equipment_types:
                    primary_type = equipment_types[0]
                    color = type_colors.get(primary_type, type_colors['default'])
                    # 如果有多个类型，在标签中显示
                    type_label = primary_type if len(equipment_types) == 1 else f"{primary_type}+{len(equipment_types)-1}"
                else:
                    color = '#BDC3C7'
                    type_label = "Unknown"
                
                # 获取 ConnectivityNode 的名称（如果有）
                conn_name = ""
                if conn_node_id in data['connectivity_nodes']:
                    conn_name = data['connectivity_nodes'][conn_node_id].get('name', '')
                
                # 优先显示设备类型，使标签更清晰
                if equipment_types:
                    # 如果有多个类型，显示主要类型和数量
                    if len(equipment_types) > 1:
                        label = f"{primary_type}\n+{len(equipment_types)-1}"
                    else:
                        label = primary_type
                else:
                    # 如果没有设备类型信息，显示名称或ID
                    label = conn_name[:8] if conn_name else conn_node_id[:8]
                
                topology_nodes.append({
                    'id': f"conn_{conn_node_id}",
                    'label': label,
                    'group': 'ConnectivityNode',
                    'color': color,
                    'shape': 'dot',
                    'size': 10,
                    'title': f"ConnectivityNode: {conn_node_id}\n设备类型: {', '.join(equipment_types) if equipment_types else 'Unknown'}\n映射到: TopologicalNode {topology_data['topological_nodes'][topo_node_id]['name']}"
                })
                
                # 添加映射边（虚线）
                mapping_edges.append({
                    'from': f"conn_{conn_node_id}",
                    'to': f"topo_{topo_node_id}",
                    'arrows': 'to',
                    'color': {'color': '#95A5A6'},
                    'dashes': True,
                    'width': 1,
                    'title': f"ConnectivityNode → TopologicalNode"
                })
    
    # 生成 HTML
    html_content = f"""<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CIM 设备连接关系可视化</title>
    <script type="text/javascript" src="https://unpkg.com/vis-network/standalone/umd/vis-network.min.js"></script>
    <style>
        body {{
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }}
        #header {{
            background: white;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 20px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }}
        h1 {{
            margin: 0 0 10px 0;
            color: #333;
        }}
        .stats {{
            display: flex;
            gap: 20px;
            margin-top: 10px;
        }}
        .stat-item {{
            padding: 10px;
            background: #f0f0f0;
            border-radius: 4px;
        }}
        .stat-label {{
            font-size: 12px;
            color: #666;
        }}
        .stat-value {{
            font-size: 20px;
            font-weight: bold;
            color: #2196F3;
        }}
        .network-container {{
            display: flex;
            gap: 20px;
            margin-bottom: 20px;
        }}
        .network-panel {{
            flex: 1;
            background: white;
            border-radius: 8px;
            padding: 15px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }}
        .network-panel h2 {{
            margin: 0 0 10px 0;
            font-size: 18px;
            color: #333;
        }}
        #network {{
            width: 100%;
            height: 600px;
            border: 1px solid #ddd;
            border-radius: 8px;
            background: white;
        }}
        #topology-network {{
            width: 100%;
            height: 600px;
            border: 1px solid #ddd;
            border-radius: 8px;
            background: white;
        }}
        #legend {{
            background: white;
            padding: 15px;
            border-radius: 8px;
            margin-top: 20px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }}
        .legend-item {{
            display: inline-block;
            margin: 5px 15px;
            padding: 5px 10px;
            border-radius: 4px;
        }}
        .legend-color {{
            display: inline-block;
            width: 20px;
            height: 20px;
            border-radius: 3px;
            margin-right: 8px;
            vertical-align: middle;
        }}
    </style>
</head>
<body>
    <div id="header">
        <h1>🔌 CIM 设备连接关系可视化</h1>
        <div class="stats">
            <div class="stat-item">
                <div class="stat-label">设备总数</div>
                <div class="stat-value">{len(data['equipment'])}</div>
            </div>
            <div class="stat-item">
                <div class="stat-label">端子总数</div>
                <div class="stat-value">{len(data['terminals'])}</div>
            </div>
            <div class="stat-item">
                <div class="stat-label">连接节点</div>
                <div class="stat-value">{len(data['connectivity_nodes'])}</div>
            </div>
            <div class="stat-item">
                <div class="stat-label">连接关系</div>
                <div class="stat-value">{len(data['connections'])}</div>
            </div>
        </div>
    </div>
    
    <div class="network-container">
        <div class="network-panel">
            <h2>🔌 物理设备网络 (Node-Breaker Model)</h2>
            <div id="network"></div>
        </div>
        {f'''
        <div class="network-panel">
            <h2>🌐 拓扑网络 (Bus-Branch Model)</h2>
            <div id="topology-network"></div>
        </div>
        ''' if topology_data else ''}
    </div>
    
    <div id="legend">
        <h3>图例</h3>
        {generate_legend(type_colors, data['equipment'])}
    </div>
    
    <script type="text/javascript">
        // 节点数据
        var nodes = new vis.DataSet({json.dumps(nodes, ensure_ascii=False, indent=2)});
        
        // 边数据
        var edges = new vis.DataSet({json.dumps(edges, ensure_ascii=False, indent=2)});
        
        // 网络配置
        var options = {{
            nodes: {{
                borderWidth: 2,
                shadow: true,
                font: {{
                    size: 12,
                    face: 'Arial'
                }}
            }},
            edges: {{
                width: 2,
                smooth: {{
                    type: 'continuous',
                    roundness: 0.5
                }},
                arrows: {{
                    to: {{
                        enabled: true,
                        scaleFactor: 0.5
                    }}
                }}
            }},
            physics: {{
                enabled: true,
                stabilization: {{
                    iterations: 200
                }},
                barnesHut: {{
                    gravitationalConstant: -2000,
                    centralGravity: 0.3,
                    springLength: 200,
                    springConstant: 0.04,
                    damping: 0.09
                }}
            }},
            interaction: {{
                hover: true,
                tooltipDelay: 100,
                zoomView: true,
                dragView: true
            }}
        }};
        
        // 创建物理设备网络
        var container = document.getElementById('network');
        var networkData = {{
            nodes: nodes,
            edges: edges
        }};
        var network = new vis.Network(container, networkData, options);
        
        // 添加事件监听
        network.on("click", function (params) {{
            if (params.nodes.length > 0) {{
                console.log("选中节点:", params.nodes[0]);
            }}
        }});
        
        {f'''
        // 创建拓扑网络
        var topologyNodes = new vis.DataSet({json.dumps(topology_nodes, ensure_ascii=False, indent=2)});
        var topologyEdges = new vis.DataSet({json.dumps(topology_edges, ensure_ascii=False, indent=2)});
        var mappingEdges = new vis.DataSet({json.dumps(mapping_edges, ensure_ascii=False, indent=2)});
        
        // 合并所有边
        var allTopologyEdges = new vis.DataSet([...topologyEdges.get(), ...mappingEdges.get()]);
        
        var topologyContainer = document.getElementById('topology-network');
        var topologyNetworkData = {{
            nodes: topologyNodes,
            edges: allTopologyEdges
        }};
        var topologyNetwork = new vis.Network(topologyContainer, topologyNetworkData, options);
        
        topologyNetwork.on("click", function (params) {{
            if (params.nodes.length > 0) {{
                console.log("选中拓扑节点:", params.nodes[0]);
            }}
        }});
        ''' if topology_data else ''}
    </script>
</body>
</html>"""
    
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(html_content)
    
    print(f"✅ HTML 文件已生成: {output_file}")

def generate_legend(type_colors, equipment):
    """生成图例 HTML"""
    # 统计实际使用的设备类型
    used_types = set()
    for eq_info in equipment.values():
        used_types.add(eq_info['type'])
    
    legend_items = []
    for eq_type in sorted(used_types):
        color = type_colors.get(eq_type, type_colors['default'])
        legend_items.append(
            f'<div class="legend-item">'
            f'<span class="legend-color" style="background-color: {color};"></span>'
            f'<span>{eq_type}</span>'
            f'</div>'
        )
    
    return ''.join(legend_items)

def main():
    xml_file = Path(__file__).parent / 'datasets' / 'MiniGrid_NodeBreaker' / 'MiniGridTestConfiguration_BC_EQ_v3.0.0.xml'
    tp_file = Path(__file__).parent / 'datasets' / 'MiniGrid_NodeBreaker' / 'MiniGridTestConfiguration_BC_TP_v3.0.0.xml'
    output_file = Path(__file__).parent / 'datasets' / 'MiniGrid_NodeBreaker' / 'equipment_network.html'
    
    print(f"📖 正在解析 CIM EQ XML 文件: {xml_file}")
    data = parse_cim_xml(xml_file)
    
    print(f"📊 EQ 解析结果:")
    print(f"  - 设备数量: {len(data['equipment'])}")
    print(f"  - 端子数量: {len(data['terminals'])}")
    print(f"  - 连接节点数量: {len(data['connectivity_nodes'])}")
    print(f"  - 连接关系数量: {len(data['connections'])}")
    
    topology_data = None
    if tp_file.exists():
        print(f"\n📖 正在解析 CIM TP XML 文件: {tp_file}")
        topology_data = parse_topology_xml(tp_file, data)
        print(f"📊 TP 解析结果:")
        print(f"  - 拓扑节点数量: {len(topology_data['topological_nodes'])}")
        print(f"  - ConnectivityNode映射数量: {len(topology_data['connectivity_to_topological'])}")
        print(f"  - Terminal映射数量: {len(topology_data['terminal_to_topological'])}")
        print(f"  - 拓扑连接数量: {len(topology_data['topology_connections'])}")
    else:
        print(f"\n⚠️  未找到 TP 文件，跳过拓扑网络可视化")
    
    print(f"\n🎨 正在生成可视化 HTML...")
    generate_html(data, output_file, topology_data)
    
    print(f"\n✨ 完成！请在浏览器中打开: {output_file}")

if __name__ == '__main__':
    main()

