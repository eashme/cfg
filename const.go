package coord_cfg

import "time"

const (
	// yaml配置文件 环境变量键名
	YamlFileEnvKey  = "CFG_YAML"
	YamlDefaultPath = "application.yaml"

	// db 相关
	MysqlUserKey = "DB_MYSQL_USER"
	DefaultMYSQLUser = "root"
	MysqlPwdKey = "DB_MYSQL_PWD"
	DefaultMYSQLPwd = "mysql"
	MysqlHostKey = "DB_MYSQL_HOST"
	DefaultMYSQLHost = "127.0.0.1"
	MysqlPortKey = "DB_MYSQL_USER"
	DefaultMYSQLPort = "3306"
	MysqlDBKey = "DB_MYSQL_DB"
	DefaultMYSQLDB = "cfg_db"

	// redis
	RedisHostKey = "DB_REDIS_HOST"
	DefaultRedisHost = "127.0.0.1"
	RedisPortKey = "DB_REDIS_PORT"
	DefaultRedisPort = "6379"
	RedisPwdKey = "DB_REDIS_PWD"
	DefaultRedisPwd = ""
	RedisDBKey = "DB_REDIS_DB"
	DefaultRedisDB = "1"

	// 缓存相关默认值
	DefaultCacheTime = 3600*time.Second
	ScanInterval     = 30 * time.Second

	// api相关配置

	// grpc
	GrpcHostKey = "SERVER_GRPC_HOST"
	GrpcPortKey = "SERVER_GRPC_PORT"
	GrpcDefaultHost = "0.0.0.0"
	GrpcDefaultPort = 7720

	// grpc 网关
	GrpcGateWayHostKey = "SERVER_GRPC_GATEWAY_HOST"
	GrpcGateWayPortKey = "SERVER_GRPC_GATEWAY_PORT"
	GrpcGateWayDefaultHost = "0.0.0.0"
	GrpcGateWayDefaultPort = 7721
)
