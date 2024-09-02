package dependency

import "time"

func (h *AppHandler) getCacheEntry(key string) (CacheEntry, bool) {
	h.CacheMutex.Lock()
	defer h.CacheMutex.Unlock()

	entry, found := h.Cache[key]
	if !found {
		return CacheEntry{}, false
	}

	if time.Since(entry.Timestamp) > entry.TTL {
		delete(h.Cache, key)
		return CacheEntry{}, false
	}

	return entry, true
}

func (h *AppHandler) setCacheEntry(key string, entry CacheEntry) {
	h.CacheMutex.Lock()
	defer h.CacheMutex.Unlock()

	entry.Timestamp = time.Now()
	entry.TTL = 5 * time.Minute
	h.Cache[key] = entry
}

func (h *AppHandler) StartCacheCleanup(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			h.cleanupCache()
		}
	}()
}

func (h *AppHandler) cleanupCache() {
	h.CacheMutex.Lock()
	defer h.CacheMutex.Unlock()

	now := time.Now()
	for key, entry := range h.Cache {
		if now.Sub(entry.Timestamp) > entry.TTL {
			delete(h.Cache, key)
		}
	}
}
