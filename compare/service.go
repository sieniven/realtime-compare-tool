package compare

import (
	"context"
	"log"
	"sync/atomic"

	"github.com/ledgerwatch/erigon-lib/common"
	"github.com/sieniven/realtime-compare-tool/kafka"
	"github.com/sieniven/realtime-compare-tool/rpc"
)

type CompareService struct {
	InitFlag   atomic.Bool
	NodeHeight atomic.Int64
	Config     CompareConfig

	KafkaConsumer *kafka.KafkaConsumer
	RpcClient     *rpc.RealtimeClient
	Logger        *log.Logger

	// Compare cache
	balanceCache   *CompareBalanceCache
	addrTokenCache *CompareAddrTokenCache

	// Channels
	HeightChan      chan int64
	AddrBalanceChan chan common.Address
	TokenHolderChan chan kafka.TokenHolderData
	ErrorChan       chan error
}

func NewCompareService(config CompareConfig, logger *log.Logger) (*CompareService, error) {
	kafkaConsumer, err := kafka.NewKafkaConsumer(config.Kafka)
	if err != nil {
		return nil, err
	}
	rpcClient, err := rpc.NewRealtimeClient(config.Rpc.RpcUrl)
	if err != nil {
		return nil, err
	}
	balanceCache, err := NewCompareBalanceCache()
	if err != nil {
		return nil, err
	}
	addrTokenCache, err := NewCompareAddrTokenCache()
	if err != nil {
		return nil, err
	}

	return &CompareService{
		InitFlag:        atomic.Bool{},
		NodeHeight:      atomic.Int64{},
		Config:          config,
		KafkaConsumer:   kafkaConsumer,
		RpcClient:       rpcClient,
		Logger:          logger,
		balanceCache:    balanceCache,
		addrTokenCache:  addrTokenCache,
		HeightChan:      make(chan int64, DefaultChannelSize),
		AddrBalanceChan: make(chan common.Address, DefaultChannelSize),
		TokenHolderChan: make(chan kafka.TokenHolderData, DefaultChannelSize),
		ErrorChan:       make(chan error, DefaultChannelSize),
	}, nil
}

func (service *CompareService) Start(ctx context.Context) error {
	// Start the kafka consumer goroutine
	go service.KafkaConsumer.ConsumeKafka(ctx, service.HeightChan, service.AddrBalanceChan, service.TokenHolderChan, service.ErrorChan, service.Logger)

	for {
		select {
		case <-ctx.Done():
			return ErrCtxCancelled
		case height := <-service.HeightChan:
			if service.NodeHeight.Load() < height {
				service.NodeHeight.Store(height)
				if !service.InitFlag.Load() {
					// Try to init compare service
					ethHeight, err := service.RpcClient.EthGetBlockNumber(ctx)
					if err != nil {
						service.Logger.Printf("error getting node height from rpc client: %v\n", err)
						continue
					}
					diff := int64(ethHeight) - height
					if diff < 0 {
						diff = -diff
					}
					if diff < DefaultHeightSyncRange {
						service.InitFlag.Store(true)
						service.Logger.Println("node heights initialized, starting compare")
					}
				} else {
					// New height. Try compare states
					go service.ProcessCompareBalanceCache(ctx)
					go service.ProcessCompareAddrTokenCache(ctx)
				}
			}
		case address := <-service.AddrBalanceChan:
			if !service.InitFlag.Load() {
				continue
			}
			skipFlag := false
			for _, skipAddress := range service.Config.SkipAddresses {
				if address.Hex() == skipAddress.Hex() {
					skipFlag = true
					break
				}
			}
			if skipFlag {
				continue
			}
			service.balanceCache.Add(address)
		case tokenHolder := <-service.TokenHolderChan:
			if !service.InitFlag.Load() {
				continue
			}
			skipFlag := false
			for _, skipAddress := range service.Config.SkipAddresses {
				if tokenHolder.Address.Hex() == skipAddress.Hex() {
					skipFlag = true
					break
				}
			}
			if skipFlag {
				continue
			}
			service.addrTokenCache.Add(tokenHolder.TokenAddress, tokenHolder.Address)
		case err := <-service.ErrorChan:
			return err
		}
	}
}
