#!/usr/bin/env python3
"""
1. 创建 Breaker Actor 的 Go 代码示例
2. 从生成的 SQL 文件中提取 Breaker 表的 CREATE TABLE 语句
3. 生成可直接执行的 SQL 脚本
"""

from pathlib import Path


def extract_breaker_table_sql(sql_file: str) -> str:
    """从 SQL 文件中提取 Breaker 表的 CREATE TABLE 语句"""
    with open(sql_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 查找 breaker_snapshots 表的定义
    start_marker = "CREATE TABLE IF NOT EXISTS ontology.grid.breaker_snapshots"
    start_idx = content.find(start_marker)
    
    if start_idx == -1:
        raise ValueError("Could not find breaker_snapshots table definition")
    
    # 找到下一个 CREATE TABLE 或文件结尾
    next_create = content.find("\n\nCREATE TABLE", start_idx)
    if next_create == -1:
        # 如果没有下一个表，取到文件结尾
        sql = content[start_idx:]
    else:
        sql = content[start_idx:next_create].strip()
    
    return sql


def create_breaker_actor_example():
    """生成创建 Breaker Actor 的 Go 代码示例"""
    go_code = """package main

import (
	"context"
	"fmt"
	"log"

	"github.com/uos-projects/uos-kernel/actors"
	"github.com/uos-projects/uos-kernel/actors/capacities"
)

func main() {
	ctx := context.Background()
	
	// 创建 Actor System
	system := actors.NewSystem(ctx)
	defer system.Shutdown()
	
	// 创建 Breaker Actor（带 Control 的 PowerSystemResource）
	breakerURI := "http://www.iec.ch/TC57/CIM#Breaker"
	breakerActor := actors.NewCIMResourceActor(
		"breaker-001",
		breakerURI,
		nil, // behavior
	)
	
	// 设置 Breaker 的属性（包括继承的属性）
	breakerActor.SetProperty("mRID", "breaker-001")
	breakerActor.SetProperty("name", "Main Breaker")
	breakerActor.SetProperty("normalOpen", false)
	breakerActor.SetProperty("open", false)
	breakerActor.SetProperty("locked", false)
	breakerActor.SetProperty("ratedCurrent", 1000.0)
	breakerActor.SetProperty("description", "Main circuit breaker")
	
	// 添加 Control Capacity（Command Capacity 用于控制开关）
	commandCapacity := capacities.NewCommandCapacity("breaker-command-001")
	breakerActor.AddCapacity(commandCapacity)
	
	// 注册到 System
	if err := system.Register(breakerActor); err != nil {
		log.Fatalf("Failed to register breaker actor: %v", err)
	}
	
	// 启动 Actor
	if err := breakerActor.Start(ctx); err != nil {
		log.Fatalf("Failed to start breaker actor: %v", err)
	}
	
	fmt.Printf("Breaker Actor created successfully:\n")
	fmt.Printf("  ID: %s\n", breakerActor.ID())
	fmt.Printf("  OWL Class URI: %s\n", breakerActor.GetOWLClassURI())
	fmt.Printf("  Resource Type: %s\n", breakerActor.ResourceType())
	
	// 显示所有属性
	props := breakerActor.GetAllProperties()
	fmt.Printf("\nProperties:\n")
	for k, v := range props {
		fmt.Printf("  %s: %v\n", k, v)
	}
	
	// 创建快照（保存到 Iceberg）
	snapshot, err := breakerActor.CreateSnapshot(1)
	if err != nil {
		log.Fatalf("Failed to create snapshot: %v", err)
	}
	
	fmt.Printf("\nSnapshot created:\n")
	fmt.Printf("  Actor ID: %s\n", snapshot.ActorID)
	fmt.Printf("  OWL Class URI: %s\n", snapshot.OWLClassURI)
	fmt.Printf("  Sequence: %d\n", snapshot.Sequence)
	fmt.Printf("  Timestamp: %s\n", snapshot.Timestamp)
}
"""
    
    output_file = Path("actors/cmd/create_breaker_example/main.go")
    output_file.parent.mkdir(parents=True, exist_ok=True)
    output_file.write_text(go_code, encoding='utf-8')
    print(f"Generated Breaker Actor example: {output_file}")


def create_sql_script(breaker_sql: str):
    """生成可直接执行的 SQL 脚本"""
    sql_script = f"""-- SQL script to create Breaker snapshot table in Iceberg
-- Generated from scripts/generated_iceberg_tables.sql

-- Create namespace
CREATE NAMESPACE IF NOT EXISTS ontology.grid;

-- Create Breaker snapshot table
{breaker_sql}

-- Verify table creation
SHOW TABLES IN ontology.grid LIKE 'breaker*';
"""
    
    output_file = Path("scripts/create_breaker_table.sql")
    output_file.write_text(sql_script, encoding='utf-8')
    print(f"Generated SQL script: {output_file}")


def main():
    print("=" * 60)
    print("Create Breaker Actor and Iceberg Table")
    print("=" * 60)
    
    # 1. 生成 Breaker Actor 示例代码
    print("\n[1/3] Generating Breaker Actor example code...")
    create_breaker_actor_example()
    
    # 2. 提取 Breaker 表的 SQL
    print("\n[2/3] Extracting Breaker table SQL...")
    sql_file = "scripts/generated_iceberg_tables.sql"
    breaker_sql = extract_breaker_table_sql(sql_file)
    
    print("Extracted SQL:")
    print("-" * 60)
    print(breaker_sql)
    print("-" * 60)
    
    # 3. 生成 SQL 脚本
    print("\n[3/3] Generating SQL script...")
    create_sql_script(breaker_sql)
    
    print("\n" + "=" * 60)
    print("Completed!")
    print("=" * 60)
    print("\nNext steps:")
    print("  1. Run the Breaker Actor example:")
    print("     cd actors/cmd/create_breaker_example && go run main.go")
    print("  2. Create table in Iceberg (using Spark SQL):")
    print("     spark-sql -f scripts/create_breaker_table.sql")
    print("  3. Or execute SQL directly in Spark:")
    print("     spark.sql('CREATE NAMESPACE IF NOT EXISTS ontology.grid')")
    print("     spark.sql('''<SQL from scripts/create_breaker_table.sql>''')")


if __name__ == "__main__":
    main()

