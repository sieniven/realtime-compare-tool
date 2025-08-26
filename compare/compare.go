package compare

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func (service *CompareService) ProcessCompareBalanceCache(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
		default:
		}

		addresses := service.balanceCache.GetAddresses()
		for _, address := range addresses {
			// Run the native balance comparison
			ethBalance, err := service.NonRtRpcClient.BalanceAt(ctx, address, nil)
			if err != nil {
				service.Logger.Printf("error getting eth balance for address %s: %v\n", address, err)
				continue
			}
			realtimeBalance, err := service.RtRpcClient.RealtimeGetBalance(ctx, address)
			if err != nil {
				service.Logger.Printf("error getting realtime balance for address %s: %v\n", address, err)
				continue
			}
			if ethBalance.Cmp(realtimeBalance) != 0 {
				count := service.balanceCache.GetCount(address)
				if count > service.Config.MismatchCount {
					service.Logger.Printf("Error in state comparator: balance mismatch at height %d for address %s, eth: %s, realtime: %s\n", service.NodeHeight.Load(), address, ethBalance, realtimeBalance)
					service.balanceCache.Remove(address)
				} else {
					service.balanceCache.AddWithCount(address, count+1)
				}
			} else {
				service.Logger.Printf("Native balance are equal at height %d for address %s\n", service.NodeHeight.Load(), address)
				service.balanceCache.Remove(address)
			}
		}

		time.Sleep(time.Duration(service.Config.CompareIntervalMS) * time.Millisecond)
	}

}

func (service *CompareService) ProcessCompareAddrTokenCache(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
		default:
		}
		tokenAddresses := service.addrTokenCache.GetTokenAddresses()
		for _, tokenAddress := range tokenAddresses {
			addresses := service.addrTokenCache.GetAddressesFromTokenAddress(tokenAddress)
			for _, address := range addresses {
				// Run the token balance comparison
				ethBalance, err := service.NonRtRpcClient.EthGetTokenBalance(ctx, address, tokenAddress)
				if err != nil {
					service.Logger.Printf("error getting eth token balance for token address %s and address %s: %v\n", tokenAddress, address, err)
					continue
				}
				realtimeBalance, err := service.RtRpcClient.RealtimeGetTokenBalance(ctx, common.Address{}, address, tokenAddress)
				if err != nil {
					service.Logger.Printf("error getting realtime token balance for token address %s and address %s: %v\n", tokenAddress, address, err)
					continue
				}
				if ethBalance.Cmp(realtimeBalance) != 0 {
					count := service.addrTokenCache.GetCount(tokenAddress, address)
					if count > service.Config.MismatchCount {
						service.Logger.Printf("Error in state comparator: balance mismatch at height %d for address %s, eth: %s, realtime: %s\n", service.NodeHeight.Load(), address, ethBalance, realtimeBalance)
						service.addrTokenCache.Remove(tokenAddress, address)
					} else {
						service.addrTokenCache.AddWithCount(tokenAddress, address, count+1)
					}
				} else {
					service.Logger.Printf("Address token balances are equal at height %d for token address %s and address %s\n", service.NodeHeight.Load(), tokenAddress, address)
					service.addrTokenCache.Remove(tokenAddress, address)
				}
			}
		}

		time.Sleep(time.Duration(service.Config.CompareIntervalMS) * time.Millisecond)
	}
}
