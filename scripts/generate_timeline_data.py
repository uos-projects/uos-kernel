#!/usr/bin/env python3
"""按语义模型生成多时间点的实体和关系数据，创建多个 Iceberg snapshots"""

from __future__ import annotations

import os
import time
from datetime import datetime, timedelta
from typing import List, Dict

from pyspark.sql import SparkSession
from pyspark.sql.functions import current_timestamp, lit
from pyspark.sql.types import TimestampType


NESSIE_URI = "http://localhost:19120/api/v2"
S3_ENDPOINT = "http://localhost:19000"
WAREHOUSE = "s3a://iceberg/warehouse"
CATALOG_NAME = "ontology"
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
    return (
        SparkSession.builder.appName("GenerateTimelineData")
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


def create_tables_if_needed(spark: SparkSession) -> None:
    """创建所有需要的表"""
    spark.sql("CREATE NAMESPACE IF NOT EXISTS ontology.grid")
    spark.sql("CREATE NAMESPACE IF NOT EXISTS ontology.relations")
    spark.sql("CREATE NAMESPACE IF NOT EXISTS ontology.assets")
    spark.sql("CREATE NAMESPACE IF NOT EXISTS ontology.metering")

    # Substation
    spark.sql(
        """
        CREATE TABLE IF NOT EXISTS ontology.grid.substation (
            entity_id STRING,
            name STRING,
            region STRING,
            nominal_voltage_kv DOUBLE,
            valid_from TIMESTAMP,
            valid_to TIMESTAMP,
            op_type STRING,
            ingestion_ts TIMESTAMP
        )
        USING ICEBERG
        PARTITIONED BY (region)
        """
    )

    # VoltageLevel
    spark.sql(
        """
        CREATE TABLE IF NOT EXISTS ontology.grid.voltage_level (
            entity_id STRING,
            substation_id STRING,
            name STRING,
            nominal_voltage_kv DOUBLE,
            valid_from TIMESTAMP,
            valid_to TIMESTAMP,
            op_type STRING,
            ingestion_ts TIMESTAMP
        )
        USING ICEBERG
        PARTITIONED BY (substation_id)
        """
    )

    # PowerTransformer
    spark.sql(
        """
        CREATE TABLE IF NOT EXISTS ontology.grid.power_transformer (
            entity_id STRING,
            name STRING,
            substation_id STRING,
            construction_kind STRING,
            valid_from TIMESTAMP,
            valid_to TIMESTAMP,
            op_type STRING,
            ingestion_ts TIMESTAMP
        )
        USING ICEBERG
        PARTITIONED BY (substation_id)
        """
    )

    # Breaker
    spark.sql(
        """
        CREATE TABLE IF NOT EXISTS ontology.grid.breaker (
            entity_id STRING,
            name STRING,
            voltage_level_id STRING,
            rated_current DOUBLE,
            normal_open BOOLEAN,
            valid_from TIMESTAMP,
            valid_to TIMESTAMP,
            op_type STRING,
            ingestion_ts TIMESTAMP
        )
        USING ICEBERG
        PARTITIONED BY (voltage_level_id)
        """
    )

    # Asset
    spark.sql(
        """
        CREATE TABLE IF NOT EXISTS ontology.assets.asset (
            entity_id STRING,
            asset_id STRING,
            serial_number STRING,
            criticality STRING,
            power_system_resource_id STRING,
            valid_from TIMESTAMP,
            valid_to TIMESTAMP,
            op_type STRING,
            ingestion_ts TIMESTAMP
        )
        USING ICEBERG
        PARTITIONED BY (criticality)
        """
    )

    # Work
    spark.sql(
        """
        CREATE TABLE IF NOT EXISTS ontology.assets.work (
            entity_id STRING,
            work_id STRING,
            priority STRING,
            status STRING,
            target_resource_id STRING,
            valid_from TIMESTAMP,
            valid_to TIMESTAMP,
            op_type STRING,
            ingestion_ts TIMESTAMP
        )
        USING ICEBERG
        PARTITIONED BY (status)
        """
    )

    # Meter
    spark.sql(
        """
        CREATE TABLE IF NOT EXISTS ontology.metering.meter (
            entity_id STRING,
            asset_id STRING,
            install_date TIMESTAMP,
            usage_point_id STRING,
            valid_from TIMESTAMP,
            valid_to TIMESTAMP,
            op_type STRING,
            ingestion_ts TIMESTAMP
        )
        USING ICEBERG
        PARTITIONED BY (usage_point_id)
        """
    )

    # UsagePoint
    spark.sql(
        """
        CREATE TABLE IF NOT EXISTS ontology.metering.usage_point (
            entity_id STRING,
            service_category STRING,
            valid_from TIMESTAMP,
            valid_to TIMESTAMP,
            op_type STRING,
            ingestion_ts TIMESTAMP
        )
        USING ICEBERG
        PARTITIONED BY (service_category)
        """
    )

    # 关系表
    spark.sql(
        """
        CREATE TABLE IF NOT EXISTS ontology.relations.substation_voltage_level (
            rel_id STRING,
            source_id STRING,
            target_id STRING,
            relation_type STRING,
            valid_from TIMESTAMP,
            valid_to TIMESTAMP,
            op_type STRING,
            ingestion_ts TIMESTAMP
        )
        USING ICEBERG
        PARTITIONED BY (relation_type)
        """
    )

    spark.sql(
        """
        CREATE TABLE IF NOT EXISTS ontology.relations.transformer_substation (
            rel_id STRING,
            source_id STRING,
            target_id STRING,
            relation_type STRING,
            valid_from TIMESTAMP,
            valid_to TIMESTAMP,
            op_type STRING,
            ingestion_ts TIMESTAMP
        )
        USING ICEBERG
        PARTITIONED BY (relation_type)
        """
    )

    spark.sql(
        """
        CREATE TABLE IF NOT EXISTS ontology.relations.breaker_voltage_level (
            rel_id STRING,
            source_id STRING,
            target_id STRING,
            relation_type STRING,
            valid_from TIMESTAMP,
            valid_to TIMESTAMP,
            op_type STRING,
            ingestion_ts TIMESTAMP
        )
        USING ICEBERG
        PARTITIONED BY (relation_type)
        """
    )

    spark.sql(
        """
        CREATE TABLE IF NOT EXISTS ontology.relations.asset_resource (
            rel_id STRING,
            source_id STRING,
            target_id STRING,
            relation_type STRING,
            valid_from TIMESTAMP,
            valid_to TIMESTAMP,
            op_type STRING,
            ingestion_ts TIMESTAMP
        )
        USING ICEBERG
        PARTITIONED BY (relation_type)
        """
    )

    spark.sql(
        """
        CREATE TABLE IF NOT EXISTS ontology.relations.work_resource (
            rel_id STRING,
            source_id STRING,
            target_id STRING,
            relation_type STRING,
            valid_from TIMESTAMP,
            valid_to TIMESTAMP,
            op_type STRING,
            ingestion_ts TIMESTAMP
        )
        USING ICEBERG
        PARTITIONED BY (relation_type)
        """
    )

    spark.sql(
        """
        CREATE TABLE IF NOT EXISTS ontology.relations.meter_usage_point (
            rel_id STRING,
            source_id STRING,
            target_id STRING,
            relation_type STRING,
            valid_from TIMESTAMP,
            valid_to TIMESTAMP,
            op_type STRING,
            ingestion_ts TIMESTAMP
        )
        USING ICEBERG
        PARTITIONED BY (relation_type)
        """
    )


def generate_timeline_data() -> List[Dict]:
    """生成多个时间点的数据"""
    base_time = datetime(2025, 1, 1, 0, 0, 0)
    far_future = datetime(9999, 12, 31, 0, 0, 0)
    
    timeline = []
    
    # 时间点 1: 2025-01-01 - 初始基础设施
    t1 = base_time
    timeline.append({
        "timestamp": t1,
        "entities": {
            "substation": [
                {"entity_id": "SS-001", "name": "北城变电站", "region": "North", "nominal_voltage_kv": 220.0},
                {"entity_id": "SS-002", "name": "港口变电站", "region": "East", "nominal_voltage_kv": 110.0},
            ],
            "voltage_level": [
                {"entity_id": "VL-001", "substation_id": "SS-001", "name": "220kV母线", "nominal_voltage_kv": 220.0},
                {"entity_id": "VL-002", "substation_id": "SS-002", "name": "110kV母线", "nominal_voltage_kv": 110.0},
            ],
        },
        "relations": [
            {"rel_id": "REL-SS001-VL001", "source_id": "SS-001", "target_id": "VL-001", "relation_type": "contains"},
            {"rel_id": "REL-SS002-VL002", "source_id": "SS-002", "target_id": "VL-002", "relation_type": "contains"},
        ],
    })
    
    # 时间点 2: 2025-02-01 - 添加变压器和断路器
    t2 = base_time + timedelta(days=31)
    timeline.append({
        "timestamp": t2,
        "entities": {
            "power_transformer": [
                {"entity_id": "PT-001", "name": "主变压器#1", "substation_id": "SS-001", "construction_kind": "threePhase"},
                {"entity_id": "PT-002", "name": "主变压器#2", "substation_id": "SS-002", "construction_kind": "threePhase"},
            ],
            "breaker": [
                {"entity_id": "BR-001", "name": "220kV断路器#1", "voltage_level_id": "VL-001", "rated_current": 2000.0, "normal_open": False},
                {"entity_id": "BR-002", "name": "110kV断路器#1", "voltage_level_id": "VL-002", "rated_current": 1600.0, "normal_open": False},
            ],
        },
        "relations": [
            {"rel_id": "REL-PT001-SS001", "source_id": "PT-001", "target_id": "SS-001", "relation_type": "locatedIn"},
            {"rel_id": "REL-PT002-SS002", "source_id": "PT-002", "target_id": "SS-002", "relation_type": "locatedIn"},
            {"rel_id": "REL-BR001-VL001", "source_id": "BR-001", "target_id": "VL-001", "relation_type": "partOf"},
            {"rel_id": "REL-BR002-VL002", "source_id": "BR-002", "target_id": "VL-002", "relation_type": "partOf"},
        ],
    })
    
    # 时间点 3: 2025-03-01 - 添加资产和运维
    t3 = base_time + timedelta(days=59)
    timeline.append({
        "timestamp": t3,
        "entities": {
            "asset": [
                {"entity_id": "AST-001", "asset_id": "AST-001", "serial_number": "SN-PT-2024-001", "criticality": "high", "power_system_resource_id": "PT-001"},
                {"entity_id": "AST-002", "asset_id": "AST-002", "serial_number": "SN-BR-2024-001", "criticality": "medium", "power_system_resource_id": "BR-001"},
            ],
            "work": [
                {"entity_id": "WK-001", "work_id": "WK-001", "priority": "high", "status": "planned", "target_resource_id": "PT-001"},
            ],
        },
        "relations": [
            {"rel_id": "REL-AST001-PT001", "source_id": "AST-001", "target_id": "PT-001", "relation_type": "realises"},
            {"rel_id": "REL-AST002-BR001", "source_id": "AST-002", "target_id": "BR-001", "relation_type": "realises"},
            {"rel_id": "REL-WK001-PT001", "source_id": "WK-001", "target_id": "PT-001", "relation_type": "targets"},
        ],
    })
    
    # 时间点 4: 2025-04-01 - 添加计量设备
    t4 = base_time + timedelta(days=90)
    timeline.append({
        "timestamp": t4,
        "entities": {
            "usage_point": [
                {"entity_id": "UP-001", "service_category": "residential"},
                {"entity_id": "UP-002", "service_category": "commercial"},
            ],
            "meter": [
                {"entity_id": "MT-001", "asset_id": "MT-001", "install_date": t4, "usage_point_id": "UP-001"},
                {"entity_id": "MT-002", "asset_id": "MT-002", "install_date": t4, "usage_point_id": "UP-002"},
            ],
        },
        "relations": [
            {"rel_id": "REL-MT001-UP001", "source_id": "MT-001", "target_id": "UP-001", "relation_type": "measuredBy"},
            {"rel_id": "REL-MT002-UP002", "source_id": "MT-002", "target_id": "UP-002", "relation_type": "measuredBy"},
        ],
    })
    
    # 时间点 5: 2025-05-01 - 新增变电站和更新关系
    t5 = base_time + timedelta(days=120)
    timeline.append({
        "timestamp": t5,
        "entities": {
            "substation": [
                {"entity_id": "SS-003", "name": "西部工业变电站", "region": "West", "nominal_voltage_kv": 220.0},
            ],
            "voltage_level": [
                {"entity_id": "VL-003", "substation_id": "SS-003", "name": "220kV母线", "nominal_voltage_kv": 220.0},
            ],
            "power_transformer": [
                {"entity_id": "PT-003", "name": "主变压器#3", "substation_id": "SS-003", "construction_kind": "threePhase"},
            ],
        },
        "relations": [
            {"rel_id": "REL-SS003-VL003", "source_id": "SS-003", "target_id": "VL-003", "relation_type": "contains"},
            {"rel_id": "REL-PT003-SS003", "source_id": "PT-003", "target_id": "SS-003", "relation_type": "locatedIn"},
        ],
    })
    
    return timeline


def write_data_at_timestamp(spark: SparkSession, data: Dict, timestamp: datetime) -> None:
    """在指定时间点写入数据"""
    far_future = datetime(9999, 12, 31, 0, 0, 0)
    
    # 写入实体
    for entity_type, entities in data.get("entities", {}).items():
        # 确定命名空间
        if entity_type in ['substation', 'voltage_level', 'power_transformer', 'breaker']:
            namespace = 'grid'
        elif entity_type in ['asset', 'work']:
            namespace = 'assets'
        else:
            namespace = 'metering'
        
        table_name = f"ontology.{namespace}.{entity_type}"
        
        records = []
        for e in entities:
            record = e.copy()
            record["valid_from"] = timestamp
            record["valid_to"] = far_future
            record["op_type"] = "insert"
            records.append(record)
        
        if records:
            df = spark.createDataFrame(records)
            # 确保时间戳列是 TimestampType
            df = df.withColumn("valid_from", lit(timestamp).cast(TimestampType()))
            df = df.withColumn("valid_to", lit(far_future).cast(TimestampType()))
            df = df.withColumn("ingestion_ts", current_timestamp())
            df.writeTo(table_name).append()
            print(f"  ✓ 写入 {len(records)} 个 {entity_type} 实体到 {table_name}")
    
    # 写入关系
    relations = data.get("relations", [])
    if relations:
        # 根据关系类型选择表
        rel_by_type = {}
        for rel in relations:
            rel_type = rel["relation_type"]
            if rel_type == "contains" and "SS" in rel["source_id"]:
                table = "ontology.relations.substation_voltage_level"
            elif rel_type == "locatedIn":
                table = "ontology.relations.transformer_substation"
            elif rel_type == "partOf":
                table = "ontology.relations.breaker_voltage_level"
            elif rel_type in ["realises"]:
                table = "ontology.relations.asset_resource"
            elif rel_type in ["targets"]:
                table = "ontology.relations.work_resource"
            elif rel_type == "measuredBy":
                table = "ontology.relations.meter_usage_point"
            else:
                table = "ontology.relations.substation_voltage_level"
            
            if table not in rel_by_type:
                rel_by_type[table] = []
            rel_by_type[table].append(rel)
        
        for table, rels in rel_by_type.items():
            records = []
            for r in rels:
                record = r.copy()
                record["valid_from"] = timestamp
                record["valid_to"] = far_future
                record["op_type"] = "insert"
                records.append(record)
            
            if records:
                df = spark.createDataFrame(records)
                df = df.withColumn("valid_from", lit(timestamp).cast(TimestampType()))
                df = df.withColumn("valid_to", lit(far_future).cast(TimestampType()))
                df = df.withColumn("ingestion_ts", current_timestamp())
                df.writeTo(table).append()
                print(f"  ✓ 写入 {len(records)} 个关系到 {table}")


def main() -> None:
    spark = build_spark()
    print("创建表结构...")
    create_tables_if_needed(spark)
    
    print("\n生成时间线数据...")
    timeline = generate_timeline_data()
    
    print(f"\n将在 {len(timeline)} 个时间点写入数据，创建多个 snapshots\n")
    
    for i, data_point in enumerate(timeline, 1):
        timestamp = data_point["timestamp"]
        print(f"[时间点 {i}/{len(timeline)}] {timestamp.strftime('%Y-%m-%d %H:%M:%S')}")
        
        write_data_at_timestamp(spark, data_point, timestamp)
        
        # 等待一下，确保创建新的 snapshot
        if i < len(timeline):
            time.sleep(2)
            print()
    
    print("\n✓ 所有数据写入完成！")
    
    # 显示所有表的 snapshot 信息
    print("\n查看所有表的 snapshots:")
    tables = [
        "ontology.grid.substation",
        "ontology.grid.voltage_level",
        "ontology.grid.power_transformer",
        "ontology.grid.breaker",
        "ontology.assets.asset",
        "ontology.assets.work",
        "ontology.metering.meter",
        "ontology.metering.usage_point",
    ]
    
    all_snapshots = []
    for table in tables:
        try:
            df = spark.table(f"{table}.snapshots")
            snapshots = df.select("snapshot_id", "committed_at", "operation").orderBy("committed_at").collect()
            for snap in snapshots:
                all_snapshots.append({
                    "table": table,
                    "snapshot_id": snap["snapshot_id"],
                    "committed_at": snap["committed_at"],
                    "operation": snap["operation"],
                })
        except Exception as e:
            print(f"  ⚠️  无法读取 {table}.snapshots: {e}")
    
    # 按时间排序所有 snapshots
    all_snapshots.sort(key=lambda x: x["committed_at"])
    
    print(f"\n总共找到 {len(all_snapshots)} 个 snapshots（跨所有表）:\n")
    for snap in all_snapshots:
        table_short = snap["table"].split(".")[-1]
        print(f"  [{snap['committed_at']}] {table_short}: Snapshot {snap['snapshot_id']} ({snap['operation']})")
    
    # 按表分组显示
    print("\n按表分组统计:")
    from collections import defaultdict
    table_counts = defaultdict(int)
    for snap in all_snapshots:
        table_counts[snap["table"]] += 1
    
    for table, count in sorted(table_counts.items()):
        table_short = table.split(".")[-1]
        print(f"  {table_short}: {count} 个 snapshots")
    
    spark.stop()


if __name__ == "__main__":
    main()

