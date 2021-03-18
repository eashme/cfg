package coord_cfg

import "time"

const (
	// yaml配置文件 环境变量键名
	YamlFileEnvKey  = "CFG_YAML"
	YamlDefaultPath = "application.yaml"

	// 缓存相关默认值
	DefaultCacheTime = 3600*time.Second
	ScanInterval     = 30 * time.Second

	// api相关配置
	GrpcHostKey = "SERVER_GRPC_HOST"
	GrpcPortKey = "SERVER_GRPC_PORT"
	GrpcDefaultHost = "0.0.0.0"
	GrpcDefaultPort = 7720

	GrpcGateWayHostKey = "SERVER_GRPC_GATEWAY_HOST"
	GrpcGateWayPortKey = "SERVER_GRPC_GATEWAY_PORT"
	GrpcGateWayDefaultHost = "0.0.0.0"
	GrpcGateWayDefaultPort = 7721
)
