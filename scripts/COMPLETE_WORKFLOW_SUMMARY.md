# 完整工作流程总结：OWL → Actor → Iceberg

## ✅ 已完成的工作

### 1. OWL 模型解析和代码生成

**生成器**: `scripts/owl_actor_table_generator.py`

- ✅ 解析了 1933 个 OWL 类
- ✅ 解析了 8988 个属性
- ✅ 识别了 172 个 PowerSystemResource 类
- ✅ 生成了通用的 `CIMResourceActor` Go 代码
- ✅ 为每个 PowerSystemResource 类生成了独立的 Iceberg 表 SQL

**生成的文件**:
- `actors/cim_resource_actor.go` - 通用的 CIM 资源 Actor
- `scripts/generated_iceberg_tables.sql` - 4388 行 SQL，包含 172 个表的定义

### 2. Breaker Actor 示例

**示例代码**: `actors/cmd/create_breaker_example/main.go`

- ✅ 创建 Breaker Actor（带 Control Capacity）
- ✅ 设置 Actor 属性（CIM 原始命名）
- ✅ 创建快照对象

### 3. Iceberg 表创建

**SQL 脚本**: `scripts/create_breaker_table.sql`

- ✅ 创建命名空间 `ontology.grid`
- ✅ 创建表 `ontology.grid.breaker_snapshots`
- ✅ 表包含所有继承的属性（来自 ProtectedSwitch、Switch、ConductingEquipment 等）

**执行方式**:
```bash
source .venv/bin/activate
python3 scripts/execute_breaker_table_sql.py
```

### 4. Actor 快照写入 Iceberg

**写入脚本**: `scripts/write_breaker_snapshot_to_iceberg.py`

- ✅ 创建 Breaker Actor 快照数据
- ✅ 将属性名从 CIM 命名转换为 snake_case
- ✅ 写入 Iceberg 表
- ✅ 创建多个快照（模拟状态变化）

**执行结果**:
- Sequence 1: `open=False`, `switch_on_count=0`
- Sequence 2: `open=True`, `switch_on_count=1`

### 5. 时间旅行查询验证

**查询结果**:
- ✅ 可以查询当前状态（sequence=2）
- ✅ 可以查询历史状态（sequence=1）
- ✅ 可以查询时间序列（所有快照）

## 📊 数据流程

```
OWL 模型 (TheCimOntology.owl)
    ↓
生成器解析 OWL
    ↓
┌─────────────────────┬─────────────────────┐
│                     │                     │
生成 Actor 代码        生成 Iceberg 表 SQL
│                     │                     │
↓                     ↓                     ↓
CIMResourceActor      breaker_snapshots     执行 SQL 创建表
    ↓                     ↑                     │
创建 Actor 实例            │                     │
    ↓                     │                     │
设置属性                  │                     │
    ↓                     │                     │
创建快照                  │                     │
    ↓                     │                     │
转换属性名                │                     │
    ↓                     │                     │
写入 Iceberg ────────────┘                     │
    ↓                                           │
时间旅行查询 ←─────────────────────────────────┘
```

## 🔑 关键设计特点

### 1. Actor 设计
- ✅ **不体现 OWL 层次关系**：使用通用的 `CIMResourceActor`
- ✅ **维护 OWL URI**：通过 `OWLClassURI` 字段维护语义引用
- ✅ **属性存储**：使用 `map[string]interface{}` 存储所有属性

### 2. 表设计
- ✅ **每个类独立的表**：每个 OWL 类（包括子类）都有独立的快照表
- ✅ **包含继承的属性**：表结构包含该类及其所有父类的 DatatypeProperty
- ✅ **时间旅行支持**：包含 `valid_from`/`valid_to` 字段

### 3. 属性映射
- ✅ **命名转换**：Actor 使用 CIM 命名（驼峰），表使用 snake_case
- ✅ **类型转换**：写入时自动转换类型（字符串、数字、时间戳）
- ✅ **字段过滤**：只写入表中存在的字段

## 📝 使用示例

### 创建 Actor 并写入快照

```python
# 1. 创建 Actor 快照数据
snapshot_record = create_breaker_snapshot_record(
    actor_id="breaker-001",
    owl_class_uri="http://www.iec.ch/TC57/CIM#Breaker",
    sequence=1,
    properties={
        "mRID": "breaker-001",
        "name": "Main Breaker",
        "normalOpen": False,
        "open": False,
        # ... 更多属性
    }
)

# 2. 写入 Iceberg
write_snapshot_to_iceberg(spark, snapshot_record)

# 3. 查询验证
query_snapshots(spark, actor_id="breaker-001")
```

### 时间旅行查询

```sql
-- 查询当前状态
SELECT * FROM ontology.grid.breaker_snapshots 
WHERE actor_id = 'breaker-001' AND sequence = 2;

-- 查询历史状态
SELECT * FROM ontology.grid.breaker_snapshots 
WHERE actor_id = 'breaker-001' AND sequence = 1;

-- 查询时间序列
SELECT * FROM ontology.grid.breaker_snapshots 
WHERE actor_id = 'breaker-001' 
ORDER BY sequence;
```

## 🎯 验证结果

### 表创建成功
```
Tables found:
  - grid.breaker
  - grid.breaker_snapshots
```

### 数据写入成功
```
Sequence 1: open=False, switch_on_count=0
Sequence 2: open=True, switch_on_count=1
```

### 查询验证成功
- ✅ 可以查询所有快照
- ✅ 可以按 sequence 查询特定快照
- ✅ 时间旅行查询正常工作

## 📚 相关文件

1. **生成器**: `scripts/owl_actor_table_generator.py`
2. **Actor 代码**: `actors/cim_resource_actor.go`
3. **示例代码**: `actors/cmd/create_breaker_example/main.go`
4. **SQL 脚本**: `scripts/create_breaker_table.sql`
5. **写入脚本**: `scripts/write_breaker_snapshot_to_iceberg.py`
6. **执行脚本**: `scripts/execute_breaker_table_sql.py`

## 🚀 下一步

1. **扩展其他 Actor 类型**：为其他 PowerSystemResource 类创建 Actor 和表
2. **实现 Go 写入逻辑**：在 Go 代码中实现直接写入 Iceberg 的功能
3. **实现快照恢复**：从 Iceberg 表恢复 Actor 状态
4. **实现事件溯源**：将 Actor 状态变化记录为事件
