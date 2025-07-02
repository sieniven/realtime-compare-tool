package compare

import "github.com/urfave/cli/v2"

var (
	// Default flags
	ConfigFlag = cli.StringFlag{
		Name:  "config",
		Usage: "Sets the configuration flags from YAML file",
		Value: "",
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
	RpcUrl = cli.StringFlag{
		Name:  "rpc.url",
		Usage: "RPC url",
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
)

var DefaultFlags = []cli.Flag{
	&ConfigFlag,
	&KafkaBootstrapServers,
	&KafkaStateTopic,
	&KafkaNonStateTopic,
	&KafkaClientID,
	&RpcUrl,
	&WsUrl,
	&MismatchCount,
	&CompareIntervalMS,
	&SkipAddresses,
}
