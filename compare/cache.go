package compare

import (
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/ledgerwatch/erigon-lib/common"
)

type CompareBalanceCache struct {
	mu    sync.RWMutex
	cache *lru.Cache[common.Address, struct{}]
}

func NewCompareBalanceCache() (*CompareBalanceCache, error) {
	cache, err := lru.NewWithEvict[common.Address, struct{}](DefaultCacheSize, nil)
	if err != nil {
		return nil, err
	}
	return &CompareBalanceCache{
		cache: cache,
	}, nil
}

func (cache *CompareBalanceCache) Add(address common.Address) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.cache.Add(address, struct{}{})
}

func (cache *CompareBalanceCache) Remove(address common.Address) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.cache.Remove(address)
}

func (cache *CompareBalanceCache) Clear() {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.cache.Purge()
}

func (cache *CompareBalanceCache) Size() int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	return cache.cache.Len()
}

func (cache *CompareBalanceCache) GetAddresses() []common.Address {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	addresses := make([]common.Address, 0, cache.cache.Len())
	for _, address := range cache.cache.Keys() {
		addresses = append(addresses, address)
	}
	return addresses
}

type CompareAddrTokenCache struct {
	mu    sync.RWMutex
	cache *lru.Cache[common.Address, map[common.Address]struct{}]
}

func NewCompareAddrTokenCache() (*CompareAddrTokenCache, error) {
	cache, err := lru.NewWithEvict[common.Address, map[common.Address]struct{}](DefaultCacheSize, nil)
	if err != nil {
		return nil, err
	}
	return &CompareAddrTokenCache{
		cache: cache,
	}, nil
}

func (cache *CompareAddrTokenCache) Add(tokenAddress common.Address, address common.Address) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	addresses, ok := cache.cache.Get(tokenAddress)
	if !ok {
		addresses = make(map[common.Address]struct{})
		cache.cache.Add(tokenAddress, addresses)
	}
	addresses[address] = struct{}{}
}

func (cache *CompareAddrTokenCache) Remove(tokenAddress common.Address, address common.Address) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	addresses, ok := cache.cache.Get(tokenAddress)
	if ok {
		delete(addresses, address)
		if len(addresses) == 0 {
			// Token address has no more pending addresses, remove it from the cache
			cache.cache.Remove(tokenAddress)
		}
	}
}

func (cache *CompareAddrTokenCache) Clear() {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.cache.Purge()
}

func (cache *CompareAddrTokenCache) Size() int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	return cache.cache.Len()
}

func (cache *CompareAddrTokenCache) GetAddressesFromTokenAddress(tokenAddress common.Address) []common.Address {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	addressSet, ok := cache.cache.Get(tokenAddress)
	if !ok {
		return nil
	}
	addresses := make([]common.Address, 0, len(addressSet))
	for address := range addressSet {
		addresses = append(addresses, address)
	}
	return addresses
}
