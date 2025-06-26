package kafka

import "github.com/ledgerwatch/erigon-lib/common"

type KafkaData struct {
	Topic string                 `json:"topic"`
	Type  string                 `json:"type"`
	Data  map[string]interface{} `json:"data"`
}

type TokenHolderData struct {
	Address      common.Address
	TokenAddress common.Address
}
