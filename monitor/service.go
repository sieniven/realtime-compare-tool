package monitor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sieniven/realtime-compare-tool/rpc/rtclient"
)

var (
	ErrCtxCancelled = fmt.Errorf("context cancelled - stopping")
)

type MonitorService struct {
	Config MonitorConfig

	RtRpcClient *rtclient.RealtimeClient
	Logger      *log.Logger

	currHeight uint64
}

func NewMonitorService(config MonitorConfig, logger *log.Logger) (*MonitorService, error) {
	ec, err := ethclient.Dial(config.Rpc.RtRpcUrl)
	if err != nil {
		return nil, err
	}
	rtRpcClient, err := rtclient.NewRealtimeClient(context.Background(), ec, config.Rpc.RtRpcUrl)
	if err != nil {
		return nil, err
	}

	return &MonitorService{
		Config:      config,
		RtRpcClient: rtRpcClient,
		Logger:      logger,
	}, nil
}

func (service *MonitorService) Start(ctx context.Context) error {
	lastBlockUpdateTime := time.Now()
	for {
		select {
		case <-ctx.Done():
			return ErrCtxCancelled
		default:
		}

		nextHeight, err := service.RtRpcClient.RealtimeBlockNumber(ctx)
		if err != nil {
			service.Logger.Printf("error getting node height from rpc client: %v\n", err)
			continue
		}

		if service.currHeight == 0 {
			service.currHeight = nextHeight
		} else if nextHeight > service.currHeight {
			// New block height. Check time taken for block update
			timeTaken := time.Since(lastBlockUpdateTime)
			if timeTaken > service.Config.ThresholdTime {
				service.Logger.Printf("time taken for block height %d update is greater than threshold: %v\n", nextHeight, timeTaken)
			}
			lastBlockUpdateTime = time.Now()
			service.currHeight = nextHeight
		}
		time.Sleep(100 * time.Millisecond)
	}
}
