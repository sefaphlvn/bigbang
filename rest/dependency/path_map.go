package dependency

import (
	"fmt"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

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
		fmt.Println("No matching GType, returning empty paths")
		return map[string]models.GTypes{}
	}
	fmt.Printf("Matched GType: %v\n", gtype)
	return paths
}
