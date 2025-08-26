package main

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sieniven/realtime-compare-tool/compare"
	"github.com/sieniven/realtime-compare-tool/kafka"
	"github.com/sieniven/realtime-compare-tool/monitor"
	"github.com/urfave/cli/v2"
)

func NewCompareConfig(ctx *cli.Context) compare.CompareConfig {
	cfg := compare.CompareConfig{
		Kafka: kafka.KafkaConfig{
			BootstrapServers: strings.Split(ctx.String(KafkaBootstrapServers.Name), ","),
			StateTopic:       ctx.String(KafkaStateTopic.Name),
			NonStateTopic:    ctx.String(KafkaNonStateTopic.Name),
			ClientID:         ctx.String(KafkaClientID.Name),
		},
		Rpc: compare.RpcConfig{
			RtRpcUrl:    ctx.String(RtRpcUrl.Name),
			NonRtRpcUrl: ctx.String(NonRtRpcUrl.Name),
			WsUrl:       ctx.String(WsUrl.Name),
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

func NewMonitorConfig(ctx *cli.Context) monitor.MonitorConfig {
	cfg := monitor.MonitorConfig{
		Rpc: monitor.RpcConfig{
			RtRpcUrl:    ctx.String(RtRpcUrl.Name),
			NonRtRpcUrl: ctx.String(NonRtRpcUrl.Name),
			WsUrl:       ctx.String(WsUrl.Name),
		},
		ThresholdTime: ctx.Duration(ThresholdTime.Name),
	}
	return cfg
}
