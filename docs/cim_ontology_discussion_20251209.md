# CIM 本体讨论总结 - 2025年12月9日

## 对话概览

本次对话主要围绕 CIM (Common Information Model) 本体中的动作类（Action Classes）展开，深入探讨了 OWL 属性定义、语义建模、以及设备与操作之间的关系。

---

## 主要讨论主题

### 1. CIM 本体中的动作和操作类

#### 问题
- 查找 `TheCimOntology.owl` 中包含 Action 和 Event 类型属性的类
- 列出所有动作或操作类

#### 发现

**动作类（Action Classes）总计 32 个：**

**开关动作类（Switching Actions）- 16 个：**
1. ClampAction - 夹钳操作
2. ClearanceAction - 许可操作
3. ControlAction - 控制执行
4. CutAction - 切口操作
5. EnergyConsumerAction - 能源消费者操作
6. EnergySourceAction - 能源源操作
7. GenericAction - 任意开关步骤
8. GroundAction - 接地操作
9. JumperAction - 跳线操作
10. MeasurementAction - 测量操作
11. ProtectiveAction - 保护动作
12. ShuntCompensatorAction - 并联补偿器操作
13. SwitchAction - 开关操作
14. TagAction - 标签操作
15. UnavailabilitySwitchAction - 不可用开关操作
16. VerificationAction - 验证操作

**其他动作类（Other Actions）- 8 个：**
- ActionRequest
- CUAllowableAction
- EndDeviceAction
- ProtectiveActionAdjustment
- ProtectiveActionCollection
- ProtectiveActionEquipment
- ProtectiveActionRegulation
- RemedialActionScheme

**操作类（Operation Classes）- 16 个：**
- ModelOperation 系列（5个）
- Operation 相关类（11个）

---

### 2. GroundAction 的属性定义

#### 问题
- GroundAction 的 description 属性是什么？

#### 答案

`GroundAction` 没有直接定义 `description` 属性，它继承自父类 `SwitchingAction`：

```xml
<owl:DatatypeProperty rdf:about="...SwitchingAction.description">
    <rdfs:domain rdf:resource="...SwitchingAction"/>
    <rdfs:range rdf:resource="http://www.w3.org/2001/XMLSchema#string"/>
    <rdfs:comment>Free text description of this activity.</rdfs:comment>
</owl:DatatypeProperty>
```

**属性说明：**
- **类型**：`DatatypeProperty`（数据类型属性）
- **值域**：`string`（字符串类型）
- **说明**：此活动的自由文本描述

---

### 3. 动作类的 Domain（域）分析

#### 问题
- 这些动作类的 domain 分别是什么？

#### 答案

在 OWL 中，domain（域）表示哪些类可以拥有指向该 Action 的属性。

**主要动作类的 Domain：**

| Action 类 | Domain 类（可以拥有该 Action 的类） |
|-----------|-----------------------------------|
| **GroundAction** | `ACLineSegment`, `ConductingEquipment`, `Ground` |
| **JumperAction** | `ACLineSegment`, `Clamp`, `ConductingEquipment`, `Jumper` |
| **ControlAction** | `Control` |
| **SwitchAction** | `Outage`, `Switch` |
| **SwitchingAction** | `Crew`, `Operator`, `SwitchingEvent`, `SwitchingPlan`, `SwitchingStep` |
| **ProtectiveAction** | `Gate`, `ProtectionEquipment`, `ProtectiveActionCollection` |
| **ClampAction** | `Clamp` |
| **ClearanceAction** | `ClearanceDocument` |
| **CutAction** | `Cut` |
| **EnergyConsumerAction** | `EnergyConsumer` |
| **EnergySourceAction** | `EnergySource` |
| **GenericAction** | `PowerSystemResource` |
| **MeasurementAction** | `Measurement` |
| **ShuntCompensatorAction** | `ShuntCompensator` |
| **TagAction** | `OperationalTag` |
| **VerificationAction** | `PowerSystemResource` |
| **EndDeviceAction** | `EndDeviceControl` |
| **ProtectiveActionAdjustment** | `ConductingEquipment`, `DCConductingEquipment`, `Measurement` |
| **ProtectiveActionEquipment** | `Equipment` |
| **ProtectiveActionRegulation** | `RegulatingControl` |
| **ProtectiveActionCollection** | `ProtectiveAction`, `StageTrigger` |
| **UnavailabilitySwitchAction** | `EquipmentUnavailabilitySchedule`, `Switch` |

---

### 4. OWL ObjectProperty 的含义

#### 问题
- `owl:ObjectProperty` 是什么意思？

#### 答案

**`owl:ObjectProperty`** 是 OWL（Web Ontology Language）中的属性类型，用于表示两个对象（类实例）之间的关系。

**核心区别：**

| 特性 | ObjectProperty | DatatypeProperty |
|------|---------------|------------------|
| **连接类型** | 对象 → 对象 | 对象 → 基本数据类型 |
| **Range 类型** | 另一个类（如 `GroundAction`） | 数据类型（如 `string`, `integer`） |
| **示例** | `ACDCConverter` → `DCTerminals` | `SwitchingAction` → `description` (string) |
| **用途** | 表示实体间的关系 | 表示对象的属性值 |

**示例：**

```xml
<!-- ObjectProperty：对象之间的关系 -->
<owl:ObjectProperty rdf:about="...ConductingEquipment.GroundingAction">
    <rdfs:domain rdf:resource="...ConductingEquipment"/>
    <rdfs:range rdf:resource="...GroundAction"/>
</owl:ObjectProperty>

<!-- DatatypeProperty：对象的属性值 -->
<owl:DatatypeProperty rdf:about="...SwitchingAction.description">
    <rdfs:domain rdf:resource="...SwitchingAction"/>
    <rdfs:range rdf:resource="http://www.w3.org/2001/XMLSchema#string"/>
</owl:DatatypeProperty>
```

---

### 5. GroundAction 作为 ConductingEquipment 的行为

#### 问题
- GroundAction 可以说是 ConductingEquipment 的一个行为吗？

#### 答案

**从语义角度：可以理解为行为**

从业务语义看，`GroundAction` 可以理解为 `ConductingEquipment` 可以执行的一种行为：

```
ConductingEquipment（导电设备）
  └─ 可以执行的行为：
      ├─ GroundingAction（接地操作）
      ├─ JumpingAction（跳线操作）
      └─ ProtectiveActionAdjustment（保护动作调整）
```

**从 OWL 建模角度：是独立的行为类**

在 OWL 中，`GroundAction` 是一个独立的类，不是 `ConductingEquipment` 的属性或方法：

- **关系建模方式**：通过 `ObjectProperty` 建立关联
- **双向关系**：
  - `ConductingEquipment.GroundingAction` → `GroundAction`
  - `GroundAction.GroundedEquipment` → `ConductingEquipment`

**两种理解方式对比：**

| 视角 | 理解方式 |
|------|---------|
| **业务语义** | "导电设备可以执行接地操作" ✅ |
| **OWL 建模** | "GroundAction 是独立的行为类，通过关系与设备关联" ✅ |
| **面向对象** | "GroundAction 不是 ConductingEquipment 的方法" ❌ |

**CIM 的设计理念：**
- 将行为建模为独立的类（如 `GroundAction`），而不是设备的方法
- 原因：
  1. 行为可以独立存在（有生命周期、时间戳等）
  2. 行为可以被多个设备共享或引用
  3. 行为可以记录历史（时间旅行）
  4. 行为可以关联更多信息（操作员、班组、计划时间等）

---

### 6. Action 作为行为参数 vs 独立实体

#### 问题
- 能不能认为 ConductingEquipment 的行为是执行 action，而 GroundAction 是行为参数？

#### 答案

**两种理解都合理：**

**理解 1：Action 作为参数（概念上）**
```
ConductingEquipment
  └─ performAction(action: GroundAction)  // 方法
      └─ GroundAction  // 参数对象
```

**理解 2：Action 作为独立实体（CIM 的方式）**
```
ConductingEquipment ←→ GroundAction
  (设备)              (操作记录)
```

**对比：**

| 维度 | Action 作为参数 | Action 作为实体（CIM） |
|------|---------------|---------------------|
| **业务语义** | ✅ 设备执行操作 | ✅ 设备关联操作记录 |
| **编程模型** | ✅ 方法调用 | ✅ 实体关联 |
| **数据持久化** | ❌ 操作记录难以独立保存 | ✅ 操作记录独立存在 |
| **历史查询** | ❌ 难以查询历史操作 | ✅ 可以时间旅行查询 |
| **生命周期** | ❌ 与设备绑定 | ✅ 独立生命周期 |

**建议：**
- **概念理解**：可以理解为"设备执行 Action，Action 是参数"，这有助于理解业务语义
- **数据建模**：将 Action 建模为独立实体，便于持久化、历史查询和关联

**类比：**
- **银行账户与交易记录**：
  - 概念上：账户执行交易（交易是参数）
  - 数据上：交易记录是独立实体，关联到账户
- **设备与操作**：
  - 概念上：设备执行操作（操作是参数）
  - 数据上：操作记录是独立实体，关联到设备

---

### 7. 更好的建模方式建议

#### 问题
- 面向对象思维方式将 `performAction(action: GroundAction)` 作为 `ConductingEquipment` 的方法来建模，有什么更好的建议？

#### 答案

面向对象的方法调用方式不适合电力系统的操作记录场景。以下是几种更合适的建模方式：

#### 7.1 事件溯源（Event Sourcing）模式

将操作建模为不可变的事件序列：

```
ConductingEquipment (当前状态)
  ↓
Event Stream (事件流)
  ├─ GroundingPlacedEvent (接地放置事件)
  ├─ GroundingRemovedEvent (接地移除事件)
  └─ StateChangedEvent (状态变更事件)
```

**优点：**
- 完整的历史记录
- 可以重建任意时间点的状态
- 支持时间旅行查询
- 符合审计要求

**实现方式：**
```python
# 事件类
class GroundingPlacedEvent:
    equipment_id: str
    timestamp: datetime
    operator: str
    kind: TempEquipActionKind
    
# 设备状态通过事件重建
class ConductingEquipment:
    def apply_event(self, event: Event):
        # 应用事件，更新状态
        pass
    
    def get_state_at(self, timestamp: datetime):
        # 重放事件到指定时间点
        pass
```

#### 7.2 命令模式（Command Pattern）+ 领域事件

将操作分为命令和事件：

```
Command (命令) → 执行 → Event (事件) → 更新状态
```

```
GroundingCommand (命令)
  ↓ 执行
GroundingExecutedEvent (领域事件)
  ↓ 更新
ConductingEquipmentState (设备状态)
```

**优点：**
- 命令可以排队、撤销、重试
- 事件可以触发其他业务逻辑
- 命令和事件分离，职责清晰

**实现方式：**
```python
class GroundingCommand:
    equipment_id: str
    action_type: str  # "place" or "remove"
    operator: str
    
    def execute(self) -> GroundingExecutedEvent:
        # 执行业务逻辑
        return GroundingExecutedEvent(...)

class GroundingExecutedEvent:
    command_id: str
    equipment_id: str
    executed_at: datetime
    result: str  # "success" or "failure"
```

#### 7.3 状态机 + 操作记录

将设备状态建模为状态机，操作作为状态转换：

```
ConductingEquipment State Machine:
  [Normal] --(place_ground)--> [Grounded]
  [Grounded] --(remove_ground)--> [Normal]
```

每个状态转换都记录操作：

```python
class EquipmentState:
    current_state: str  # "normal", "grounded", "maintenance"
    state_history: List[StateTransition]
    
class StateTransition:
    from_state: str
    to_state: str
    action: GroundAction  # 触发转换的操作
    timestamp: datetime
    operator: str
```

**优点：**
- 状态转换清晰
- 可以验证操作的合法性（状态机约束）
- 操作记录与状态变更关联

#### 7.4 CQRS（命令查询职责分离）

将写操作（命令）和读操作（查询）分离：

```
Write Side (命令侧):
  GroundingCommand → CommandHandler → Event Store
  
Read Side (查询侧):
  Event Store → Read Model → EquipmentView
```

**优点：**
- 写模型专注于业务规则
- 读模型可以优化查询性能
- 可以有不同的数据模型

**实现方式：**
```python
# 命令侧
class GroundingCommandHandler:
    def handle(self, command: GroundingCommand):
        # 验证业务规则
        # 生成事件
        event_store.append(event)

# 查询侧
class EquipmentQueryService:
    def get_equipment_with_actions(self, equipment_id: str):
        # 从读模型查询
        return equipment_view
```

#### 7.5 时间序列数据模型

将操作记录作为时间序列数据：

```
Time Series Data:
  equipment_id | timestamp | action_type | state | operator
  CE-001       | 10:30:00  | place       | grounded | 张三
  CE-001       | 14:00:00  | remove      | normal  | 李四
```

**优点：**
- 适合时间旅行查询
- 可以按时间范围聚合
- 支持流式处理

#### 7.6 领域驱动设计（DDD）的聚合根模式

将设备和操作记录作为聚合：

```
EquipmentAggregate (聚合根)
  ├─ ConductingEquipment (实体)
  └─ ActionHistory (值对象集合)
      ├─ GroundAction (值对象)
      ├─ JumperAction (值对象)
      └─ ...
```

**优点：**
- 业务边界清晰
- 保证一致性
- 操作记录属于设备聚合

#### 7.7 推荐方案：事件溯源 + 状态机组合

结合事件溯源和状态机：

```python
class ConductingEquipment:
    """设备聚合根"""
    def __init__(self, equipment_id: str):
        self.equipment_id = equipment_id
        self.current_state = "normal"
        self.events: List[DomainEvent] = []
    
    def place_ground(self, command: PlaceGroundCommand):
        """放置接地操作"""
        # 1. 验证状态机约束
        if self.current_state != "normal":
            raise InvalidStateTransition()
        
        # 2. 创建领域事件
        event = GroundingPlacedEvent(
            equipment_id=self.equipment_id,
            timestamp=datetime.now(),
            operator=command.operator,
            kind=command.kind
        )
        
        # 3. 应用事件（更新状态）
        self.apply_event(event)
        
        # 4. 记录事件
        self.events.append(event)
        
        return event
    
    def apply_event(self, event: DomainEvent):
        """应用事件，更新状态"""
        if isinstance(event, GroundingPlacedEvent):
            self.current_state = "grounded"
        elif isinstance(event, GroundingRemovedEvent):
            self.current_state = "normal"
    
    def get_state_at(self, timestamp: datetime):
        """时间旅行：获取指定时间点的状态"""
        state = "normal"
        for event in self.events:
            if event.timestamp <= timestamp:
                self._apply_event_to_state(event, state)
        return state
```

**为什么这种方式更好？**

1. **业务语义清晰**：操作是事件，不是方法调用
2. **历史完整**：所有事件都被记录，可以重建任意状态
3. **审计友好**：每个操作都有完整的事件记录
4. **时间旅行**：可以查询任意时间点的设备状态
5. **扩展性好**：新的事件类型可以轻松添加
6. **状态一致性**：通过状态机保证操作合法性

**与 CIM 的对应关系：**

- CIM 的 `GroundAction` → 领域事件（`GroundingPlacedEvent`）
- CIM 的 `ConductingEquipment` → 聚合根（`EquipmentAggregate`）
- CIM 的 `SwitchingAction.executedDateTime` → 事件的时间戳
- CIM 的 `SwitchingAction.operator` → 事件的元数据

这种方式既保留了 CIM 的语义，又提供了更好的编程模型和查询能力。

---

### 8. ObjectProperty 的基数（Cardinality）问题

#### 问题
- 根据定义，`ConductingEquipment.GroundingAction` 好像只能记录一个 GroundingAction？

#### 答案

**关键点：没有基数约束 = 允许多值**

在 OWL 中，如果没有明确指定基数约束（`owl:cardinality`、`owl:maxCardinality`、`owl:minCardinality`），ObjectProperty **默认允许多个值**。

**对比示例：**

```xml
<!-- 示例 1：复数形式（明确表示多个） -->
<owl:ObjectProperty rdf:about="...ACDCConverter.DCTerminals">
    <rdfs:comment>A converter has two DC converter terminals.</rdfs:comment>
</owl:ObjectProperty>

<!-- 示例 2：单数形式（但同样允许多值） -->
<owl:ObjectProperty rdf:about="...ConductingEquipment.GroundingAction">
    <!-- 没有基数约束，默认允许多个值 -->
</owl:ObjectProperty>
```

**RDF/OWL 的多值特性：**

在 RDF 中，一个属性可以有多个值：

```xml
<!-- 一个设备可以有多个接地操作 -->
<ConductingEquipment rdf:ID="CE-001">
    <GroundingAction rdf:resource="#GA-001"/>  <!-- 操作 1 -->
    <GroundingAction rdf:resource="#GA-002"/>  <!-- 操作 2 -->
    <GroundingAction rdf:resource="#GA-003"/>  <!-- 操作 3 -->
</ConductingEquipment>
```

**如何限制为单值？**

如果要限制为只能有一个值，需要添加基数约束：

```xml
<owl:ObjectProperty rdf:about="...ConductingEquipment.GroundingAction">
    <owl:maxCardinality rdf:datatype="&xsd;nonNegativeInteger">1</owl:maxCardinality>
</owl:ObjectProperty>
```

但在 CIM 中，没有这样的约束，所以允许多个值。

**总结：**
- ✅ 虽然属性名是单数形式（`GroundingAction`），但语义上允许多值
- ✅ 没有基数约束时，默认允许多个值
- ✅ 一个 `ConductingEquipment` 可以关联多个 `GroundAction`（历史操作记录）

---

## 关键概念总结

### 1. OWL 属性类型
- **ObjectProperty**：对象之间的关系（如设备 → 操作）
- **DatatypeProperty**：对象的属性值（如操作 → 描述字符串）

### 2. Domain 和 Range
- **Domain（域）**：属性的起点，表示哪些类的实例可以拥有这个属性
- **Range（值域）**：属性的终点，表示属性值必须是什么类型
  - ObjectProperty 的 range：另一个类
  - DatatypeProperty 的 range：基本数据类型

### 3. 基数（Cardinality）
- **默认行为**：没有基数约束时，属性允许多个值
- **单值限制**：需要明确指定 `owl:maxCardinality` 为 1

### 4. CIM 的建模理念
- **行为作为独立实体**：Action 类不是设备的方法，而是独立的操作记录
- **历史记录**：支持时间旅行查询，可以记录所有历史操作
- **关联关系**：通过 ObjectProperty 建立设备与操作之间的双向关系

### 5. 更好的建模方式
- **事件溯源（Event Sourcing）**：将操作建模为不可变的事件序列，支持完整历史记录和时间旅行
- **命令模式 + 领域事件**：分离命令和事件，支持命令排队、撤销、重试
- **状态机 + 操作记录**：通过状态机验证操作合法性，记录状态转换历史
- **CQRS（命令查询职责分离）**：分离写模型和读模型，优化查询性能
- **时间序列数据模型**：将操作记录作为时间序列，适合时间范围查询和流式处理
- **DDD 聚合根模式**：将设备和操作记录作为聚合，保证业务边界和一致性
- **推荐方案**：事件溯源 + 状态机组合，既保留 CIM 语义，又提供更好的编程模型

---

## 实际应用示例

### 设备与操作的关联

```xml
<!-- 设备 -->
<ConductingEquipment rdf:ID="CE-001">
    <name>220kV 母线</name>
</ConductingEquipment>

<!-- 操作记录 1 -->
<GroundAction rdf:ID="GA-001">
    <GroundedEquipment rdf:resource="#CE-001"/>
    <executedDateTime>2025-01-15T10:30:00</executedDateTime>
    <kind rdf:resource="#TempEquipActionKind.place"/>
    <description>在 220kV 母线接地</description>
</GroundAction>

<!-- 操作记录 2 -->
<GroundAction rdf:ID="GA-002">
    <GroundedEquipment rdf:resource="#CE-001"/>
    <executedDateTime>2025-01-20T14:00:00</executedDateTime>
    <kind rdf:resource="#TempEquipActionKind.remove"/>
    <description>移除 220kV 母线接地</description>
</GroundAction>

<!-- 设备关联多个操作 -->
<ConductingEquipment rdf:ID="CE-001">
    <GroundingAction rdf:resource="#GA-001"/>
    <GroundingAction rdf:resource="#GA-002"/>
</ConductingEquipment>
```

### SPARQL 查询示例

```sparql
# 查询一个设备的所有接地操作
SELECT ?equipment ?action ?time ?description
WHERE {
    ?equipment rdf:type cim:ConductingEquipment .
    ?equipment cim:ConductingEquipment.GroundingAction ?action .
    ?action cim:SwitchingAction.executedDateTime ?time .
    ?action cim:SwitchingAction.description ?description .
}
ORDER BY ?time
```

---

## 相关文件

- `TheCimOntology.owl` - CIM 本体定义文件
- `scripts/list_action_operation_classes.py` - 列出所有动作和操作类的脚本
- `scripts/extract_action_event_classes.py` - 提取包含 Action/Event 属性的类
- `scripts/extract_action_domains.py` - 提取动作类的 Domain 信息

---

## 参考资料

- [OWL 2 Web Ontology Language Primer](https://www.w3.org/TR/owl2-primer/)
- [RDF Schema 1.1](https://www.w3.org/TR/rdf-schema/)
- [IEC 61970/61968 CIM](https://www.iec.ch/iec61970)

---

**生成时间**：2025年12月9日  
**讨论主题**：CIM 本体中的动作类（Action Classes）和 OWL 属性定义

