#!/usr/bin/env python3
"""
查找 STR_MEAS_1 实体
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
    
    print("=" * 80)
    print("查找 STR_MEAS_1")
    print("=" * 80)
    
    # 加载数据库
    session, model = load(str(db_path))
    
    # 查找 StringMeasurement
    str_meas = session.query(model.StringMeasurement).filter(
        model.StringMeasurement.name == 'STR_MEAS_1'
    ).first()
    
    if str_meas:
        print(f"\n找到 StringMeasurement: {str_meas.name}")
        print(f"mRID: {str_meas.mRID}")
        print(f"描述: {str_meas.description}")
        print(f"测量类型: {str_meas.measurementType}")
        print(f"相位: {str_meas.phases}")
        print(f"单位倍数: {str_meas.unitMultiplier}")
        print(f"单位符号: {str_meas.unitSymbol}")
        
        # 查找关联的 PowerSystemResource
        if hasattr(str_meas, 'PowerSystemResource') and str_meas.PowerSystemResource:
            psr = str_meas.PowerSystemResource
            print(f"\n关联的 PowerSystemResource:")
            print(f"  名称: {psr.name}")
            print(f"  mRID: {psr.mRID}")
            print(f"  类型: {type(psr).__name__}")
        
        # 查找关联的 Terminal
        if hasattr(str_meas, 'Terminal') and str_meas.Terminal:
            terminal = str_meas.Terminal
            print(f"\n关联的 Terminal:")
            print(f"  名称: {terminal.name}")
            print(f"  mRID: {terminal.mRID}")
        
        # 查找关联的 StringMeasurementValue
        if hasattr(model, 'StringMeasurementValue'):
            str_values = session.query(model.StringMeasurementValue).filter(
                model.StringMeasurementValue.StringMeasurement_id == str_meas.id
            ).all()
            
            if str_values:
                print(f"\n关联的 StringMeasurementValue:")
                for val in str_values:
                    print(f"  名称: {val.name}")
                    print(f"  mRID: {val.mRID}")
                    print(f"  值: {val.value}")
                    print(f"  时间戳: {val.timeStamp}")
                    print(f"  传感器精度: {val.sensorAccuracy}")
                    
                    # 查找 MeasurementValueSource
                    if hasattr(val, 'MeasurementValueSource') and val.MeasurementValueSource:
                        source = val.MeasurementValueSource
                        print(f"  数据源: {source.name}")
        
        # 打印完整对象信息
        print("\n" + "=" * 80)
        print("完整对象属性:")
        print("=" * 80)
        for attr in dir(str_meas):
            if not attr.startswith('_') and not callable(getattr(str_meas, attr)):
                try:
                    value = getattr(str_meas, attr)
                    if value is not None:
                        print(f"  {attr}: {value}")
                except:
                    pass
    else:
        print("\n未找到 STR_MEAS_1")
        
        # 列出所有 StringMeasurement
        all_str_meas = session.query(model.StringMeasurement).all()
        print(f"\n数据集中的所有 StringMeasurement ({len(all_str_meas)} 个):")
        for sm in all_str_meas:
            print(f"  - {sm.name} (mRID: {sm.mRID})")

if __name__ == '__main__':
    main()


