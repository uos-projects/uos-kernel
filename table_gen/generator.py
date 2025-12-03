"""生成 Iceberg CREATE TABLE SQL 语句"""

from __future__ import annotations

from typing import List, Optional

from mapping.mapper import TableSchema, Column


class SQLGenerator:
    """SQL 生成器"""
    
    @staticmethod
    def generate_create_namespace_sql(namespace: str) -> str:
        """生成创建命名空间的 SQL"""
        full_namespace = f"ontology.{namespace}"
        return f"CREATE NAMESPACE IF NOT EXISTS {full_namespace}"
    
    @staticmethod
    def generate_create_table_sql(schema: TableSchema) -> str:
        """生成 CREATE TABLE SQL"""
        full_table_name = schema.get_full_table_name()
        
        # 生成列定义
        column_defs = []
        for col in schema.columns:
            nullable = "" if col.nullable else " NOT NULL"
            comment = f" COMMENT '{col.comment}'" if col.comment else ""
            column_defs.append(f"    {col.name} {col.type}{nullable}{comment}")
        
        columns_sql = ",\n".join(column_defs)
        
        # 生成分区子句
        partition_clause = ""
        if schema.partition_fields:
            partition_fields_str = ", ".join(schema.partition_fields)
            partition_clause = f"\nPARTITIONED BY ({partition_fields_str})"
        
        # 生成注释
        comment_clause = ""
        if schema.comment:
            comment_clause = f"\nCOMMENT '{schema.comment}'"
        
        sql = f"""CREATE TABLE IF NOT EXISTS {full_table_name} (
{columns_sql}
)
USING ICEBERG{partition_clause}{comment_clause}"""
        
        return sql
    
    @staticmethod
    def generate_drop_table_sql(schema: TableSchema) -> str:
        """生成 DROP TABLE SQL"""
        full_table_name = schema.get_full_table_name()
        return f"DROP TABLE IF EXISTS {full_table_name}"
    
    @staticmethod
    def generate_table_sqls(schema: TableSchema) -> List[str]:
        """生成完整的表创建 SQL 列表（包括命名空间）"""
        sqls = []
        
        # 1. 创建命名空间
        namespace_sql = SQLGenerator.generate_create_namespace_sql(schema.namespace)
        sqls.append(namespace_sql)
        
        # 2. 创建表
        table_sql = SQLGenerator.generate_create_table_sql(schema)
        sqls.append(table_sql)
        
        return sqls


if __name__ == "__main__":
    # 测试 SQL 生成
    from mapping.mapper import EntityMapper
    
    mapper = EntityMapper()
    schema = mapper.map_entity_by_name("Substation")
    
    if schema:
        print("=== 生成的 SQL ===\n")
        sqls = SQLGenerator.generate_table_sqls(schema)
        for sql in sqls:
            print(sql)
            print()

