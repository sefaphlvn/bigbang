package dependency

import (
	"fmt"
	"time"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

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

func (h *AppHandler) getResourceData(collection, name, project string) (string, string) {
	cacheKey := fmt.Sprintf("%s|%s|%s", collection, name, project)

	if cacheEntry, found := h.getCacheEntry(cacheKey); found {
		return cacheEntry.ID, cacheEntry.JSON
	}

	resource, err := resources.GetResourceNGeneral(h.Context, collection, name, project)
	if err != nil {
		h.Context.Logger.Debugf("Error fetching resource: %v", err)
		return "", ""
	}

	resourceID := resource.ID.Hex()
	jsonResource := resources.ConvertToJSON(resource, h.Context.Logger)

	h.setCacheEntry(cacheKey, CacheEntry{
		ID:   resourceID,
		JSON: jsonResource,
	})

	return resourceID, jsonResource
}

func getDynamicJsonPaths(gtype models.GTypes) map[string]models.GTypes {
	paths := gtype.GetUpstreamPaths()

	if len(paths) == 0 {
		//fmt.Println("No matching GType, returning empty paths")
		return map[string]models.GTypes{}
	}

	if gtype == models.VirtualHost {
		paths = addPrefixToPaths(paths, "#.")
	}

	// fmt.Printf("Matched GType: %v\n", gtype)
	return paths
}

func addPrefixToPaths(paths map[string]models.GTypes, prefix string) map[string]models.GTypes {
	newPaths := make(map[string]models.GTypes)
	for path, gtype := range paths {
		newPaths[prefix+path] = gtype
	}
	return newPaths
}
