package compare

import (
	"strings"

	"github.com/ledgerwatch/erigon-lib/common"
	"github.com/sieniven/realtime-compare-tool/kafka"
	"github.com/urfave/cli/v2"
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
	RpcUrl string
	WsUrl  string
}

func NewCompareConfig(ctx *cli.Context) CompareConfig {
	cfg := CompareConfig{
		Kafka: kafka.KafkaConfig{
			BootstrapServers: strings.Split(ctx.String(KafkaBootstrapServers.Name), ","),
			StateTopic:       ctx.String(KafkaStateTopic.Name),
			NonStateTopic:    ctx.String(KafkaNonStateTopic.Name),
			ClientID:         ctx.String(KafkaClientID.Name),
		},
		Rpc: RpcConfig{
			RpcUrl: ctx.String(RpcUrl.Name),
			WsUrl:  ctx.String(WsUrl.Name),
		},
		MismatchCount:     ctx.Int(MismatchCount.Name),
		CompareIntervalMS: ctx.Int(CompareIntervalMS.Name),
		SkipAddresses:     make([]common.Address, 0),
	}

	addrsHex := strings.Split(ctx.String(SkipAddresses.Name), ",")
	for _, addrHex := range addrsHex {
		cfg.SkipAddresses = append(cfg.SkipAddresses, common.HexToAddress(addrHex))
	}

	return cfg
}
