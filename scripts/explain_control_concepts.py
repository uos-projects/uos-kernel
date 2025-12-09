#!/usr/bin/env python3
"""
解释 Control 对象的三个关键概念：
1. 控制类型（controlType）
2. 管理对象（PowerSystemResource）
3. 关联值（MeasurementValue）
"""
import sys
from pathlib import Path

# 添加 cimpyorm 到 sys.path
project_root = Path(__file__).parent.parent
sys.path.insert(0, str(project_root / 'cimpyorm'))

from cimpyorm import load

def main():
    # 数据库路径
    db_path = project_root / 'cimpyorm' / 'cimpyorm' / 'res' / 'datasets' / 'FullGrid' / 'out.db'
    
    print("=" * 100)
    print("Control 对象关键概念解释")
    print("=" * 100)
    
    # 加载数据库
    session, model = load(str(db_path))
    
    print("\n" + "=" * 100)
    print("1. 控制类型（controlType）")
    print("=" * 100)
    print("""
控制类型（controlType）是一个字符串，用于指定 Control 对象的控制类型。
它描述了该控制命令的具体用途和操作类型。

常见控制类型包括：
  - ThreePhaseActivePower: 三相有功功率控制
  - SwitchPosition: 开关位置控制（开/关）
  - VoltageSetPoint: 电压设定点控制
  - FrequencySetPoint: 频率设定点控制
  - TieLineFlow: 联络线潮流控制
  - GeneratorVoltageSetPoint: 发电机电压设定点
  - BreakerOn/Off: 断路器开关控制

控制类型的作用：
  - 标识控制命令的语义含义
  - 帮助系统理解如何执行该控制
  - 用于控制命令的分类和路由
    """)
    
    # 展示实际的控制类型
    print("\nFullGrid 数据集中的实际控制类型：")
    print("-" * 100)
    controls = session.query(model.Control).all()
    control_types = {}
    for ctrl in controls:
        ctrl_type = ctrl.controlType if hasattr(ctrl, 'controlType') else 'N/A'
        ctrl_name = ctrl.name
        if ctrl_type not in control_types:
            control_types[ctrl_type] = []
        control_types[ctrl_type].append(ctrl_name)
    
    for ctrl_type, names in sorted(control_types.items()):
        print(f"  {ctrl_type}:")
        for name in names:
            print(f"    - {name}")
    
    print("\n" + "=" * 100)
    print("2. 管理对象（PowerSystemResource）")
    print("=" * 100)
    print("""
管理对象（PowerSystemResource）是 Control 对象所控制的电力系统资源。

关系说明：
  - Control 通过 PowerSystemResource.Controls 属性关联到 PowerSystemResource
  - 一个 PowerSystemResource 可以有多个 Control（一对多关系）
  - 一个 Control 只能关联一个 PowerSystemResource（多对一关系）

管理对象的类型可以是：
  - Equipment（设备）: 如 Breaker（断路器）、SynchronousMachine（同步电机）
  - EquipmentContainer（设备容器）: 如 Substation（变电站）、VoltageLevel（电压等级）
  - ControlArea（控制区域）: 用于自动发电控制
  - 其他 PowerSystemResource 的子类

作用：
  - 标识控制命令的目标对象
  - 建立控制命令与电力系统资源的关联
  - 支持控制命令的执行和状态跟踪
    """)
    
    # 展示实际的管理对象
    print("\nFullGrid 数据集中的实际管理对象：")
    print("-" * 100)
    for ctrl in controls:
        if hasattr(ctrl, 'PowerSystemResource_id') and ctrl.PowerSystemResource_id:
            psr_id = ctrl.PowerSystemResource_id
            psr = session.query(model.PowerSystemResource).filter(
                model.PowerSystemResource.id == psr_id
            ).first()
            if psr:
                print(f"\n  Control: {ctrl.name}")
                print(f"    管理对象: {psr.name}")
                print(f"    类型: {psr.type_}")
                print(f"    描述: {psr.description if hasattr(psr, 'description') and psr.description else 'N/A'}")
    
    print("\n" + "=" * 100)
    print("3. 关联值（MeasurementValue）")
    print("=" * 100)
    print("""
关联值（MeasurementValue）是 Control 对象关联的测量值或控制值。

不同类型的 Control 关联不同类型的 MeasurementValue：

  a) AccumulatorReset（累加器复位）
     - 关联: AccumulatorValue（累加器值）
     - 作用: 复位累加器的计数值
     - 示例: 复位电能表的累计电量

  b) Command（命令）
     - 关联: DiscreteValue（离散值）
     - 作用: 设置离散状态值（如开关位置 0/1）
     - 示例: 设置断路器位置为"开"或"关"

  c) RaiseLowerCommand（升降命令）
     - 关联: AnalogValue（模拟值）
     - 作用: 通过脉冲增加或减少设定值
     - 示例: 通过脉冲调整发电机功率设定点

  d) SetPoint（设定点）
     - 关联: AnalogValue（模拟值）
     - 作用: 直接设置模拟量的目标值
     - 示例: 设置发电机功率目标值

关联值的作用：
  - 存储控制命令的具体数值
  - 记录控制命令的执行结果
  - 支持控制命令的状态查询和历史追溯
    """)
    
    # 展示实际的关联值
    print("\nFullGrid 数据集中的实际关联值：")
    print("-" * 100)
    
    # AccumulatorReset -> AccumulatorValue
    acc_resets = session.query(model.AccumulatorReset).all()
    for acc_reset in acc_resets:
        if hasattr(acc_reset, 'AccumulatorValue_id') and acc_reset.AccumulatorValue_id:
            acc_value = session.query(model.AccumulatorValue).filter(
                model.AccumulatorValue.id == acc_reset.AccumulatorValue_id
            ).first()
            if acc_value:
                print(f"\n  Control: {acc_reset.name} (AccumulatorReset)")
                print(f"    关联值类型: AccumulatorValue")
                print(f"    关联值名称: {acc_value.name}")
                print(f"    值: {acc_value.value}")
                print(f"    作用: 复位累加器 {acc_value.Accumulator.name if hasattr(acc_value, 'Accumulator') else 'N/A'} 的计数值")
    
    # Command -> DiscreteValue
    commands = session.query(model.Command).all()
    for cmd in commands:
        if hasattr(cmd, 'DiscreteValue_id') and cmd.DiscreteValue_id:
            disc_value = session.query(model.DiscreteValue).filter(
                model.DiscreteValue.id == cmd.DiscreteValue_id
            ).first()
            if disc_value:
                print(f"\n  Control: {cmd.name} (Command)")
                print(f"    关联值类型: DiscreteValue")
                print(f"    关联值名称: {disc_value.name}")
                print(f"    值: {disc_value.value}")
                print(f"    作用: 设置离散状态值（如开关位置）")
    
    # RaiseLowerCommand -> AnalogValue
    rl_commands = session.query(model.RaiseLowerCommand).all()
    for rl_cmd in rl_commands:
        if hasattr(rl_cmd, 'AnalogValue_id') and rl_cmd.AnalogValue_id:
            analog_value = session.query(model.AnalogValue).filter(
                model.AnalogValue.id == rl_cmd.AnalogValue_id
            ).first()
            if analog_value:
                print(f"\n  Control: {rl_cmd.name} (RaiseLowerCommand)")
                print(f"    关联值类型: AnalogValue")
                print(f"    关联值名称: {analog_value.name}")
                print(f"    值: {analog_value.value}")
                if hasattr(rl_cmd, 'maxValue'):
                    print(f"    最大值: {rl_cmd.maxValue}")
                if hasattr(rl_cmd, 'minValue'):
                    print(f"    最小值: {rl_cmd.minValue}")
                print(f"    作用: 通过脉冲增加或减少设定值")
    
    print("\n" + "=" * 100)
    print("总结")
    print("=" * 100)
    print("""
Control 对象的三个关键概念构成了一个完整的控制命令模型：

1. 控制类型（controlType）: "做什么" - 定义控制命令的语义
2. 管理对象（PowerSystemResource）: "对谁做" - 定义控制命令的目标
3. 关联值（MeasurementValue）: "怎么做" - 定义控制命令的具体数值

这三个概念共同描述了：
  - 控制命令的意图（控制类型）
  - 控制命令的目标（管理对象）
  - 控制命令的参数（关联值）

例如，一个完整的控制命令可能是：
  - 控制类型: "ThreePhaseActivePower"（三相有功功率控制）
  - 管理对象: "BE-G4"（同步电机）
  - 关联值: "100.0 MW"（功率设定值）
  
这表示：对同步电机 BE-G4 执行三相有功功率控制，设定值为 100.0 MW。
    """)

if __name__ == '__main__':
    main()

