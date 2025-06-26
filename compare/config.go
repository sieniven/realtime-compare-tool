package compare

import "github.com/sieniven/realtime-compare-tool/kafka"

type CompareConfig struct {
	Kafka kafka.KafkaConfig
	Rpc   RpcConfig

	// Compare configs
	MismatchCount int
}

type RpcConfig struct {
	RpcUrl string
	WsUrl  string
}
