# Actor 系统设计讨论总结

**日期**: 2025年12月9日  
**主题**: Go Actor 系统实现与 CIM 模型映射

---

## 目录

1. [轻量级 Go-Actor 系统实现](#1-轻量级-go-actor-系统实现)
2. [CIM 模型映射设计](#2-cim-模型映射设计)
3. [代码重构与优化](#3-代码重构与优化)
4. [Actor 之间消息传递](#4-actor-之间消息传递)
5. [PowerSystemResource 关联关系分析](#5-powersystemresource-关联关系分析)

---

## 1. 轻量级 Go-Actor 系统实现

### 1.1 核心组件

创建了一个轻量级的 Go Actor 系统，包含以下核心组件：

- **BaseActor**: Actor 的基础实现，包含邮箱（mailbox）和消息处理循环
- **System**: Actor 系统管理器，负责 Actor 的注册、查找和消息路由
- **ActorBehavior**: Actor 行为接口，定义消息处理逻辑
- **Message**: 消息类型接口

### 1.2 文件结构

```
actors/
├── actor.go                    # 基础 Actor 实现
├── system.go                   # Actor 系统管理器
├── actor_ref.go                # ActorRef 实现（Actor 之间通信）
├── capacity_factory.go         # Capacity 工厂
├── resource_actor.go           # PowerSystemResourceActor
├── capacities/                 # Capacity 实现目录
│   ├── accumulator_reset.go
│   ├── command.go
│   ├── raise_lower_command.go
│   └── set_point.go
└── cmd/                        # 示例程序
    ├── example/
    ├── resource_example/
    └── actor_communication_example/
```

---

## 2. CIM 模型映射设计

### 2.1 核心映射关系

将 CIM (Common Information Model) 模型映射到 Actor 系统：

- **PowerSystemResource** → `PowerSystemResourceActor`
  - 每个 PowerSystemResource 或其子类（如 Breaker、SynchronousMachine、ACLineSegment）对应一个 Actor
  - Actor 代表电力系统中的物理或逻辑资源

- **Control** → `Capacity` (interface)
  - 每个 Control 类或其子类对应一种能力（Capability）
  - 命名转换：`AccumulatorReset` → `AccumulatorResetCapacity`
  - 其他示例：
    - `Command` → `CommandCapacity`
    - `RaiseLowerCommand` → `RaiseLowerCommandCapacity`
    - `SetPoint` → `SetPointCapacity`

- **关联关系** → Actor 具有能力
  - 如果 `Control` 与 `PowerSystemResource` 存在关联（通过 `Control.PowerSystemResource` 属性），则对应的 Actor 具有该 Control 对应的 Capacity
  - 例如：`BE-Line_2` (ACLineSegment) 关联了 `ACC_RESET_1` (AccumulatorReset)，则 `BE-Line_2` Actor 具有 `AccumulatorResetCapacity`

### 2.2 Capacity Interface 设计

使用 Go 的 interface 实现 Capacity：

```go
type Capacity interface {
    Name() string
    CanHandle(msg Message) bool
    Execute(ctx context.Context, msg Message) error
    ControlID() string
}
```

### 2.3 设计优势

1. **动态能力发现**: 通过查询关联关系确定 Actor 具备的能力
2. **类型安全**: 能力与 Control 类型一一对应
3. **可扩展性**: 新增 Control 类型只需实现 Capacity interface
4. **符合领域模型**: 完美映射 CIM 的语义关系

---

## 3. 代码重构与优化

### 3.1 文件拆分决策

**问题**: 不同的 Capacity 和 Message 是否需要拆分到不同文件？

**决策**:
- ✅ **Capacity 拆分**: 每个 Capacity 实现独立文件，放在 `capacities/` 目录
- ✅ **Message 拆分**: 消息类型与对应的 Capacity 放在同一个文件中
- ✅ **保持通用 Actor**: 不為 PowerSystemResource 的每个子类创建不同的 Actor 类型

**理由**:
- 可扩展性：CIM 有 22+ 个 Control 类，每个对应一个 Capacity
- 关注点分离：每个 Capacity 独立，代码更清晰
- 维护性：修改某个 Capacity 不影响其他文件
- 避免类型爆炸：通过 Capacity 组合实现差异化，而不是通过继承

### 3.2 删除不必要的文件

- **删除 `actors/capacity.go`**: 只是类型别名，没有实际使用
- **删除 `actors/example.go`**: 示例代码移到测试文件和独立的示例程序中

### 3.3 CapacityFactory 位置

**决策**: 放在 `actors/` 目录下

**理由**:
- 它是 actors 系统的工具/服务
- 用户通过 actors 包使用它
- capacities 包应该只包含 Capacity 实现，保持单一职责

---

## 4. Actor 之间消息传递

### 4.1 实现 ActorRef 模式

**问题**: Actor 之间可以发消息吗？

**实现**: 添加了 ActorRef 模式，支持 Actor 之间的消息传递

**核心组件**:

1. **ActorRef** (`actors/actor_ref.go`):
   - Actor 的引用，用于 Actor 之间发送消息
   - `Send()` / `Tell()` - 通过 ActorRef 发送消息

2. **BaseActor 增强**:
   - `SetSystem()` - 设置 System 引用（由 System 自动调用）
   - `GetRef()` - 获取其他 Actor 的 ActorRef
   - `SendTo()` / `Tell()` - 直接向其他 Actor 发送消息

3. **System 增强**:
   - `GetRef()` - 获取 ActorRef
   - `MustGetRef()` - 获取 ActorRef（不存在则 panic）
   - 自动设置 Actor 的 System 引用

### 4.2 使用方式

**方式1：通过 System 发送（外部代码）**
```go
system.Send("actor1", "Hello")
```

**方式2：Actor 之间直接通信**
```go
// 在 Actor 的 Handle 方法中
baseActor.SendTo("target-actor", "Hello")
// 或
ref, _ := baseActor.GetRef("target-actor")
ref.Tell("Hello")
```

**方式3：通过 ActorRef**
```go
ref, _ := system.GetRef("actor1")
ref.Send("Hello")
```

---

## 5. PowerSystemResource 关联关系分析

### 5.1 分析目标

使用 cimpyorm 查找 CIM16 schema 中 PowerSystemResource 的所有关联关系，除了 Control 之外还关联了什么。

### 5.2 分析结果

**PowerSystemResource 关联了 4 个不同的类**:

1. **Control** (Controls)
   - 反向属性：PowerSystemResource
   - 基数：0..n
   - 说明：控制命令（已实现为 Capacity）

2. **Measurement** (Measurements)
   - 反向属性：PowerSystemResource
   - 基数：0..n
   - 说明：测量值（如电压、电流、功率等）

3. **Location** (Location)
   - 反向属性：PowerSystemResources
   - 基数：0..1
   - 说明：地理位置信息

4. **DiagramObject** (DiagramObjects)
   - 反向属性：IdentifiedObject
   - 基数：0..n
   - 说明：用于图形化显示的对象

### 5.3 数据类型属性

PowerSystemResource 还包含以下数据类型属性：
- `description` (String) - 描述
- `energyIdentCodeEic` (String) - 能源标识码
- `mRID` (String) - 模型资源标识符
- `name` (String) - 名称
- `shortName` (String) - 短名称

### 5.4 映射建议

这些关联关系可以映射到 Actor 系统中：

- **Control** → Capacity（已实现）
- **Measurement** → 可以添加 Measurement 相关的能力或属性
- **Location** → 可以作为 Actor 的属性
- **DiagramObject** → 可以作为 Actor 的可视化属性

---

## 6. 关键设计决策总结

### 6.1 Capacity 使用 Interface 实现

✅ **决策**: 使用 Go 的 interface 实现 Capacity

**优势**:
- 类型安全：interface 约束能力实现
- 可扩展：新增 Control 类型只需实现 Capacity interface
- 解耦：Actor 与具体能力实现解耦
- 动态发现：运行时根据 CIM 关联关系添加能力

### 6.2 不為 PowerSystemResource 子类创建不同的 Actor 类型

✅ **决策**: 保持通用的 `PowerSystemResourceActor`

**理由**:
- 行为差异通过 Capacity 体现：不同资源类型的差异在于能力组合，而非 Actor 本身
- 避免类型爆炸：CIM 中 PowerSystemResource 子类很多
- 符合组合优于继承：通过 Capacity 组合实现差异化
- 灵活性：运行时动态添加/移除能力

### 6.3 文件组织

✅ **决策**: 
- Capacity 拆分到 `capacities/` 目录，每个 Capacity 一个文件
- Message 与对应的 Capacity 放在同一个文件中
- CapacityFactory 放在 `actors/` 目录下

---

## 7. 实现的功能

### 7.1 核心功能

- ✅ 基础 Actor 系统
- ✅ PowerSystemResourceActor（资源 Actor）
- ✅ Capacity 接口和实现
- ✅ CapacityFactory（能力工厂）
- ✅ 消息路由到 Capacity
- ✅ 动态能力管理（添加/移除/查询）
- ✅ Actor 之间消息传递（ActorRef 模式）

### 7.2 已实现的 Capacity

- `AccumulatorResetCapacity` - 累加器复位能力
- `CommandCapacity` - 命令能力
- `RaiseLowerCommandCapacity` - 升降命令能力
- `SetPointCapacity` - 设定点能力

### 7.3 示例程序

- `cmd/example/` - 基础 Actor 使用示例
- `cmd/resource_example/` - PowerSystemResourceActor 使用示例
- `cmd/actor_communication_example/` - Actor 之间通信示例

---

## 8. 未来扩展方向

### 8.1 基于 PowerSystemResource 关联关系的扩展

根据分析结果，可以考虑以下扩展：

1. **Measurement 支持**
   - 添加 Measurement 相关的 Capacity 或属性
   - 支持测量值的查询和订阅

2. **Location 支持**
   - 在 PowerSystemResourceActor 中添加 Location 属性
   - 支持地理位置查询和 GIS 集成

3. **DiagramObject 支持**
   - 添加可视化相关的属性
   - 支持图形化显示

### 8.2 其他扩展方向

- [ ] 从 CIM 数据自动创建 Actor 和 Capacity
- [ ] Actor 监控和指标
- [ ] 消息超时和重试
- [ ] Actor 监督（Supervision）
- [ ] 分布式 Actor 支持
- [ ] 消息优先级
- [ ] Capacity 的持久化

---

## 9. 相关脚本文件

本次分析过程中创建的脚本文件：

1. **scripts/find_control_classes_in_schema.py**: 查找 CIM16 schema 中所有 Control 类及其关联类型
2. **scripts/find_powersystemresource_relations.py**: 查找 PowerSystemResource 的所有关联关系

---

## 10. 关键发现

1. **Control 类数量**: CIM16 schema 中有 22 个 Control 相关的类
2. **PowerSystemResource 关联**: 除了 Control，还关联了 Measurement、Location、DiagramObject
3. **Capacity 设计**: 使用 interface 实现，支持动态扩展
4. **Actor 通信**: 通过 ActorRef 模式实现 Actor 之间的消息传递

---

**文档生成时间**: 2025年12月9日  
**讨论主题**: Go Actor 系统实现、CIM 模型映射、代码重构优化

