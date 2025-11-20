#!/usr/bin/env python3
"""为 CIM 裁剪模型创建首批 Iceberg 表并写入示例数据"""

from __future__ import annotations

import os
from datetime import datetime
from typing import List, Dict

from pyspark.sql import SparkSession
from pyspark.sql.functions import current_timestamp, lit


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
        SparkSession.builder.appName("BootstrapIceberg")
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


def create_namespaces(spark: SparkSession) -> None:
    spark.sql("CREATE NAMESPACE IF NOT EXISTS ontology.grid")
    spark.sql("CREATE NAMESPACE IF NOT EXISTS ontology.relations")


def create_tables(spark: SparkSession) -> None:
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


def sample_entities() -> Dict[str, List[Dict[str, object]]]:
    now = datetime.utcnow()
    far_future = datetime(9999, 12, 31)
    return {
        "substation": [
            {
                "entity_id": "SS-001",
                "name": "北城变电站",
                "region": "North",
                "nominal_voltage_kv": 220.0,
                "valid_from": now,
                "valid_to": far_future,
                "op_type": "insert",
            },
            {
                "entity_id": "SS-002",
                "name": "港口变电站",
                "region": "East",
                "nominal_voltage_kv": 110.0,
                "valid_from": now,
                "valid_to": far_future,
                "op_type": "insert",
            },
        ],
        "voltage_level": [
            {
                "entity_id": "VL-001",
                "substation_id": "SS-001",
                "name": "220kV母线",
                "nominal_voltage_kv": 220.0,
                "valid_from": now,
                "valid_to": far_future,
                "op_type": "insert",
            },
            {
                "entity_id": "VL-002",
                "substation_id": "SS-002",
                "name": "110kV母线",
                "nominal_voltage_kv": 110.0,
                "valid_from": now,
                "valid_to": far_future,
                "op_type": "insert",
            },
        ],
        "relations": [
            {
                "rel_id": "REL-SS-001",
                "source_id": "SS-001",
                "target_id": "VL-001",
                "relation_type": "contains",
                "valid_from": now,
                "valid_to": far_future,
                "op_type": "insert",
            },
            {
                "rel_id": "REL-SS-002",
                "source_id": "SS-002",
                "target_id": "VL-002",
                "relation_type": "contains",
                "valid_from": now,
                "valid_to": far_future,
                "op_type": "insert",
            },
        ],
    }


def write_initial_data(spark: SparkSession) -> None:
    datasets = sample_entities()
    for table_name, records in datasets.items():
        df = spark.createDataFrame(records).withColumn("ingestion_ts", current_timestamp())
        target = (
            f"ontology.grid.{table_name}"
            if table_name in {"substation", "voltage_level"}
            else "ontology.relations.substation_voltage_level"
        )
        df.writeTo(target).append()


def main() -> None:
    spark = build_spark()
    create_namespaces(spark)
    create_tables(spark)
    write_initial_data(spark)
    spark.stop()


if __name__ == "__main__":
    main()

