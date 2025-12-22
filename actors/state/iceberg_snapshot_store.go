package state

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// IcebergSnapshotStore 使用 Apache Iceberg 作为快照存储后端
// 通过 Spark Thrift Server（使用 Thrift 协议）与 Iceberg 表交互
type IcebergSnapshotStore struct {
	// Thrift 客户端（优先使用）
	thriftClient *ThriftClient

	// beeline 客户端（后备方案，如果 Thrift 不可用）
	jdbcURL     string
	beelinePath string

	// 命名空间（例如：ontology.grid）
	namespace string

	// 资源ID到表名的映射（缓存）
	resourceToTable map[string]string

	// 使用 Thrift 还是 beeline
	useThrift bool
}

// IcebergSnapshotStoreConfig 配置结构
type IcebergSnapshotStoreConfig struct {
	// Spark Thrift Server 主机和端口
	Host string // 例如：localhost
	Port int    // 二进制 Thrift 模式默认 10000，HTTP 模式默认 10001

	// 数据库名（默认：default）
	Database string

	// 命名空间（例如：ontology.grid）
	Namespace string

	// 用户名（可选）
	Username string

	// 密码（可选）
	Password string

	// 是否使用 Thrift 协议（默认：true，如果失败会回退到 beeline）
	UseThrift bool

	// beeline 命令路径（可选，Thrift 不可用时使用）
	BeelinePath string
}

// NewIcebergSnapshotStore 创建新的 Iceberg 快照存储
func NewIcebergSnapshotStore(config IcebergSnapshotStoreConfig) (*IcebergSnapshotStore, error) {
	host := config.Host
	if host == "" {
		host = "localhost"
	}

	port := config.Port
	if port == 0 {
		port = 10000 // Thrift 二进制模式默认端口
	}

	database := config.Database
	if database == "" {
		database = "default"
	}

	useThrift := config.UseThrift
	if !useThrift {
		useThrift = true // 默认尝试使用 Thrift
	}

	store := &IcebergSnapshotStore{
		namespace:       config.Namespace,
		resourceToTable: make(map[string]string),
		useThrift:       false, // 先设为 false，如果 Thrift 成功再设为 true
	}

	// 尝试使用 Thrift 客户端
	if useThrift {
		thriftClient, err := NewThriftClient(host, port, config.Username, config.Password)
		if err == nil {
			// Thrift 客户端创建成功
			store.thriftClient = thriftClient
			store.useThrift = true

			// 测试连接
			if err := store.ping(context.Background()); err == nil {
				return store, nil
			}
			// Thrift 连接失败，回退到 beeline
			store.thriftClient = nil
			store.useThrift = false
		}
	}

	// 回退到 beeline 客户端
	httpPort := config.Port
	if httpPort == 0 || httpPort == 10000 {
		httpPort = 10001 // HTTP 模式端口
	}

	jdbcURL := fmt.Sprintf("jdbc:hive2://%s:%d/%s;transportMode=http;httpPath=cliservice", host, httpPort, database)
	if config.Username != "" {
		jdbcURL += fmt.Sprintf(";user=%s", config.Username)
	}
	if config.Password != "" {
		jdbcURL += fmt.Sprintf(";password=%s", config.Password)
	}

	beelinePath := config.BeelinePath
	if beelinePath == "" {
		beelinePath = "beeline"
	}

	store.jdbcURL = jdbcURL
	store.beelinePath = beelinePath

	// 测试连接
	if err := store.ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to Spark Thrift Server (both Thrift and beeline failed): %w", err)
	}

	return store, nil
}

// ping 测试与 Spark Thrift Server 的连接
func (s *IcebergSnapshotStore) ping(ctx context.Context) error {
	if s.useThrift && s.thriftClient != nil {
		return s.thriftClient.ExecuteSQL(ctx, "SELECT 1")
	}
	return s.executeSQLBeeline(ctx, "SELECT 1")
}

// Save 保存快照到 Iceberg 表
func (s *IcebergSnapshotStore) Save(ctx context.Context, resourceID string, state interface{}) error {
	// 使用反射或类型断言提取字段（避免循环导入 actors 包）
	// 假设 state 是一个包含特定字段的结构体或 map
	var owlClassURI string
	var sequence int64
	var timestamp time.Time
	var properties map[string]interface{}
	var stateData interface{}

	// 尝试从 map 提取
	if snapshotMap, ok := state.(map[string]interface{}); ok {
		owlClassURI, _ = snapshotMap["OWLClassURI"].(string)
		if seq, ok := snapshotMap["Sequence"].(int64); ok {
			sequence = seq
		} else if seq, ok := snapshotMap["Sequence"].(int); ok {
			sequence = int64(seq)
		}
		if ts, ok := snapshotMap["Timestamp"].(time.Time); ok {
			timestamp = ts
		}
		properties, _ = snapshotMap["Properties"].(map[string]interface{})
		stateData = snapshotMap["State"]
	} else {
		// 尝试使用反射（更通用但性能稍差）
		// 这里简化处理，要求调用者传入 map
		return fmt.Errorf("state must be map[string]interface{} with fields: OWLClassURI, Sequence, Timestamp, Properties, State, got %T", state)
	}

	// 验证必需字段
	if owlClassURI == "" {
		return fmt.Errorf("OWLClassURI is required")
	}
	if timestamp.IsZero() {
		timestamp = time.Now()
	}
	if properties == nil {
		properties = make(map[string]interface{})
	}

	// 从 OWL URI 提取表名（例如：Breaker -> breaker_snapshots）
	tableName := s.getTableName(owlClassURI)
	fullTableName := fmt.Sprintf("%s.%s", s.namespace, tableName)

	// 缓存资源ID到表名的映射
	s.resourceToTable[resourceID] = fullTableName

	// 确保表存在
	if err := s.ensureTableExists(ctx, fullTableName, owlClassURI, properties); err != nil {
		return fmt.Errorf("failed to ensure table exists: %w", err)
	}

	// 构建 INSERT SQL
	sql := s.buildInsertSQL(fullTableName, resourceID, owlClassURI, sequence, timestamp, properties, stateData)

	// 执行 SQL
	return s.executeSQL(ctx, sql)
}

// Load 加载最新的快照
func (s *IcebergSnapshotStore) Load(ctx context.Context, resourceID string) (interface{}, error) {
	// 需要先查询表名，这里简化处理：尝试从常见的表查询
	// 实际应该维护 resourceID -> tableName 的映射
	return s.LoadAt(ctx, resourceID, time.Now().UnixMilli())
}

// LoadAt 加载指定时间点的快照（时间旅行查询）
func (s *IcebergSnapshotStore) LoadAt(ctx context.Context, resourceID string, timestamp int64) (interface{}, error) {
	// 从缓存获取表名
	fullTableName, ok := s.resourceToTable[resourceID]
	if !ok {
		// 如果缓存中没有，尝试通过元数据查询
		// 简化处理：查询所有表找到包含该 resourceID 的表
		tables, err := s.listTables(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list tables: %w", err)
		}

		// 尝试从每个表查询
		for _, tableName := range tables {
			fullTableName = fmt.Sprintf("%s.%s", s.namespace, tableName)
			snapshot, err := s.querySnapshotFromTable(ctx, fullTableName, resourceID, timestamp)
			if err == nil && snapshot != nil {
				s.resourceToTable[resourceID] = fullTableName // 缓存
				return snapshot, nil
			}
		}

		return nil, fmt.Errorf("snapshot not found for resource %s: table mapping not found", resourceID)
	}

	// 从已知的表查询
	return s.querySnapshotFromTable(ctx, fullTableName, resourceID, timestamp)
}

// querySnapshotFromTable 从指定表查询快照
func (s *IcebergSnapshotStore) querySnapshotFromTable(ctx context.Context, fullTableName, resourceID string, timestamp int64) (interface{}, error) {
	// 转换时间戳为 SQL 时间格式
	timestampTime := time.UnixMilli(timestamp)
	timestampStr := timestampTime.Format("2006-01-02 15:04:05")

	// 构建时间旅行查询 SQL
	sql := fmt.Sprintf(`
		SELECT * FROM %s
		WHERE actor_id = '%s'
		AND valid_from <= TIMESTAMP '%s'
		AND (valid_to IS NULL OR valid_to > TIMESTAMP '%s')
		ORDER BY sequence DESC
		LIMIT 1
	`, fullTableName, resourceID, timestampStr, timestampStr)

	result, err := s.querySQL(ctx, sql)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("no snapshot found")
	}

	return s.parseSnapshotFromRow(result, resourceID)
}

// listTables 列出命名空间中的所有表
func (s *IcebergSnapshotStore) listTables(ctx context.Context) ([]string, error) {
	// 简化实现：返回常见的表名列表
	// 实际应该执行 SHOW TABLES IN namespace
	return []string{"breaker_snapshots", "synchronousmachine_snapshots", "transformer_snapshots"}, nil
}

// ensureTableExists 确保表存在，如果不存在则创建
func (s *IcebergSnapshotStore) ensureTableExists(ctx context.Context, tableName string, owlClassURI string, properties map[string]interface{}) error {
	// 检查表是否存在
	checkSQL := fmt.Sprintf("SHOW TABLES IN %s LIKE '%s'", s.namespace, strings.Split(tableName, ".")[1])
	_, err := s.querySQL(ctx, checkSQL)
	if err == nil {
		// 表已存在
		return nil
	}

	// 创建命名空间（如果不存在）
	namespaceSQL := fmt.Sprintf("CREATE NAMESPACE IF NOT EXISTS %s", s.namespace)
	_ = s.executeSQL(ctx, namespaceSQL) // 忽略错误，可能已存在

	// 构建 CREATE TABLE SQL
	createSQL := s.buildCreateTableSQL(tableName)
	return s.executeSQL(ctx, createSQL)
}

// buildCreateTableSQL 构建 CREATE TABLE SQL
func (s *IcebergSnapshotStore) buildCreateTableSQL(tableName string) string {
	// 基础列
	columns := []string{
		"actor_id STRING",
		"owl_class_uri STRING",
		"sequence BIGINT",
		"timestamp TIMESTAMP",
		"valid_from TIMESTAMP",
		"valid_to TIMESTAMP",
		"op_type STRING",
		"ingestion_ts TIMESTAMP",
	}

	// 添加属性列（从 snapshot.Properties 推断）
	// 简化处理：使用 JSON 列存储所有属性
	columns = append(columns, "properties STRING") // JSON 格式

	// 添加状态列（JSON 格式）
	columns = append(columns, "state STRING") // JSON 格式

	// 构建 SQL
	sql := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			%s
		)
		USING ICEBERG
		PARTITIONED BY (actor_id)
		TBLPROPERTIES (
			'write.format.default' = 'parquet',
			'write.parquet.compression-codec' = 'zstd'
		)
	`, tableName, strings.Join(columns, ",\n\t\t"))

	return sql
}

// buildInsertSQL 构建 INSERT SQL
func (s *IcebergSnapshotStore) buildInsertSQL(tableName string, resourceID string, owlClassURI string, sequence int64, timestamp time.Time, properties map[string]interface{}, state interface{}) string {
	// 序列化属性为 JSON
	propertiesJSON, _ := json.Marshal(properties)
	propertiesStr := strings.ReplaceAll(string(propertiesJSON), "'", "''") // SQL 转义

	// 序列化状态为 JSON
	stateJSON, _ := json.Marshal(state)
	stateStr := strings.ReplaceAll(string(stateJSON), "'", "''") // SQL 转义

	// 时间戳
	timestampStr := timestamp.Format("2006-01-02 15:04:05")
	nowStr := time.Now().Format("2006-01-02 15:04:05")

	// 构建 INSERT SQL
	sql := fmt.Sprintf(`
		INSERT INTO %s (
			actor_id, owl_class_uri, sequence, timestamp,
			valid_from, valid_to, op_type, ingestion_ts,
			properties, state
		) VALUES (
			'%s', '%s', %d, TIMESTAMP '%s',
			TIMESTAMP '%s', NULL, 'INSERT', TIMESTAMP '%s',
			'%s', '%s'
		)
	`, tableName,
		resourceID,
		owlClassURI,
		sequence,
		timestampStr,
		timestampStr,
		nowStr,
		propertiesStr,
		stateStr)

	return sql
}

// executeSQL 执行 SQL 命令
func (s *IcebergSnapshotStore) executeSQL(ctx context.Context, sql string) error {
	if s.useThrift && s.thriftClient != nil {
		return s.thriftClient.ExecuteSQL(ctx, sql)
	}
	return s.executeSQLBeeline(ctx, sql)
}

// executeSQLBeeline 通过 beeline 执行 SQL（后备方案）
func (s *IcebergSnapshotStore) executeSQLBeeline(ctx context.Context, sql string) error {
	// 使用 beeline 执行 SQL
	// beeline -u "jdbc:hive2://..." -e "SQL"
	cmd := exec.CommandContext(ctx, s.beelinePath,
		"-u", s.jdbcURL,
		"-e", sql,
		"--silent=true",
		"--showHeader=false",
		"--outputformat=tsv2",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("beeline execution failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// querySQL 查询 SQL 并返回结果
func (s *IcebergSnapshotStore) querySQL(ctx context.Context, sql string) (map[string]interface{}, error) {
	if s.useThrift && s.thriftClient != nil {
		results, err := s.thriftClient.QuerySQL(ctx, sql)
		if err != nil {
			return nil, err
		}
		if len(results) > 0 {
			return results[0], nil
		}
		return nil, fmt.Errorf("no results returned")
	}
	return s.querySQLBeeline(ctx, sql)
}

// querySQLBeeline 通过 beeline 查询 SQL（后备方案）
func (s *IcebergSnapshotStore) querySQLBeeline(ctx context.Context, sql string) (map[string]interface{}, error) {
	// 使用 beeline 执行查询并返回 JSON 格式
	// beeline -u "jdbc:hive2://..." -e "SQL" --outputformat=json
	cmd := exec.CommandContext(ctx, s.beelinePath,
		"-u", s.jdbcURL,
		"-e", sql,
		"--silent=true",
		"--showHeader=false",
		"--outputformat=json",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("beeline query failed: %w\nOutput: %s", err, string(output))
	}

	// 解析 JSON 输出
	// beeline 的 JSON 输出格式可能是数组或对象
	var result map[string]interface{}

	// 尝试解析为 JSON
	var jsonData interface{}
	if err := json.Unmarshal(output, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to parse beeline output as JSON: %w\nOutput: %s", err, string(output))
	}

	// 如果输出是数组，取第一个元素
	if arr, ok := jsonData.([]interface{}); ok && len(arr) > 0 {
		if obj, ok := arr[0].(map[string]interface{}); ok {
			result = obj
		} else {
			return nil, fmt.Errorf("unexpected JSON array element type")
		}
	} else if obj, ok := jsonData.(map[string]interface{}); ok {
		result = obj
	} else {
		return nil, fmt.Errorf("unexpected JSON output format")
	}

	return result, nil
}

// parseSnapshotFromRow 从查询结果解析 Snapshot（返回 map，避免循环导入）
func (s *IcebergSnapshotStore) parseSnapshotFromRow(row map[string]interface{}, resourceID string) (interface{}, error) {
	// 解析 JSON 字段
	propertiesJSON, ok := row["properties"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid properties field")
	}

	var properties map[string]interface{}
	if err := json.Unmarshal([]byte(propertiesJSON), &properties); err != nil {
		return nil, fmt.Errorf("failed to parse properties: %w", err)
	}

	stateJSON, ok := row["state"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid state field")
	}

	var state interface{}
	if err := json.Unmarshal([]byte(stateJSON), &state); err != nil {
		return nil, fmt.Errorf("failed to parse state: %w", err)
	}

	// 解析 OWL Class URI
	owlClassURI, ok := row["owl_class_uri"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid owl_class_uri field")
	}

	// 解析 Sequence（从字符串转换为 int64）
	sequenceStr, ok := row["sequence"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid sequence field")
	}
	var sequence int64
	if _, err := fmt.Sscanf(sequenceStr, "%d", &sequence); err != nil {
		return nil, fmt.Errorf("failed to parse sequence: %w", err)
	}

	// 解析 Timestamp（从字符串转换为 time.Time）
	timestampStr, ok := row["timestamp"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid timestamp field")
	}
	timestamp, err := time.Parse("2006-01-02 15:04:05", timestampStr)
	if err != nil {
		// 尝试其他格式
		timestamp, err = time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}
	}

	// 构建 Snapshot map（避免循环导入）
	snapshot := map[string]interface{}{
		"ActorID":     resourceID,
		"OWLClassURI": owlClassURI,
		"Sequence":    sequence,
		"Timestamp":   timestamp,
		"Properties":  properties,
		"State":       state,
	}

	return snapshot, nil
}

// getTableName 从 OWL URI 提取表名
func (s *IcebergSnapshotStore) getTableName(owlURI string) string {
	// 提取本地名称（例如：Breaker）
	localName := extractLocalName(owlURI)

	// 转换为 snake_case 并添加 _snapshots 后缀
	tableName := camelToSnake(localName) + "_snapshots"
	return tableName
}

// extractLocalName 从 OWL URI 提取本地名称
func extractLocalName(uri string) string {
	for i := len(uri) - 1; i >= 0; i-- {
		if uri[i] == '#' || uri[i] == '/' {
			return uri[i+1:]
		}
	}
	return uri
}

// camelToSnake 将驼峰命名转换为蛇形命名
func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// Snapshot 实现 SnapshotStore 接口的其他方法
func (s *IcebergSnapshotStore) Snapshot() (interface{}, error) {
	return nil, fmt.Errorf("Snapshot() not supported for IcebergSnapshotStore")
}

func (s *IcebergSnapshotStore) Restore(snapshot interface{}) error {
	return fmt.Errorf("Restore() not supported for IcebergSnapshotStore")
}
