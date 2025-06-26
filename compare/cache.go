package compare

import (
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/ledgerwatch/erigon-lib/common"
)

type CompareBalanceCache struct {
	mu    sync.RWMutex
	cache *lru.Cache[common.Address, int]
}

func NewCompareBalanceCache() (*CompareBalanceCache, error) {
	cache, err := lru.NewWithEvict[common.Address, int](DefaultCacheSize, nil)
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

	// Only add if cache miss
	if _, ok := cache.cache.Get(address); !ok {
		cache.cache.Add(address, 0)
	}
}

func (cache *CompareBalanceCache) AddWithCount(address common.Address, count int) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	// AddWithCount overrides the current count in the current cache
	cache.cache.Add(address, count)
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

func (cache *CompareBalanceCache) GetCount(address common.Address) int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	count, _ := cache.cache.Get(address)
	return count
}

func (cache *CompareBalanceCache) GetAddresses() []common.Address {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	addresses := make([]common.Address, 0, cache.cache.Len())
	addresses = append(addresses, cache.cache.Keys()...)
	return addresses
}

type CompareAddrTokenCache struct {
	mu    sync.RWMutex
	cache *lru.Cache[common.Address, map[common.Address]int]
}

func NewCompareAddrTokenCache() (*CompareAddrTokenCache, error) {
	cache, err := lru.NewWithEvict[common.Address, map[common.Address]int](DefaultCacheSize, nil)
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
		addresses = make(map[common.Address]int)
		cache.cache.Add(tokenAddress, addresses)
	}
	// Only add if cache miss
	if _, ok := addresses[address]; !ok {
		addresses[address] = 0
	}
}

func (cache *CompareAddrTokenCache) AddWithCount(tokenAddress common.Address, address common.Address, count int) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	addresses, ok := cache.cache.Get(tokenAddress)
	if !ok {
		addresses = make(map[common.Address]int)
		cache.cache.Add(tokenAddress, addresses)
	}
	// AddWithCount overrides the current count in the current cache
	addresses[address] = count
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

func (cache *CompareAddrTokenCache) GetCount(tokenAddress common.Address, address common.Address) int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	addresses, ok := cache.cache.Get(tokenAddress)
	if ok {
		count := addresses[address]
		return count
	}
	return 0
}

func (cache *CompareAddrTokenCache) GetTokenAddresses() []common.Address {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	tokenAddresses := make([]common.Address, 0, cache.cache.Len())
	tokenAddresses = append(tokenAddresses, cache.cache.Keys()...)
	return tokenAddresses
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
