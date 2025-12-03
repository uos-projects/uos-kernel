# 语义模型管理和自动映射功能说明

## 📋 功能概述

实现了从 CIM 语义模型（`cim_scope.yaml`）到 Iceberg 表的自动映射和表创建功能，参考 Apache Atlas 的设计理念。

## ✅ 已实现功能

### 1. 语义模型解析模块 (`ontology/`)

- **`model.py`**: 定义语义模型数据结构
  - `Entity`: 实体类型
  - `Package`: 包定义
  - `Relationship`: 关系定义
  - `SemanticModel`: 完整语义模型

- **`parser.py`**: 解析 `cim_scope.yaml`
  - 解析包、实体、关系
  - 支持继承关系处理
  - 提供查询接口

### 2. 映射规则配置 (`mapping/rules.py`)

- **命名空间映射**: 包 ID → 命名空间
  - `core_topology` → `grid`
  - `asset_mgmt` → `assets`
  - `metering_market` → `metering`

- **属性类型映射**: CIM 属性 → Iceberg 字段类型
  - `mRID` → `STRING`
  - `nominalVoltage` → `DOUBLE`
  - `normalOpen` → `BOOLEAN`
  - `installDate` → `TIMESTAMP`

- **字段名转换**: 驼峰命名 → 下划线命名
  - `nominalVoltage` → `nominal_voltage_kv`
  - `normalOpen` → `normal_open`

- **分区策略**: 实体类型 → 分区字段
  - `Substation` → `['region']`
  - `VoltageLevel` → `['substation_id']`
  - `Breaker` → `['voltage_level_id']`

### 3. 自动映射引擎 (`mapping/mapper.py`)

- **`EntityMapper`**: 实体到表的映射器
  - 自动确定命名空间
  - 生成表名（驼峰转下划线）
  - 收集所有属性（包括继承的）
  - 根据关系自动添加外键字段
  - 添加通用时间旅行字段
  - 确定分区字段

### 4. SQL 生成器 (`table_gen/generator.py`)

- **`SQLGenerator`**: 生成 CREATE TABLE SQL
  - 生成命名空间创建 SQL
  - 生成表创建 SQL（包含列、分区、注释）
  - 支持 DROP TABLE SQL

### 5. API 接口 (`server/app.py`)

#### 语义模型查询 API

- `GET /api/ontology/packages` - 获取所有包
- `GET /api/ontology/entities` - 获取所有实体（支持按包过滤）
- `GET /api/ontology/entity/{entity_name}` - 获取实体详情

#### 映射和表生成 API

- `GET /api/mapping/entity/{entity_name}` - 获取实体的表结构（不创建）
- `POST /api/tables/create/{entity_name}` - 创建表（支持 dry_run）
- `GET /api/tables/mappings` - 获取所有映射关系

### 6. 可视化 UI (`web/ontology.html`)

- **左侧边栏**: 包和实体树形结构
- **主内容区**: 实体详情展示
  - 属性列表
  - 关系列表
  - Iceberg 表结构预览
  - SQL 预览
  - 表创建按钮

## 🚀 使用方法

### 1. 启动服务

```bash
cd /home/chun/Develop/uos-projects/uos-kernel
source .venv/bin/activate
uvicorn server.app:app --host 0.0.0.0 --port 8000
```

### 2. 访问语义模型管理界面

打开浏览器访问: http://localhost:8000/ontology.html

### 3. 使用 API

#### 查询语义模型

```bash
# 获取所有包
curl http://localhost:8000/api/ontology/packages

# 获取所有实体
curl http://localhost:8000/api/ontology/entities

# 获取实体详情
curl http://localhost:8000/api/ontology/entity/Substation
```

#### 生成表结构

```bash
# 获取表结构（不创建）
curl http://localhost:8000/api/mapping/entity/Substation

# 预览 SQL（dry run）
curl -X POST "http://localhost:8000/api/tables/create/Substation?dry_run=true"

# 实际创建表
curl -X POST "http://localhost:8000/api/tables/create/Substation?dry_run=false"
```

### 4. 使用 Python 脚本

```python
from ontology.parser import get_semantic_model
from mapping.mapper import EntityMapper
from table_gen.generator import SQLGenerator

# 获取语义模型
model = get_semantic_model()

# 创建映射器
mapper = EntityMapper(model)

# 映射实体到表结构
schema = mapper.map_entity_by_name("Substation")

# 生成 SQL
sqls = SQLGenerator.generate_table_sqls(schema)
for sql in sqls:
    print(sql)
```

## 📊 映射示例

### 输入：语义模型（YAML）

```yaml
entities:
  - name: Substation
    key_attributes: [mRID, name, description, region]
```

### 输出：Iceberg 表

```sql
CREATE TABLE IF NOT EXISTS ontology.grid.substation (
    description STRING COMMENT 'CIM attribute: description',
    entity_id STRING COMMENT 'CIM attribute: mRID',
    name STRING COMMENT 'CIM attribute: name',
    region STRING COMMENT 'CIM attribute: region',
    valid_from TIMESTAMP COMMENT 'Common time travel field',
    valid_to TIMESTAMP COMMENT 'Common time travel field',
    op_type STRING COMMENT 'Common time travel field',
    ingestion_ts TIMESTAMP COMMENT 'Common time travel field'
)
USING ICEBERG
PARTITIONED BY (region)
COMMENT 'Table for CIM entity: Substation'
```

## 🔧 配置说明

### 修改映射规则

编辑 `mapping/rules.py`:

```python
# 添加新的命名空间映射
NAMESPACE_MAPPING["new_package"] = "new_namespace"

# 添加新的属性类型映射
ATTRIBUTE_TYPE_MAPPING["newAttribute"] = "STRING"

# 添加新的分区策略
PARTITION_STRATEGY["NewEntity"] = ["partition_field"]
```

## 📁 文件结构

```
uos-kernel/
├── ontology/
│   ├── __init__.py
│   ├── model.py          # 数据模型
│   ├── parser.py          # YAML 解析器
│   └── cim_scope.yaml     # 语义模型定义
│
├── mapping/
│   ├── __init__.py
│   ├── rules.py           # 映射规则配置
│   └── mapper.py          # 映射引擎
│
├── table_gen/
│   ├── __init__.py
│   └── generator.py       # SQL 生成器
│
├── server/
│   └── app.py             # API 服务（已扩展）
│
└── web/
    └── ontology.html      # 语义模型可视化 UI
```

## 🎯 核心特性

1. **自动映射**: 从语义模型自动生成 Iceberg 表结构
2. **继承支持**: 自动处理实体继承关系
3. **关系处理**: 根据关系自动添加外键字段
4. **分区策略**: 自动确定分区字段
5. **通用字段**: 自动添加时间旅行字段
6. **可视化**: Web UI 展示语义模型和表结构
7. **API 驱动**: RESTful API 支持程序化操作

## 🔮 未来改进

1. **映射关系持久化**: 存储映射关系到数据库
2. **版本管理**: 支持语义模型版本管理
3. **关系表自动生成**: 自动生成关系表
4. **图形化关系视图**: 使用 D3.js 展示实体关系图
5. **映射规则可视化编辑**: UI 编辑映射规则
6. **表结构版本管理**: 跟踪表结构变更历史

## 📝 注意事项

1. **分区字段**: 确保分区字段在列中存在，否则会自动过滤
2. **外键字段**: 根据 `partOf` 关系自动添加外键字段（如 `substation_id`）
3. **字段名冲突**: 如果属性名映射后与通用字段冲突，会跳过通用字段
4. **继承属性**: 会递归收集父类的所有属性

## 🐛 已知问题

1. VoltageLevel 的分区字段需要从关系中推断（已实现）
2. 某些实体可能没有合适的分区字段，会使用默认策略

