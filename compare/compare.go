package compare

import (
	"context"
)

func (service *CompareService) ProcessCompareBalanceCache(ctx context.Context) {
	addresses := service.balanceCache.GetAddresses()
	for _, address := range addresses {
		// Run the native balance comparison
		ethBalance, err := service.RpcClient.EthGetBalance(address, "latest")
		if err != nil {
			service.Logger.Printf("error getting eth balance for address %s: %v\n", address, err)
			continue
		}
		realtimeBalance, err := service.RpcClient.RealtimeGetBalance(address)
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
		}
	}
}

func (service *CompareService) ProcessCompareAddrTokenCache(ctx context.Context) {
	tokenAddresses := service.addrTokenCache.GetTokenAddresses()
	for _, tokenAddress := range tokenAddresses {
		addresses := service.addrTokenCache.GetAddressesFromTokenAddress(tokenAddress)
		for _, address := range addresses {
			// Run the token balance comparison
			ethBalance, err := service.RpcClient.EthGetTokenBalance(ctx, address, tokenAddress)
			if err != nil {
				service.Logger.Printf("error getting eth token balance for token address %s and address %s: %v\n", tokenAddress, address, err)
				continue
			}
			realtimeBalance, err := service.RpcClient.RealtimeGetTokenBalance(ctx, address, tokenAddress)
			if err != nil {
				service.Logger.Printf("error getting realtime token balance for token address %s and address %s: %v\n", tokenAddress, address, err)
				continue
			}
			if ethBalance.Cmp(realtimeBalance) != 0 {
				count := service.addrTokenCache.GetCount(tokenAddress, address)
				if count > service.Config.MismatchCount {
					service.Logger.Printf("Error in state comparator: balance mismatch at height %d for address %s, eth: %s, realtime: %s\n", service.NodeHeight.Load(), address, ethBalance, realtimeBalance)
					service.balanceCache.Remove(address)
				} else {
					service.balanceCache.AddWithCount(address, count+1)
				}
			}
		}
	}
}
