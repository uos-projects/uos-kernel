# CIMPyORM FullGrid 数据集分析总结

**日期**: 2025-12-09  
**数据集**: FullGrid (cimpyorm/cimpyorm/res/datasets/FullGrid)

## 目录

1. [CIM Schema 中的 Control 关联关系](#1-cim-schema-中的-control-关联关系)
2. [数据集中的 Control 和 Measurement 实体](#2-数据集中的-control-和-measurement-实体)
3. [启动 CIMPyORM 并加载 FullGrid 数据集](#3-启动-cimpyorm-并加载-fullgrid-数据集)
4. [类层次结构分析](#4-类层次结构分析)
5. [Control 对象及其管理对象](#5-control-对象及其管理对象)
6. [Control 对象关键概念解释](#6-control-对象关键概念解释)

---

## 1. CIM Schema 中的 Control 关联关系

### 1.1 关联类型

在 `cimpyorm/cimpyorm/res/schemata/CIM16` 目录下的 RDF 文件中，**`PowerSystemResource`** 类型关联到 `Control`。

### 1.2 关联关系详情

#### PowerSystemResource.Controls
- **文件位置**: `EquipmentProfileCoreOperationRDFSAugmented-v2_4_15-4Jul2016.rdf`
- **Domain**: `PowerSystemResource`
- **Range**: `Control`
- **基数**: 0..n（一个资源可以有多个控制）
- **说明**: "Regulating device governed by this control output."

#### Control.PowerSystemResource（反向关联）
- **Domain**: `Control`
- **Range**: `PowerSystemResource`
- **基数**: 0..1（一个控制可以关联一个资源）
- **说明**: "The controller outputs used to actually govern a regulating device."

### 1.3 关系总结

- **关联类型**: `PowerSystemResource` ↔ `Control`
- **关系**: 双向关联（互逆）
- **基数**: 
  - `PowerSystemResource` → `Control`: 0..n
  - `Control` → `PowerSystemResource`: 0..1

---

## 2. 数据集中的 Control 和 Measurement 实体

### 2.1 Control 相关实体

在 `cimpyorm/cimpyorm/res/datasets/FullGrid` 数据集中，包含以下 Control 子类的实例：

#### a) AccumulatorReset（累加器复位）
- **文件**: `20171002T0930Z_BE_EQ_4.xml`
- **示例**: ACC_RESET_1
- **控制类型**: ThreePhaseActivePower
- **关联对象**: BE-Line_2 (ACLineSegment)

#### b) Command（命令）
- **示例**: CMD_1
- **控制类型**: SwitchPosition
- **关联对象**: CIRCB-1230991526 (Breaker)

#### c) RaiseLowerCommand（升降命令）
- **示例**: CMD_RL_1
- **控制类型**: ThreePhaseActivePower
- **关联对象**: BE-G4 (SynchronousMachine)

### 2.2 Measurement 相关实体

在 FullGrid 数据集中，包含以下 Measurement 子类的实例：

#### a) Accumulator（累加器测量）
- **示例**: ACCUM_1
- **测量类型**: ThreePhasePower
- **关联对象**: BE-Line_2 (ACLineSegment)

#### b) Analog（模拟量测量）
- **示例**: ANA_1
- **测量类型**: ThreePhaseActivePower
- **关联对象**: 相关 PowerSystemResource

#### c) Discrete（离散量测量）
- **示例**: DIS_1
- **测量类型**: SwitchPosition
- **关联对象**: 相关 PowerSystemResource

#### d) StringMeasurement（字符串测量）
- **示例**: STR_MEAS_1
- **测量类型**: ThreePhasePower
- **关联对象**: SVC_1 (StaticVarCompensator)

### 2.3 Measurement 值相关实体

- **AccumulatorValue**: 累加器值
- **AnalogValue**: 模拟量值
- **DiscreteValue**: 离散量值
- **StringMeasurementValue**: 字符串测量值

### 2.4 Measurement 辅助实体

- **MeasurementValueSource**: 测量值源（如 ICCP）
- **MeasurementValueQuality**: 测量值质量（基于 IEC 61850 质量码）

---

## 3. 启动 CIMPyORM 并加载 FullGrid 数据集

### 3.1 环境准备

创建虚拟环境并安装依赖：

```bash
python3 -m venv venv
source venv/bin/activate
pip install -r cimpyorm/requirements.txt
```

### 3.2 加载数据集

使用脚本 `scripts/load_fullgrid.py` 加载 FullGrid 数据集：

```python
from cimpyorm import load, parse

# 数据集路径
fullgrid_path = 'cimpyorm/cimpyorm/res/datasets/FullGrid'
db_path = fullgrid_path / 'out.db'

# 加载或解析数据集
if db_path.exists():
    session, model = load(str(db_path))
else:
    session, model = parse(str(fullgrid_path))
```

### 3.3 数据集统计

- **总对象数**: 815 个
- **生成类数**: 397 个 CIM 类
- **数据库**: SQLite (`out.db`)

### 3.4 数据集中的实体统计

- **Terminal**: 144 个
- **ConductingEquipment**: 94 个
- **Equipment**: 116 个
- **PowerSystemResource**: 188 个

---

## 4. 类层次结构分析

### 4.1 总体统计

- **总类数**: 397 个
- **根类数**: 28 个
- **最大根类**: `cim_IdentifiedObject` 有 366 个子类

### 4.2 主要根类

1. **cim_IdentifiedObject** - 基础标识对象（366 个子类）
2. **cim_CurveData** - 曲线数据
3. **cim_Quality61850** - IEC 61850 质量码（1 个子类）
4. **cim_TapChangerTablePoint** - 分接头表点（2 个子类）

### 4.3 重要类的层次结构

#### Measurement（测量）层次
```
cim_Measurement
├── cim_Accumulator (累加器)
├── cim_Analog (模拟量)
├── cim_Discrete (离散量)
└── cim_StringMeasurement (字符串测量)
```

#### Control（控制）层次
```
cim_Control
├── cim_AccumulatorReset (累加器复位)
├── cim_AnalogControl (模拟控制)
│   ├── cim_RaiseLowerCommand (升降命令)
│   └── cim_SetPoint (设定点)
└── cim_Command (命令)
```

#### PowerSystemResource（电力系统资源）层次
```
cim_PowerSystemResource
├── cim_ConnectivityNodeContainer
│   ├── cim_EquipmentContainer
│   │   ├── cim_Bay (间隔)
│   │   ├── cim_Substation (变电站)
│   │   └── cim_VoltageLevel (电压等级)
│   └── cim_EquivalentNetwork
├── cim_ControlArea (控制区域)
├── cim_Equipment (设备)
│   ├── cim_ConductingEquipment (导电设备)
│   ├── cim_DCConductingEquipment (直流导电设备)
│   └── cim_GeneratingUnit (发电机组)
└── cim_RegulatingControl (调节控制)
```

### 4.4 完整类层次结构文件

完整的类层次结构已保存到：`scripts/fullgrid_class_hierarchy.txt`（680 行）

---

## 5. Control 对象及其管理对象

### 5.1 FullGrid 数据集中的 Control 对象

总计：**4 个 Control 对象**（去重后）

#### 1. AccumulatorReset（累加器复位）- 1 个
- **名称**: ACC_RESET_1
- **控制类型**: ThreePhaseActivePower（三相有功功率）
- **管理对象**: BE-Line_2 (ACLineSegment - 交流线路段)
- **关联值**: AccumulatorValue (ACC_VALUE_1, 值: 100)

#### 2. RaiseLowerCommand（升降命令）- 1 个
- **名称**: CMD_RL_1
- **控制类型**: ThreePhaseActivePower（三相有功功率）
- **管理对象**: BE-G4 (SynchronousMachine - 同步电机)
- **关联值**: AnalogValue (ANA_VALUE_1, 值: 100.0, 范围: 39.99-99.99)

#### 3. SetPoint（设定点）- 1 个
- **名称**: SET_PNT_1
- **控制类型**: ThreePhaseActivePower（三相有功功率）
- **管理对象**: BE-G1 (SynchronousMachine - 同步电机)

#### 4. Command（命令）- 1 个
- **名称**: CMD_1
- **控制类型**: SwitchPosition（开关位置）
- **管理对象**: CIRCB-1230991526 (Breaker - 断路器)
- **关联值**: DiscreteValue (DIS_VALUE_1, 值: 0)

### 5.2 管理对象类型分布

1. **cim_SynchronousMachine（同步电机）**: 2 个 Control
   - RaiseLowerCommand: 1 个
   - SetPoint: 1 个

2. **cim_ACLineSegment（交流线路段）**: 1 个 Control
   - AccumulatorReset: 1 个

3. **cim_Breaker（断路器）**: 1 个 Control
   - Command: 1 个

---

## 6. Control 对象关键概念解释

### 6.1 控制类型（controlType）

**含义**: 控制类型是一个字符串属性，用于指定 Control 对象的控制类型。它描述了该控制命令的具体用途和操作类型。

**作用**:
- 标识控制命令的语义含义
- 帮助系统理解如何执行该控制
- 用于控制命令的分类和路由

**常见控制类型**:
- `ThreePhaseActivePower`: 三相有功功率控制
- `SwitchPosition`: 开关位置控制（开/关）
- `VoltageSetPoint`: 电压设定点控制
- `FrequencySetPoint`: 频率设定点控制
- `TieLineFlow`: 联络线潮流控制
- `GeneratorVoltageSetPoint`: 发电机电压设定点
- `BreakerOn/Off`: 断路器开关控制

**FullGrid 数据集中的实际控制类型**:
- `SwitchPosition`: CMD_1
- `ThreePhaseActivePower`: ACC_RESET_1, CMD_RL_1, SET_PNT_1

### 6.2 管理对象（PowerSystemResource）

**含义**: 管理对象是 Control 对象所控制的电力系统资源，即控制命令的目标对象。

**关系说明**:
- Control 通过 `PowerSystemResource.Controls` 属性关联到 PowerSystemResource
- 一个 PowerSystemResource 可以有多个 Control（一对多关系）
- 一个 Control 只能关联一个 PowerSystemResource（多对一关系）

**管理对象的类型可以是**:
- **Equipment（设备）**: 如 Breaker（断路器）、SynchronousMachine（同步电机）
- **EquipmentContainer（设备容器）**: 如 Substation（变电站）、VoltageLevel（电压等级）
- **ControlArea（控制区域）**: 用于自动发电控制
- 其他 PowerSystemResource 的子类

**作用**:
- 标识控制命令的目标对象
- 建立控制命令与电力系统资源的关联
- 支持控制命令的执行和状态跟踪

**FullGrid 数据集中的实际管理对象**:
- BE-Line_2 (ACLineSegment) - 关联 ACC_RESET_1
- CIRCB-1230991526 (Breaker) - 关联 CMD_1
- BE-G4 (SynchronousMachine) - 关联 CMD_RL_1
- BE-G1 (SynchronousMachine) - 关联 SET_PNT_1

### 6.3 关联值（MeasurementValue）

**含义**: 关联值是 Control 对象关联的测量值或控制值，用于存储控制命令的具体数值。

**不同类型的 Control 关联不同类型的 MeasurementValue**:

| Control 类型 | 关联值类型 | 作用 | 示例 |
|-------------|-----------|------|------|
| **AccumulatorReset** | AccumulatorValue | 复位累加器的计数值 | 复位电能表的累计电量 |
| **Command** | DiscreteValue | 设置离散状态值（0/1） | 设置断路器位置为"开"或"关" |
| **RaiseLowerCommand** | AnalogValue | 通过脉冲增加或减少设定值 | 通过脉冲调整发电机功率设定点 |
| **SetPoint** | AnalogValue | 直接设置模拟量的目标值 | 设置发电机功率目标值 |

**作用**:
- 存储控制命令的具体数值
- 记录控制命令的执行结果
- 支持控制命令的状态查询和历史追溯

**FullGrid 数据集中的实际关联值**:
- ACC_RESET_1 → AccumulatorValue (ACC_VALUE_1, 值: 100)
- CMD_1 → DiscreteValue (DIS_VALUE_1, 值: 0)
- CMD_RL_1 → AnalogValue (ANA_VALUE_1, 值: 100.0, 范围: 39.99-99.99)

### 6.4 完整示例

一个完整的控制命令示例：

```
Control: CMD_RL_1 (RaiseLowerCommand)
├─ 控制类型: "ThreePhaseActivePower"  ← "做什么"（三相有功功率控制）
├─ 管理对象: "BE-G4" (SynchronousMachine)  ← "对谁做"（同步电机）
└─ 关联值: AnalogValue (100.0 MW, 范围: 39.99-99.99)  ← "怎么做"（功率设定值）
```

**含义**: 对同步电机 BE-G4 执行三相有功功率控制，通过脉冲将功率设定值调整到 100.0 MW。

### 6.5 总结

Control 对象的三个关键概念构成了一个完整的控制命令模型：

1. **控制类型（controlType）**: "做什么" - 定义控制命令的语义
2. **管理对象（PowerSystemResource）**: "对谁做" - 定义控制命令的目标
3. **关联值（MeasurementValue）**: "怎么做" - 定义控制命令的具体数值

这三个概念共同描述了：
- 控制命令的意图（控制类型）
- 控制命令的目标（管理对象）
- 控制命令的参数（关联值）

---

## 7. 相关脚本文件

本次分析过程中创建的脚本文件：

1. **scripts/load_fullgrid.py**: 加载 FullGrid 数据集
2. **scripts/find_str_meas.py**: 查找 STR_MEAS_1 实体
3. **scripts/show_class_hierarchy.py**: 展示类层次结构
4. **scripts/show_control_objects.py**: 展示 Control 对象及其管理对象
5. **scripts/explain_control_concepts.py**: 解释 Control 对象关键概念

## 8. 输出文件

1. **scripts/fullgrid_class_hierarchy.txt**: 完整的类层次结构（680 行）
2. **cimpyorm/cimpyorm/res/datasets/FullGrid/out.db**: SQLite 数据库文件

---

## 9. 关键发现

1. **Control 关联关系**: PowerSystemResource 是唯一关联到 Control 的类型
2. **数据集实体**: FullGrid 数据集包含 4 个 Control 对象和 4 个 Measurement 对象
3. **类层次结构**: 共 397 个类，其中 cim_IdentifiedObject 有 366 个子类
4. **控制模式**: Control 对象通过控制类型、管理对象和关联值三个维度完整描述控制命令

---

## 10. 后续工作建议

1. 深入分析 Measurement 对象及其关联关系
2. 研究 PowerSystemResource 的完整层次结构
3. 分析 Control 和 Measurement 的协同工作机制
4. 探索 CIM 模型在电力系统中的应用场景

---

**文档生成时间**: 2025-12-09  
**数据集版本**: CIM16  
**分析工具**: CIMPyORM

