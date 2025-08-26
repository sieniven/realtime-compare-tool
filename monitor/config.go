package monitor

import "time"

type MonitorConfig struct {
	Rpc             RpcConfig
	ThresholdTime   time.Duration
	LogBlockTimeOut bool
}

type RpcConfig struct {
	RtRpcUrl    string
	NonRtRpcUrl string
	WsUrl       string
}
