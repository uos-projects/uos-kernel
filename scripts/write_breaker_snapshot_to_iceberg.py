#!/usr/bin/env python3
"""
将 Breaker Actor 快照写入 Iceberg 表
演示完整的 Actor → Snapshot → Iceberg 流程
"""

import os
import re
from datetime import datetime, timezone
from typing import Dict, Any
from pyspark.sql import SparkSession
from pyspark.sql.functions import current_timestamp, lit
from pyspark.sql.types import TimestampType, StringType, LongType

# Iceberg 配置
NESSIE_URI = "http://localhost:19120/api/v2"
S3_ENDPOINT = "http://localhost:19000"
WAREHOUSE = "s3a://iceberg/warehouse"
BRANCH = "main"
ICEBERG_PACKAGES = ",".join(
    [
        "org.apache.iceberg:iceberg-spark-runtime-3.5_2.12:1.5.0",
        "org.apache.iceberg:iceberg-aws-bundle:1.5.0",
    ]
)

os.environ.setdefault("AWS_REGION", "us-east-1")
os.environ.setdefault("AWS_ACCESS_KEY_ID", "iceberg")
os.environ.setdefault("AWS_SECRET_ACCESS_KEY", "iceberg_password")


def build_spark() -> SparkSession:
    """构建 SparkSession"""
    return (
        SparkSession.builder.appName("WriteBreakerSnapshot")
        .config("spark.sql.catalog.ontology", "org.apache.iceberg.spark.SparkCatalog")
        .config("spark.jars.packages", ICEBERG_PACKAGES)
        .config(
            "spark.sql.catalog.ontology.catalog-impl",
            "org.apache.iceberg.nessie.NessieCatalog",
        )
        .config("spark.sql.catalog.ontology.uri", NESSIE_URI)
        .config("spark.sql.catalog.ontology.ref", BRANCH)
        .config("spark.sql.catalog.ontology.warehouse", WAREHOUSE)
        .config("spark.sql.catalog.ontology.io-impl", "org.apache.iceberg.aws.s3.S3FileIO")
        .config("spark.sql.catalog.ontology.s3.endpoint", S3_ENDPOINT)
        .config("spark.sql.catalog.ontology.s3.access-key-id", "iceberg")
        .config("spark.sql.catalog.ontology.s3.secret-access-key", "iceberg_password")
        .config("spark.sql.catalog.ontology.s3.path-style-access", "true")
        .config("spark.sql.catalog.ontology.s3.region", "us-east-1")
        .config("spark.hadoop.fs.s3a.endpoint", S3_ENDPOINT)
        .config("spark.hadoop.fs.s3a.path.style.access", "true")
        .config("spark.hadoop.fs.s3a.access.key", "iceberg")
        .config("spark.hadoop.fs.s3a.secret.key", "iceberg_password")
        .config("spark.sql.extensions", "org.apache.iceberg.spark.extensions.IcebergSparkSessionExtensions")
        .getOrCreate()
    )


def camel_to_snake(name: str) -> str:
    """将驼峰命名转换为 snake_case"""
    s1 = re.sub('(.)([A-Z][a-z]+)', r'\1_\2', name)
    return re.sub('([a-z0-9])([A-Z])', r'\1_\2', s1).lower()


def convert_properties_to_table_fields(properties: Dict[str, Any]) -> Dict[str, Any]:
    """将 Actor 属性（CIM 命名）转换为表字段（snake_case）"""
    table_fields = {}
    
    for prop_name, prop_value in properties.items():
        # 转换为 snake_case
        field_name = camel_to_snake(prop_name)
        table_fields[field_name] = prop_value
    
    return table_fields


def create_breaker_snapshot_record(
    actor_id: str,
    owl_class_uri: str,
    sequence: int,
    properties: Dict[str, Any],
    snapshot_id: str = None,
) -> Dict[str, Any]:
    """创建 Breaker 快照记录"""
    now = datetime.now(timezone.utc)
    far_future = datetime(9999, 12, 31, 23, 59, 59, tzinfo=timezone.utc)
    
    # 转换属性名为表字段名
    table_fields = convert_properties_to_table_fields(properties)
    
    # 构建记录
    record = {
        "actor_id": actor_id,
        "owl_class_uri": owl_class_uri,
        "sequence": sequence,
        "timestamp": now,
        "snapshot_id": snapshot_id or f"{actor_id}-snapshot-{sequence}",
        # 时间旅行字段
        "valid_from": now,
        "valid_to": far_future,
        "op_type": "insert",
        "ingestion_ts": now,
    }
    
    # 添加属性字段（确保所有表字段都有值）
    # 表字段列表（从 SQL 中提取）
    table_field_names = {
        "locked", "normal_open", "open", "retained", "switch_on_count",
        "switch_on_date", "aggregate", "in_service", "network_analysis_enabled",
        "normally_in_service", "alias_name", "description", "m_rid", "name"
    }
    
    # 先设置所有字段为 None
    for field_name in table_field_names:
        record[field_name] = None
    
    # 然后从转换后的属性中填充值
    for field_name, value in table_fields.items():
        if field_name in table_field_names:
            # 转换值为字符串（因为表中字段都是 STRING 类型）
            record[field_name] = str(value) if value is not None else None
    
    return record


def write_snapshot_to_iceberg(spark: SparkSession, snapshot_record: Dict[str, Any]):
    """将快照记录写入 Iceberg 表"""
    table_name = "ontology.grid.breaker_snapshots"
    
    # 准备数据：确保所有字段都是字符串类型（因为表中字段都是 STRING）
    clean_record = {}
    for key, value in snapshot_record.items():
        if value is None:
            clean_record[key] = None
        elif isinstance(value, datetime):
            # 时间戳字段保持为 datetime 对象
            clean_record[key] = value
        elif isinstance(value, (int, float)):
            # 数字字段转为字符串
            clean_record[key] = str(value)
        elif isinstance(value, bool):
            # 布尔值转为字符串
            clean_record[key] = str(value).lower()
        else:
            clean_record[key] = str(value) if value is not None else None
    
    # 使用 SQL INSERT 语句写入（避免类型推断问题）
    # 构建字段列表（按表定义顺序）
    columns = [
        "actor_id", "owl_class_uri", "sequence", "timestamp", "snapshot_id",
        "locked", "normal_open", "open", "retained", "switch_on_count",
        "switch_on_date", "aggregate", "in_service", "network_analysis_enabled",
        "normally_in_service", "alias_name", "description", "m_rid", "name",
        "valid_from", "valid_to", "op_type", "ingestion_ts"
    ]
    
    # 定义字段类型映射
    field_types = {
        "sequence": "bigint",
        "timestamp": "timestamp",
        "valid_from": "timestamp",
        "valid_to": "timestamp",
        "ingestion_ts": "timestamp",
    }
    
    # 构建值列表
    values = []
    for col in columns:
        val = clean_record.get(col)
        if val is None:
            values.append("NULL")
        elif col in field_types and field_types[col] == "timestamp":
            # 时间戳字段
            if isinstance(val, datetime):
                ts_str = val.strftime("%Y-%m-%d %H:%M:%S")
                values.append(f"TIMESTAMP '{ts_str}'")
            else:
                values.append("NULL")
        elif col in field_types and field_types[col] == "bigint":
            # 数字字段（不转字符串）
            if isinstance(val, (int, str)):
                values.append(str(val))
            else:
                values.append("NULL")
        elif isinstance(val, str):
            # 字符串字段（转义单引号）
            escaped = val.replace("'", "''")
            values.append(f"'{escaped}'")
        else:
            # 其他类型转为字符串
            values.append(f"'{str(val)}'")
    
    # 构建 INSERT 语句
    columns_str = ", ".join(columns)
    values_str = ", ".join(values)
    
    insert_sql = f"INSERT INTO {table_name} ({columns_str}) VALUES ({values_str})"
    
    # 执行 INSERT
    spark.sql(insert_sql)
    
    print(f"  ✓ Snapshot written to {table_name}")


def query_snapshots(spark: SparkSession, actor_id: str = None):
    """查询快照数据"""
    table_name = "ontology.grid.breaker_snapshots"
    
    if actor_id:
        query = f"SELECT * FROM {table_name} WHERE actor_id = '{actor_id}' ORDER BY sequence DESC"
    else:
        query = f"SELECT * FROM {table_name} ORDER BY sequence DESC LIMIT 10"
    
    print(f"\n=== Querying snapshots ===")
    print(f"SQL: {query}\n")
    
    df = spark.sql(query)
    df.show(truncate=False)


def main():
    print("=" * 60)
    print("Write Breaker Actor Snapshot to Iceberg")
    print("=" * 60)
    
    # 1. 构建 SparkSession
    print("\n[1/4] Building SparkSession...")
    spark = build_spark()
    print("  ✓ SparkSession created")
    
    # 2. 创建 Breaker Actor 快照数据（模拟 Actor 的属性）
    print("\n[2/4] Creating Breaker Actor snapshot...")
    
    # Actor 属性（CIM 原始命名）
    actor_properties = {
        "mRID": "breaker-001",
        "name": "Main Breaker",
        "normalOpen": False,
        "open": False,
        "locked": False,
        "description": "Main circuit breaker for substation",
        "in_service": True,
        "normally_in_service": True,
        "network_analysis_enabled": True,
        "aggregate": False,
        "retained": True,
        "switch_on_count": 0,
    }
    
    # 创建快照记录
    snapshot_record = create_breaker_snapshot_record(
        actor_id="breaker-001",
        owl_class_uri="http://www.iec.ch/TC57/CIM#Breaker",
        sequence=1,
        properties=actor_properties,
    )
    
    print(f"  Actor ID: {snapshot_record['actor_id']}")
    print(f"  OWL Class URI: {snapshot_record['owl_class_uri']}")
    print(f"  Sequence: {snapshot_record['sequence']}")
    print(f"  Properties: {len(actor_properties)} properties")
    
    # 3. 写入 Iceberg 表
    print("\n[3/4] Writing snapshot to Iceberg...")
    try:
        write_snapshot_to_iceberg(spark, snapshot_record)
    except Exception as e:
        print(f"  ✗ Error: {e}")
        spark.stop()
        return
    
    # 4. 查询验证
    print("\n[4/4] Querying snapshots...")
    query_snapshots(spark, actor_id="breaker-001")
    
    # 5. 创建第二个快照（模拟状态变化）
    print("\n[Bonus] Creating second snapshot (state change)...")
    actor_properties_updated = actor_properties.copy()
    actor_properties_updated["open"] = True  # 断路器打开
    actor_properties_updated["switch_on_count"] = 1
    
    snapshot_record_2 = create_breaker_snapshot_record(
        actor_id="breaker-001",
        owl_class_uri="http://www.iec.ch/TC57/CIM#Breaker",
        sequence=2,
        properties=actor_properties_updated,
    )
    
    write_snapshot_to_iceberg(spark, snapshot_record_2)
    query_snapshots(spark, actor_id="breaker-001")
    
    # 6. 时间旅行查询示例
    print("\n[Bonus] Time travel query example:")
    print("  SELECT * FROM ontology.grid.breaker_snapshots")
    print("  FOR SYSTEM_TIME AS OF <timestamp>")
    print("  WHERE actor_id = 'breaker-001'")
    
    spark.stop()
    
    print("\n" + "=" * 60)
    print("Completed!")
    print("=" * 60)


if __name__ == "__main__":
    main()
