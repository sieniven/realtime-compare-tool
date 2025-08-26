package compare

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/sieniven/realtime-compare-tool/kafka"
)

type CompareConfig struct {
	Kafka kafka.KafkaConfig
	Rpc   RpcConfig

	// Compare configs
	MismatchCount     int
	CompareIntervalMS int
	SkipAddresses     []common.Address
}

type RpcConfig struct {
	RtRpcUrl    string
	NonRtRpcUrl string
	WsUrl       string
}
