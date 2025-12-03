#!/usr/bin/env python3
"""
查看 MinIO 存储桶内容的脚本
"""

import boto3
from botocore.client import Config
from botocore.exceptions import ClientError
from datetime import datetime


def format_size(size_bytes):
    """格式化文件大小"""
    for unit in ['B', 'KB', 'MB', 'GB', 'TB']:
        if size_bytes < 1024.0:
            return f"{size_bytes:.2f} {unit}"
        size_bytes /= 1024.0
    return f"{size_bytes:.2f} PB"


def list_bucket_contents(s3_client, bucket_name, prefix=""):
    """列出存储桶中的对象"""
    try:
        paginator = s3_client.get_paginator('list_objects_v2')
        pages = paginator.paginate(Bucket=bucket_name, Prefix=prefix)
        
        total_size = 0
        total_files = 0
        
        print(f"\n存储桶: {bucket_name}")
        print(f"前缀: {prefix if prefix else '(全部)'}")
        print("=" * 80)
        
        for page in pages:
            if 'Contents' not in page:
                continue
                
            for obj in page['Contents']:
                total_files += 1
                size = obj['Size']
                total_size += size
                
                # 格式化时间
                last_modified = obj['LastModified'].strftime('%Y-%m-%d %H:%M:%S')
                
                print(f"{obj['Key']:<60} {format_size(size):>12}  {last_modified}")
        
        print("=" * 80)
        print(f"总计: {total_files} 个文件, 总大小: {format_size(total_size)}")
        
    except ClientError as e:
        print(f"错误: {e}")


def list_buckets(s3_client):
    """列出所有存储桶"""
    try:
        response = s3_client.list_buckets()
        print("\n可用的存储桶:")
        print("=" * 80)
        for bucket in response['Buckets']:
            creation_date = bucket['CreationDate'].strftime('%Y-%m-%d %H:%M:%S')
            print(f"  {bucket['Name']:<30} 创建时间: {creation_date}")
        print("=" * 80)
        return [b['Name'] for b in response['Buckets']]
    except ClientError as e:
        print(f"错误: {e}")
        return []


def main():
    # MinIO 配置
    endpoint_url = "http://localhost:19000"
    access_key = "iceberg"
    secret_key = "iceberg_password"
    region = "us-east-1"
    
    # 创建 S3 客户端
    s3_client = boto3.client(
        's3',
        endpoint_url=endpoint_url,
        aws_access_key_id=access_key,
        aws_secret_access_key=secret_key,
        region_name=region,
        config=Config(signature_version='s3v4')
    )
    
    print("MinIO 内容查看器")
    print("=" * 80)
    
    # 列出所有存储桶
    buckets = list_buckets(s3_client)
    
    if not buckets:
        print("\n没有找到任何存储桶")
        return
    
    # 列出每个存储桶的内容
    for bucket_name in buckets:
        print(f"\n{'=' * 80}")
        list_bucket_contents(s3_client, bucket_name)
        
        # 如果是 iceberg 存储桶，也列出子目录结构
        if bucket_name == "iceberg":
            print(f"\n{'=' * 80}")
            print("存储桶结构概览:")
            print("=" * 80)
            
            # 列出顶层目录
            prefixes = set()
            paginator = s3_client.get_paginator('list_objects_v2')
            pages = paginator.paginate(Bucket=bucket_name, Delimiter='/')
            
            for page in pages:
                if 'CommonPrefixes' in page:
                    for prefix_info in page['CommonPrefixes']:
                        prefix = prefix_info['Prefix']
                        prefixes.add(prefix)
                        print(f"  📁 {prefix}")
                
                if 'Contents' in page:
                    for obj in page['Contents']:
                        if not obj['Key'].endswith('/'):
                            # 提取目录路径
                            parts = obj['Key'].split('/')
                            if len(parts) > 1:
                                dir_prefix = '/'.join(parts[:-1]) + '/'
                                if dir_prefix not in prefixes:
                                    prefixes.add(dir_prefix)
                                    print(f"  📁 {dir_prefix}")


if __name__ == "__main__":
    main()


