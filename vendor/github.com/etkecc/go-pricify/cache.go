package pricify

import "sync"

var (
	cached  *Data
	cachemu sync.Mutex
)

func getCache() *Data {
	cachemu.Lock()
	defer cachemu.Unlock()

	if cached == nil {
		return nil
	}
	return cached
}

func setCache(data *Data) {
	cachemu.Lock()
	defer cachemu.Unlock()

	cached = data
}
