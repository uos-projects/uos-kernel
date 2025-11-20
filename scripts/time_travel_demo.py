#!/usr/bin/env python3
"""示例：按时间点/快照回溯实体与关系"""

from __future__ import annotations

import argparse
import os
os.environ.setdefault("AWS_REGION", "us-east-1")
os.environ.setdefault("AWS_ACCESS_KEY_ID", "iceberg")
os.environ.setdefault("AWS_SECRET_ACCESS_KEY", "iceberg_password")


from pyspark.sql import SparkSession


def build_spark() -> SparkSession:
    iceberg_packages = ",".join(
        [
            "org.apache.iceberg:iceberg-spark-runtime-3.5_2.12:1.5.0",
            "org.apache.iceberg:iceberg-aws-bundle:1.5.0",
        ]
    )
    return (
        SparkSession.builder.appName("TimeTravelDemo")
        .config("spark.sql.catalog.ontology", "org.apache.iceberg.spark.SparkCatalog")
        .config("spark.jars.packages", iceberg_packages)
        .config(
            "spark.sql.catalog.ontology.catalog-impl",
            "org.apache.iceberg.nessie.NessieCatalog",
        )
        .config("spark.sql.catalog.ontology.uri", "http://localhost:19120/api/v2")
        .config("spark.sql.catalog.ontology.ref", "main")
        .config("spark.sql.catalog.ontology.warehouse", "s3a://iceberg/warehouse")
        .config("spark.sql.catalog.ontology.io-impl", "org.apache.iceberg.aws.s3.S3FileIO")
        .config("spark.sql.catalog.ontology.s3.endpoint", "http://localhost:19000")
        .config("spark.sql.catalog.ontology.s3.access-key-id", "iceberg")
        .config("spark.sql.catalog.ontology.s3.secret-access-key", "iceberg_password")
        .config("spark.sql.catalog.ontology.s3.path-style-access", "true")
        .config("spark.sql.catalog.ontology.s3.region", "us-east-1")
        .config("spark.hadoop.fs.s3a.endpoint", "http://localhost:19000")
        .config("spark.hadoop.fs.s3a.path.style.access", "true")
        .config("spark.hadoop.fs.s3a.access.key", "iceberg")
        .config("spark.hadoop.fs.s3a.secret.key", "iceberg_password")
        .config("spark.sql.extensions", "org.apache.iceberg.spark.extensions.IcebergSparkSessionExtensions")
        .getOrCreate()
    )


def query_by_snapshot(spark: SparkSession, snapshot_id: int) -> None:
    print(f"== Substation @ snapshot {snapshot_id} ==")
    spark.sql(
        f"""
        SELECT *
        FROM ontology.grid.substation VERSION AS OF {snapshot_id}
        """
    ).show(truncate=False)

    print(f"== Relations @ snapshot {snapshot_id} ==")
    spark.sql(
        f"""
        SELECT *
        FROM ontology.relations.substation_voltage_level VERSION AS OF {snapshot_id}
        """
    ).show(truncate=False)


def query_by_timestamp(spark: SparkSession, timestamp: str) -> None:
    print(f"== Substation @ {timestamp} ==")
    spark.sql(
        f"""
        SELECT *
        FROM ontology.grid.substation TIMESTAMP AS OF '{timestamp}'
        """
    ).show(truncate=False)

    spark.sql(
        f"""
        SELECT *
        FROM ontology.relations.substation_voltage_level TIMESTAMP AS OF '{timestamp}'
        """
    ).show(truncate=False)


def query_business_time(spark: SparkSession, timestamp: str) -> None:
    print(f"== Substation business view @ {timestamp} ==")
    spark.sql(
        f"""
        SELECT *
        FROM ontology.grid.substation
        WHERE '{timestamp}' BETWEEN valid_from AND valid_to
        """
    ).show(truncate=False)

    print(f"== Contains relations business view @ {timestamp} ==")
    spark.sql(
        f"""
        SELECT *
        FROM ontology.relations.substation_voltage_level
        WHERE '{timestamp}' BETWEEN valid_from AND valid_to
        """
    ).show(truncate=False)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Iceberg Time Travel Demo")
    parser.add_argument("--snapshot-id", type=int, help="指定 Iceberg snapshot id")
    parser.add_argument("--timestamp", type=str, help="技术时间（ISO8601）")
    parser.add_argument("--biz-time", type=str, help="业务时间（ISO8601）")
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    spark = build_spark()

    if args.snapshot_id is not None:
        query_by_snapshot(spark, args.snapshot_id)
    if args.timestamp is not None:
        query_by_timestamp(spark, args.timestamp)
    if args.biz_time is not None:
        query_business_time(spark, args.biz_time)

    if not any([args.snapshot_id, args.timestamp, args.biz_time]):
        print("未提供参数，默认展示当前业务视图")
        query_business_time(spark, "9999-12-31T00:00:00Z")

    spark.stop()


if __name__ == "__main__":
    main()

