#!/usr/bin/env python3
"""
查找 FullGrid 数据集中所有 Control 对象及其所有关联对象
"""
import sys
from pathlib import Path

# 添加 cimpyorm 到 sys.path
project_root = Path(__file__).parent.parent
sys.path.insert(0, str(project_root / 'cimpyorm'))

from cimpyorm import load

def get_object_info(obj, session, model):
    """获取对象的基本信息"""
    info = {
        'name': getattr(obj, 'name', 'N/A'),
        'mRID': getattr(obj, 'mRID', 'N/A'),
        'type': getattr(obj, 'type_', 'N/A'),
    }
    if hasattr(obj, 'description') and obj.description:
        info['description'] = obj.description
    return info

def find_all_control_relations():
    """查找所有 Control 对象及其关联对象"""
    # 数据库路径
    db_path = project_root / 'cimpyorm' / 'cimpyorm' / 'res' / 'datasets' / 'FullGrid' / 'out.db'
    
    print("=" * 100)
    print("FullGrid 数据集中所有 Control 对象及其关联对象")
    print("=" * 100)
    
    # 加载数据库
    session, model = load(str(db_path))
    
    # 查找所有 Control 子类
    control_classes = [
        ('AccumulatorReset', model.AccumulatorReset),
        ('RaiseLowerCommand', model.RaiseLowerCommand),
        ('SetPoint', model.SetPoint),
        ('Command', model.Command),
        ('AnalogControl', model.AnalogControl),
        ('Control', model.Control),
    ]
    
    all_controls = []
    displayed_ids = set()
    
    # 收集所有 Control 对象
    for class_name, class_obj in control_classes:
        if not hasattr(model, class_name):
            continue
        
        controls = session.query(class_obj).all()
        for control in controls:
            if control.id not in displayed_ids:
                displayed_ids.add(control.id)
                all_controls.append((class_name, control))
    
    print(f"\n找到 {len(all_controls)} 个 Control 对象\n")
    
    # 遍历每个 Control 对象，查找所有关联对象
    for idx, (class_name, control) in enumerate(all_controls, 1):
        print("=" * 100)
        print(f"[{idx}] {class_name}: {control.name}")
        print("=" * 100)
        
        # Control 基本信息
        print("\n【Control 基本信息】")
        print(f"  名称: {control.name}")
        print(f"  mRID: {control.mRID}")
        print(f"  类型: {class_name}")
        
        if hasattr(control, 'description') and control.description:
            print(f"  描述: {control.description}")
        if hasattr(control, 'controlType') and control.controlType:
            print(f"  控制类型: {control.controlType}")
        if hasattr(control, 'timeStamp') and control.timeStamp:
            print(f"  时间戳: {control.timeStamp}")
        if hasattr(control, 'operationInProgress') and control.operationInProgress is not None:
            print(f"  操作进行中: {control.operationInProgress}")
        
        # 1. 关联的 PowerSystemResource（管理对象）
        print("\n【关联的管理对象 (PowerSystemResource)】")
        if hasattr(control, 'PowerSystemResource_id') and control.PowerSystemResource_id:
            psr_id = control.PowerSystemResource_id
            psr = session.query(model.PowerSystemResource).filter(
                model.PowerSystemResource.id == psr_id
            ).first()
            
            if psr:
                psr_info = get_object_info(psr, session, model)
                print(f"  ✓ 找到关联对象:")
                print(f"    名称: {psr_info['name']}")
                print(f"    mRID: {psr_info['mRID']}")
                print(f"    类型: {psr_info['type']}")
                if 'description' in psr_info:
                    print(f"    描述: {psr_info['description']}")
                
                # 进一步查找具体类型
                # Equipment
                if hasattr(model, 'Equipment'):
                    equipment = session.query(model.Equipment).filter(
                        model.Equipment.id == psr_id
                    ).first()
                    if equipment:
                        print(f"    设备类型: {equipment.type_}")
                
                # ConductingEquipment
                if hasattr(model, 'ConductingEquipment'):
                    cond_eq = session.query(model.ConductingEquipment).filter(
                        model.ConductingEquipment.id == psr_id
                    ).first()
                    if cond_eq:
                        print(f"    导电设备类型: {cond_eq.type_}")
                        # 如果是 Breaker
                        if hasattr(model, 'Breaker'):
                            breaker = session.query(model.Breaker).filter(
                                model.Breaker.id == psr_id
                            ).first()
                            if breaker:
                                print(f"    具体类型: Breaker (断路器)")
                        # 如果是 SynchronousMachine
                        if hasattr(model, 'SynchronousMachine'):
                            sync_machine = session.query(model.SynchronousMachine).filter(
                                model.SynchronousMachine.id == psr_id
                            ).first()
                            if sync_machine:
                                print(f"    具体类型: SynchronousMachine (同步电机)")
                        # 如果是 ACLineSegment
                        if hasattr(model, 'ACLineSegment'):
                            line_seg = session.query(model.ACLineSegment).filter(
                                model.ACLineSegment.id == psr_id
                            ).first()
                            if line_seg:
                                print(f"    具体类型: ACLineSegment (交流线路段)")
            else:
                print(f"  ✗ 未找到 (ID: {psr_id})")
        else:
            print("  ✗ 无关联对象")
        
        # 2. 关联的 MeasurementValue（根据 Control 类型不同）
        print("\n【关联的测量值/控制值 (MeasurementValue)】")
        
        # AccumulatorReset -> AccumulatorValue
        if class_name == 'AccumulatorReset' and hasattr(control, 'AccumulatorValue_id'):
            acc_value_id = control.AccumulatorValue_id
            if acc_value_id and hasattr(model, 'AccumulatorValue'):
                acc_value = session.query(model.AccumulatorValue).filter(
                    model.AccumulatorValue.id == acc_value_id
                ).first()
                if acc_value:
                    print(f"  ✓ AccumulatorValue:")
                    print(f"    名称: {acc_value.name}")
                    print(f"    mRID: {acc_value.mRID}")
                    if hasattr(acc_value, 'value') and acc_value.value is not None:
                        print(f"    值: {acc_value.value}")
                    if hasattr(acc_value, 'timeStamp') and acc_value.timeStamp:
                        print(f"    时间戳: {acc_value.timeStamp}")
        
        # Command -> DiscreteValue
        if class_name == 'Command' and hasattr(control, 'DiscreteValue_id'):
            disc_value_id = control.DiscreteValue_id
            if disc_value_id and hasattr(model, 'DiscreteValue'):
                disc_value = session.query(model.DiscreteValue).filter(
                    model.DiscreteValue.id == disc_value_id
                ).first()
                if disc_value:
                    print(f"  ✓ DiscreteValue:")
                    print(f"    名称: {disc_value.name}")
                    print(f"    mRID: {disc_value.mRID}")
                    if hasattr(disc_value, 'value') and disc_value.value is not None:
                        print(f"    值: {disc_value.value}")
                    if hasattr(disc_value, 'timeStamp') and disc_value.timeStamp:
                        print(f"    时间戳: {disc_value.timeStamp}")
        
        # RaiseLowerCommand -> AnalogValue
        if class_name == 'RaiseLowerCommand' and hasattr(control, 'AnalogValue_id'):
            analog_value_id = control.AnalogValue_id
            if analog_value_id and hasattr(model, 'AnalogValue'):
                analog_value = session.query(model.AnalogValue).filter(
                    model.AnalogValue.id == analog_value_id
                ).first()
                if analog_value:
                    print(f"  ✓ AnalogValue:")
                    print(f"    名称: {analog_value.name}")
                    print(f"    mRID: {analog_value.mRID}")
                    if hasattr(analog_value, 'value') and analog_value.value is not None:
                        print(f"    值: {analog_value.value}")
                    if hasattr(analog_value, 'timeStamp') and analog_value.timeStamp:
                        print(f"    时间戳: {analog_value.timeStamp}")
                    # RaiseLowerCommand 可能有范围限制
                    if hasattr(control, 'maxValue') and control.maxValue is not None:
                        print(f"    最大值: {control.maxValue}")
                    if hasattr(control, 'minValue') and control.minValue is not None:
                        print(f"    最小值: {control.minValue}")
        
        # SetPoint -> AnalogValue
        if class_name == 'SetPoint' and hasattr(control, 'AnalogValue_id'):
            analog_value_id = control.AnalogValue_id
            if analog_value_id and hasattr(model, 'AnalogValue'):
                analog_value = session.query(model.AnalogValue).filter(
                    model.AnalogValue.id == analog_value_id
                ).first()
                if analog_value:
                    print(f"  ✓ AnalogValue:")
                    print(f"    名称: {analog_value.name}")
                    print(f"    mRID: {analog_value.mRID}")
                    if hasattr(analog_value, 'value') and analog_value.value is not None:
                        print(f"    值: {analog_value.value}")
                    if hasattr(analog_value, 'timeStamp') and analog_value.timeStamp:
                        print(f"    时间戳: {analog_value.timeStamp}")
        
        # 如果没有找到关联值
        if not any([
            (class_name == 'AccumulatorReset' and hasattr(control, 'AccumulatorValue_id') and control.AccumulatorValue_id),
            (class_name == 'Command' and hasattr(control, 'DiscreteValue_id') and control.DiscreteValue_id),
            (class_name == 'RaiseLowerCommand' and hasattr(control, 'AnalogValue_id') and control.AnalogValue_id),
            (class_name == 'SetPoint' and hasattr(control, 'AnalogValue_id') and control.AnalogValue_id),
        ]):
            print("  ✗ 无关联值")
        
        # 3. 查找关联的 Measurement（如果有）
        print("\n【关联的测量 (Measurement)】")
        # 通过 PowerSystemResource 查找关联的 Measurement
        if hasattr(control, 'PowerSystemResource_id') and control.PowerSystemResource_id:
            psr_id = control.PowerSystemResource_id
            # 查找关联到同一个 PowerSystemResource 的 Measurement
            measurements = []
            
            # Accumulator
            if hasattr(model, 'Accumulator'):
                accums = session.query(model.Accumulator).filter(
                    model.Accumulator.PowerSystemResource_id == psr_id
                ).all()
                measurements.extend([('Accumulator', m) for m in accums])
            
            # Analog
            if hasattr(model, 'Analog'):
                analogs = session.query(model.Analog).filter(
                    model.Analog.PowerSystemResource_id == psr_id
                ).all()
                measurements.extend([('Analog', m) for m in analogs])
            
            # Discrete
            if hasattr(model, 'Discrete'):
                discretes = session.query(model.Discrete).filter(
                    model.Discrete.PowerSystemResource_id == psr_id
                ).all()
                measurements.extend([('Discrete', m) for m in discretes])
            
            # StringMeasurement
            if hasattr(model, 'StringMeasurement'):
                str_meas = session.query(model.StringMeasurement).filter(
                    model.StringMeasurement.PowerSystemResource_id == psr_id
                ).all()
                measurements.extend([('StringMeasurement', m) for m in str_meas])
            
            if measurements:
                print(f"  ✓ 找到 {len(measurements)} 个关联的 Measurement:")
                for meas_type, meas in measurements:
                    meas_info = get_object_info(meas, session, model)
                    print(f"    - {meas_type}: {meas_info['name']} (mRID: {meas_info['mRID']})")
                    if hasattr(meas, 'measurementType') and meas.measurementType:
                        print(f"      测量类型: {meas.measurementType}")
            else:
                print("  ✗ 无关联的 Measurement")
        else:
            print("  ✗ 无关联的 Measurement（无管理对象）")
        
        print()
    
    # 统计信息
    print("=" * 100)
    print("统计信息")
    print("=" * 100)
    
    # 按类型统计
    type_count = {}
    for class_name, control in all_controls:
        type_count[class_name] = type_count.get(class_name, 0) + 1
    
    print("\nControl 类型分布:")
    for ctrl_type, count in sorted(type_count.items()):
        print(f"  {ctrl_type}: {count} 个")
    
    # 统计关联的 PowerSystemResource 类型
    psr_type_count = {}
    for class_name, control in all_controls:
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
    
    print("\n关联的 PowerSystemResource 类型分布:")
    for psr_type, control_types in sorted(psr_type_count.items()):
        print(f"\n  {psr_type}:")
        for ctrl_type, count in sorted(control_types.items()):
            print(f"    - {ctrl_type}: {count} 个")
    
    print("\n" + "=" * 100)

if __name__ == '__main__':
    find_all_control_relations()

