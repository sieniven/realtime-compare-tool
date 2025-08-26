package main

import (
	"time"

	"github.com/urfave/cli/v2"
)

var (
	// Default flags
	ConfigFlag = cli.StringFlag{
		Name:  "config",
		Usage: "Sets the configuration flags from YAML file",
		Value: "",
	}
	// Program flags
	CompareFlag = cli.StringFlag{
		Name:  "compare",
		Usage: "Run compare program",
		Value: "false",
	}
	MonitorFlag = cli.StringFlag{
		Name:  "monitor",
		Usage: "Run monitor program",
		Value: "false",
	}
	// Kafka flags
	KafkaBootstrapServers = cli.StringFlag{
		Name:  "kafka.bootstrap-servers",
		Usage: "Kafka bootstrap servers",
		Value: "",
	}
	KafkaStateTopic = cli.StringFlag{
		Name:  "kafka.state-topic",
		Usage: "Kafka state topic",
		Value: "",
	}
	KafkaNonStateTopic = cli.StringFlag{
		Name:  "kafka.non-state-topic",
		Usage: "Kafka non state topic",
		Value: "",
	}
	KafkaClientID = cli.StringFlag{
		Name:  "kafka.client-id",
		Usage: "Kafka client id",
		Value: "",
	}
	// RPC flags
	RtRpcUrl = cli.StringFlag{
		Name:  "rpc.rt-url",
		Usage: "Realtime RPC url",
		Value: "",
	}
	NonRtRpcUrl = cli.StringFlag{
		Name:  "rpc.non-rt-url",
		Usage: "Non realtime RPC url",
		Value: "",
	}
	WsUrl = cli.StringFlag{
		Name:  "ws.url",
		Usage: "WS url",
		Value: "",
	}
	// Compare flags
	MismatchCount = cli.IntFlag{
		Name:  "compare.mismatch-count",
		Usage: "Mismatch count",
		Value: 0,
	}
	CompareIntervalMS = cli.IntFlag{
		Name:  "compare.interval-ms",
		Usage: "Compare time interval in milliseconds",
		Value: 1000,
	}
	SkipAddresses = cli.StringFlag{
		Name:  "compare.skip-addresses",
		Usage: "Skip addresses",
		Value: "",
	}
	// Monitor flags
	ThresholdTime = cli.DurationFlag{
		Name:  "monitor.threshold-time",
		Usage: "Threshold time",
		Value: 1 * time.Second,
	}
	LogBlockTimeOut = cli.BoolFlag{
		Name:  "monitor.log-block-timeout",
		Usage: "Log block timeout",
		Value: false,
	}
)

var DefaultFlags = []cli.Flag{
	&ConfigFlag,
	&CompareFlag,
	&MonitorFlag,
	&KafkaBootstrapServers,
	&KafkaStateTopic,
	&KafkaNonStateTopic,
	&KafkaClientID,
	&RtRpcUrl,
	&NonRtRpcUrl,
	&WsUrl,
	&MismatchCount,
	&CompareIntervalMS,
	&SkipAddresses,
	&ThresholdTime,
	&LogBlockTimeOut,
}
