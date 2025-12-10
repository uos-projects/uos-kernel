#!/usr/bin/env python3
"""
使用 cimpyorm 加载 FullGrid 数据集
"""
import sys
from pathlib import Path

# 添加 cimpyorm 到 sys.path
project_root = Path(__file__).parent.parent
sys.path.insert(0, str(project_root / 'cimpyorm'))

import cimpyorm
from cimpyorm import parse, load, stats
from cimpyorm.backends import SQLite

def main():
    # FullGrid 数据集路径
    fullgrid_path = project_root / 'cimpyorm' / 'cimpyorm' / 'res' / 'datasets' / 'FullGrid'
    db_path = fullgrid_path / 'out.db'
    
    print("=" * 80)
    print("加载 FullGrid 数据集")
    print("=" * 80)
    print(f"数据集路径: {fullgrid_path}")
    print(f"数据库路径: {db_path}")
    print()
    
    # 检查数据库是否已存在
    if db_path.exists():
        print("发现已存在的数据库，尝试加载...")
        try:
            session, model = load(str(db_path))
            print("✓ 成功加载已存在的数据库")
        except Exception as e:
            print(f"✗ 加载失败: {e}")
            print("重新解析数据集...")
            session, model = parse(str(fullgrid_path))
    else:
        print("数据库不存在，开始解析数据集...")
        session, model = parse(str(fullgrid_path))
        print("✓ 解析完成")
    
    print()
    print("=" * 80)
    print("数据集统计信息")
    print("=" * 80)
    
    # 获取统计信息
    stats_df = stats(session)
    print(stats_df.head(20))
    
    print()
    print("=" * 80)
    print("示例查询")
    print("=" * 80)
    
    # 查询一些示例数据
    # 查询 Measurement 相关类
    if hasattr(model, 'Accumulator'):
        accumulators = session.query(model.Accumulator).limit(5).all()
        print(f"\nAccumulator (累加器) 数量: {session.query(model.Accumulator).count()}")
        if accumulators:
            print(f"示例: {accumulators[0].name}")
    
    if hasattr(model, 'Analog'):
        analogs = session.query(model.Analog).limit(5).all()
        print(f"\nAnalog (模拟量) 数量: {session.query(model.Analog).count()}")
        if analogs:
            print(f"示例: {analogs[0].name}")
    
    if hasattr(model, 'Discrete'):
        discretes = session.query(model.Discrete).limit(5).all()
        print(f"\nDiscrete (离散量) 数量: {session.query(model.Discrete).count()}")
        if discretes:
            print(f"示例: {discretes[0].name}")
    
    if hasattr(model, 'StringMeasurement'):
        string_measurements = session.query(model.StringMeasurement).limit(5).all()
        print(f"\nStringMeasurement (字符串测量) 数量: {session.query(model.StringMeasurement).count()}")
        if string_measurements:
            print(f"示例: {string_measurements[0].name}")
    
    # 查询 Control 相关类
    if hasattr(model, 'Command'):
        commands = session.query(model.Command).limit(5).all()
        print(f"\nCommand (命令) 数量: {session.query(model.Command).count()}")
        if commands:
            print(f"示例: {commands[0].name}")
    
    if hasattr(model, 'AccumulatorReset'):
        accumulator_resets = session.query(model.AccumulatorReset).limit(5).all()
        print(f"\nAccumulatorReset (累加器复位) 数量: {session.query(model.AccumulatorReset).count()}")
        if accumulator_resets:
            print(f"示例: {accumulator_resets[0].name}")
    
    if hasattr(model, 'RaiseLowerCommand'):
        raise_lower_commands = session.query(model.RaiseLowerCommand).limit(5).all()
        print(f"\nRaiseLowerCommand (升降命令) 数量: {session.query(model.RaiseLowerCommand).count()}")
        if raise_lower_commands:
            print(f"示例: {raise_lower_commands[0].name}")
    
    print()
    print("=" * 80)
    print("cimpyorm 已成功启动！")
    print("=" * 80)
    print("\n提示: 可以在 Python 交互式环境中使用以下代码继续操作:")
    print(f"  from cimpyorm import load")
    print(f"  session, model = load('{db_path}')")
    print()
    
    return session, model

if __name__ == '__main__':
    session, model = main()


