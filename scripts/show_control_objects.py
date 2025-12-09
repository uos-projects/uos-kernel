#!/usr/bin/env python3
"""
展示 Control 对象及其关联的管理对象（PowerSystemResource）
"""
import sys
from pathlib import Path

# 添加 cimpyorm 到 sys.path
project_root = Path(__file__).parent.parent
sys.path.insert(0, str(project_root / 'cimpyorm'))

from cimpyorm import load

def print_control_info(control, session, model, indent=""):
    """打印 Control 对象的详细信息"""
    print(f"{indent}名称: {control.name}")
    print(f"{indent}mRID: {control.mRID}")
    if hasattr(control, 'description') and control.description:
        print(f"{indent}描述: {control.description}")
    if hasattr(control, 'controlType') and control.controlType:
        print(f"{indent}控制类型: {control.controlType}")
    if hasattr(control, 'timeStamp') and control.timeStamp:
        print(f"{indent}时间戳: {control.timeStamp}")
    if hasattr(control, 'operationInProgress') and control.operationInProgress is not None:
        print(f"{indent}操作进行中: {control.operationInProgress}")
    if hasattr(control, 'unitMultiplier') and control.unitMultiplier:
        unit_mult = control.unitMultiplier
        if hasattr(unit_mult, 'name'):
            print(f"{indent}单位倍数: {unit_mult.name}")
    if hasattr(control, 'unitSymbol') and control.unitSymbol:
        unit_sym = control.unitSymbol
        if hasattr(unit_sym, 'name'):
            print(f"{indent}单位符号: {unit_sym.name}")
    
    # 查找关联的 PowerSystemResource（管理对象）
    if hasattr(control, 'PowerSystemResource_id') and control.PowerSystemResource_id:
        psr_id = control.PowerSystemResource_id
        # 查询 PowerSystemResource
        psr = session.query(model.PowerSystemResource).filter(
            model.PowerSystemResource.id == psr_id
        ).first()
        
        if psr:
            print(f"{indent}关联的管理对象 (PowerSystemResource):")
            print(f"{indent}  ├─ 名称: {psr.name}")
            print(f"{indent}  ├─ mRID: {psr.mRID}")
            print(f"{indent}  ├─ 类型: {psr.type_}")
            if hasattr(psr, 'description') and psr.description:
                print(f"{indent}  ├─ 描述: {psr.description}")
            
            # 如果是设备，显示更多信息
            if hasattr(model, 'Equipment'):
                equipment = session.query(model.Equipment).filter(
                    model.Equipment.id == psr_id
                ).first()
                if equipment:
                    print(f"{indent}  └─ 设备类型: {equipment.type_}")
            
            # 如果是 ConductingEquipment，显示更多信息
            if hasattr(model, 'ConductingEquipment'):
                cond_eq = session.query(model.ConductingEquipment).filter(
                    model.ConductingEquipment.id == psr_id
                ).first()
                if cond_eq:
                    print(f"{indent}  └─ 导电设备类型: {cond_eq.type_}")
        else:
            print(f"{indent}关联的管理对象: 未找到 (ID: {psr_id})")
    else:
        print(f"{indent}关联的管理对象: 无")

def main():
    # 数据库路径
    db_path = project_root / 'cimpyorm' / 'cimpyorm' / 'res' / 'datasets' / 'FullGrid' / 'out.db'
    
    print("=" * 100)
    print("Control 对象及其关联的管理对象")
    print("=" * 100)
    
    # 加载数据库
    session, model = load(str(db_path))
    
    # 查找所有 Control 子类的实例（按继承层次从子类到父类）
    control_classes = [
        ('AccumulatorReset', model.AccumulatorReset),
        ('RaiseLowerCommand', model.RaiseLowerCommand),
        ('SetPoint', model.SetPoint),
        ('Command', model.Command),
        ('AnalogControl', model.AnalogControl),
        ('Control', model.Control),
    ]
    
    # 用于跟踪已显示的 Control ID，避免重复
    displayed_ids = set()
    total_count = 0
    
    for class_name, class_obj in control_classes:
        if not hasattr(model, class_name):
            continue
            
        controls = session.query(class_obj).all()
        if controls:
            # 过滤掉已经显示过的对象
            new_controls = [c for c in controls if c.id not in displayed_ids]
            if new_controls:
                total_count += len(new_controls)
                print(f"\n{'=' * 100}")
                print(f"{class_name} ({len(new_controls)} 个实例)")
                print(f"{'=' * 100}")
                
                for i, control in enumerate(new_controls, 1):
                    displayed_ids.add(control.id)
                    print(f"\n[{i}] {class_name} 实例:")
                    print("-" * 100)
                    print_control_info(control, session, model, "  ")
                
                # 如果是 AccumulatorReset，显示关联的 AccumulatorValue
                if class_name == 'AccumulatorReset' and hasattr(control, 'AccumulatorValue_id'):
                    acc_value_id = control.AccumulatorValue_id
                    if acc_value_id and hasattr(model, 'AccumulatorValue'):
                        acc_value = session.query(model.AccumulatorValue).filter(
                            model.AccumulatorValue.id == acc_value_id
                        ).first()
                        if acc_value:
                            print("  关联的 AccumulatorValue:")
                            print(f"    名称: {acc_value.name}")
                            print(f"    值: {acc_value.value}")
                
                # 如果是 Command，显示关联的 DiscreteValue
                if class_name == 'Command' and hasattr(control, 'DiscreteValue_id'):
                    disc_value_id = control.DiscreteValue_id
                    if disc_value_id and hasattr(model, 'DiscreteValue'):
                        disc_value = session.query(model.DiscreteValue).filter(
                            model.DiscreteValue.id == disc_value_id
                        ).first()
                        if disc_value:
                            print("  关联的 DiscreteValue:")
                            print(f"    名称: {disc_value.name}")
                            print(f"    值: {disc_value.value}")
                
                # 如果是 RaiseLowerCommand，显示关联的 AnalogValue
                if class_name == 'RaiseLowerCommand' and hasattr(control, 'AnalogValue_id'):
                    analog_value_id = control.AnalogValue_id
                    if analog_value_id and hasattr(model, 'AnalogValue'):
                        analog_value = session.query(model.AnalogValue).filter(
                            model.AnalogValue.id == analog_value_id
                        ).first()
                        if analog_value:
                            print("  关联的 AnalogValue:")
                            print(f"    名称: {analog_value.name}")
                            print(f"    值: {analog_value.value}")
                            if hasattr(control, 'maxValue'):
                                print(f"    最大值: {control.maxValue}")
                            if hasattr(control, 'minValue'):
                                print(f"    最小值: {control.minValue}")
    
    print(f"\n{'=' * 100}")
    print(f"总计: {total_count} 个 Control 对象")
    print(f"{'=' * 100}")
    
    # 统计信息
    print("\n统计信息:")
    print("-" * 100)
    
    # 统计每个 Control 类型关联的 PowerSystemResource 类型（只统计已显示的对象）
    psr_type_count = {}
    
    for class_name, class_obj in control_classes:
        if not hasattr(model, class_name):
            continue
            
        controls = session.query(class_obj).all()
        # 只统计已显示的对象
        for control in controls:
            if control.id in displayed_ids:
                if hasattr(control, 'PowerSystemResource_id') and control.PowerSystemResource_id:
                    psr_id = control.PowerSystemResource_id
                    psr = session.query(model.PowerSystemResource).filter(
                        model.PowerSystemResource.id == psr_id
                    ).first()
                    if psr:
                        psr_type = psr.type_
                        if psr_type not in psr_type_count:
                            psr_type_count[psr_type] = {}
                        if class_name not in psr_type_count[psr_type]:
                            psr_type_count[psr_type][class_name] = 0
                        psr_type_count[psr_type][class_name] += 1
    
    print("\nControl 类型关联的 PowerSystemResource 类型统计:")
    for psr_type, control_types in sorted(psr_type_count.items()):
        print(f"\n  {psr_type}:")
        for ctrl_type, count in sorted(control_types.items()):
            print(f"    - {ctrl_type}: {count} 个")

if __name__ == '__main__':
    main()

