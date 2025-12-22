#!/bin/bash
# beeline wrapper script using Docker container
# Usage: ./beeline.sh [beeline options]
# Example: ./beeline.sh -u "jdbc:hive2://localhost:10001/default;transportMode=http;httpPath=cliservice" -e "SELECT 1"

docker exec -i iceberg-spark-thrift /opt/spark/bin/beeline "$@"
