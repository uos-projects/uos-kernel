#!/usr/bin/env python3
"""
测试 cimpyorm 库的使用
"""

import sys
from pathlib import Path

# 添加 cimpyorm 到路径
sys.path.insert(0, str(Path(__file__).parent.parent / 'cimpyorm'))

from cimpyorm import parse, describe, stats, lint
from cimpyorm.backends import InMemory, SQLite


def test_minigrid_dataset():
    """测试 MiniGrid_BusBranch 数据集"""
    print("=" * 60)
    print("测试 cimpyorm - MiniGrid_BusBranch 数据集")
    print("=" * 60)
    
    dataset_path = Path(__file__).parent.parent / 'cimpyorm' / 'cimpyorm' / 'res' / 'datasets' / 'MiniGrid_BusBranch'
    
    print(f"\n📂 数据集路径: {dataset_path}")
    print("开始解析...")
    
    # 使用内存数据库（更快）
    session, model = parse(str(dataset_path), backend=InMemory, silence_tqdm=True)
    
    print("✅ 数据集解析成功！")
    
    # 基本统计
    print("\n📊 数据集统计:")
    print(f"  - Equipment: {session.query(model.Equipment).count()}")
    print(f"  - Substation: {session.query(model.Substation).count()}")
    print(f"  - VoltageLevel: {session.query(model.VoltageLevel).count()}")
    print(f"  - Terminal: {session.query(model.Terminal).count()}")
    print(f"  - ACLineSegment: {session.query(model.ACLineSegment).count()}")
    print(f"  - PowerTransformer: {session.query(model.PowerTransformer).count()}")
    print(f"  - ConnectivityNode: {session.query(model.ConnectivityNode).count()}")
    
    # 查询变电站
    print("\n🏭 变电站信息:")
    substations = session.query(model.Substation).all()
    for sub in substations:
        print(f"\n  {sub.name}")
        if hasattr(sub, 'VoltageLevels'):
            for vl in sub.VoltageLevels:
                print(f"    └─ {vl.name}")
                # 查询该电压等级下的设备
                if hasattr(vl, 'Equipments'):
                    equip_count = len(list(vl.Equipments))
                    print(f"       设备数量: {equip_count}")
    
    # 查询线路
    print("\n🔌 交流线路段 (ACLineSegment):")
    ac_lines = session.query(model.ACLineSegment).limit(5).all()
    for line in ac_lines:
        print(f"  - {line.name}")
        print(f"    电阻 r={line.r if hasattr(line, 'r') else 'N/A'} Ω")
        print(f"    电抗 x={line.x if hasattr(line, 'x') else 'N/A'} Ω")
        if hasattr(line, 'BaseVoltage') and line.BaseVoltage:
            print(f"    基准电压: {line.BaseVoltage.name}")
        if hasattr(line, 'Terminals'):
            terminals = list(line.Terminals)
            print(f"    端子数: {len(terminals)}")
    
    # 查询变压器
    print("\n� transformer 变压器:")
    transformers = session.query(model.PowerTransformer).limit(3).all()
    for trans in transformers:
        print(f"  - {trans.name}")
        if hasattr(trans, 'PowerTransformerEnds'):
            ends = list(trans.PowerTransformerEnds)
            print(f"    绕组数: {len(ends)}")
            for end in ends:
                if hasattr(end, 'ratedU'):
                    print(f"      └─ 额定电压: {end.ratedU} kV")
    
    # 查询端子连接关系
    print("\n🔗 端子连接关系示例:")
    terminals = session.query(model.Terminal).limit(5).all()
    for term in terminals:
        if hasattr(term, 'ConductingEquipment') and term.ConductingEquipment:
            equip = term.ConductingEquipment
            print(f"  Terminal {term.sequenceNumber if hasattr(term, 'sequenceNumber') else 'N/A'} -> {equip.name if hasattr(equip, 'name') else type(equip).__name__}")
            if hasattr(term, 'ConnectivityNode') and term.ConnectivityNode:
                cn = term.ConnectivityNode
                print(f"    └─ 连接到: {cn.name if hasattr(cn, 'name') else 'ConnectivityNode'}")
    
    # 使用 stats 函数
    print("\n📈 详细统计信息:")
    try:
        stats_df = stats(session)
        print(stats_df.head(10))
    except Exception as e:
        print(f"统计功能出错: {e}")
    
    return session, model


def test_fullgrid_dataset():
    """测试 FullGrid 数据集（较大）"""
    print("\n" + "=" * 60)
    print("测试 FullGrid 数据集")
    print("=" * 60)
    
    dataset_path = Path(__file__).parent.parent / 'cimpyorm' / 'cimpyorm' / 'res' / 'datasets' / 'FullGrid'
    
    print(f"\n📂 数据集路径: {dataset_path}")
    print("开始解析（这可能需要一些时间）...")
    
    try:
        session, model = parse(str(dataset_path), backend=InMemory, silence_tqdm=True)
        
        print("✅ FullGrid 数据集解析成功！")
        print(f"\n📊 FullGrid 统计:")
        print(f"  - Equipment: {session.query(model.Equipment).count()}")
        print(f"  - Substation: {session.query(model.Substation).count()}")
        print(f"  - Terminal: {session.query(model.Terminal).count()}")
        
    except Exception as e:
        print(f"❌ FullGrid 解析失败: {e}")
        import traceback
        traceback.print_exc()


if __name__ == '__main__':
    # 测试 MiniGrid
    session, model = test_minigrid_dataset()
    
    # 可选：测试 FullGrid（较大，可能较慢）
    # test_fullgrid_dataset()
    
    print("\n" + "=" * 60)
    print("✅ 测试完成！")
    print("=" * 60)

