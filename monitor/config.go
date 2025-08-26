package monitor

import "time"

type MonitorConfig struct {
	Rpc           RpcConfig
	ThresholdTime time.Duration
}

type RpcConfig struct {
	RtRpcUrl    string
	NonRtRpcUrl string
	WsUrl       string
}
