package kafka

import "github.com/IBM/sarama"

var (
	DEFAULT_VERSION = sarama.V2_1_0_0
)

type KafkaConfig struct {
	ClientID         string
	BootstrapServers []string
	StateTopic       string
	NonStateTopic    string
}
