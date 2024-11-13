package dependency

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

func NewDependencyHandler(context *db.AppContext) *AppHandler {
	return &AppHandler{
		Context:      context,
		Dependencies: &Graph{},
		Cache:        make(map[string]CacheEntry),
	}
}

func (h *AppHandler) GetResourceDependencies(requestDetails models.RequestDetails) (*Graph, error) {
	activeResource := Depend{
		Collection: requestDetails.Collection,
		Name:       requestDetails.Name,
		Gtype:      requestDetails.GType,
		Project:    requestDetails.Project,
		First:      true,
	}

	h.Dependencies = &Graph{}
	h.ProcessResource(activeResource)

	return h.Dependencies, nil
}

func (h *AppHandler) CallUpstreamFunction(activeResource Depend) (Node, []Depend) {
	return GenericUpstreamCollector(h, activeResource)
}

func (h *AppHandler) CallDownstreamFunction(activeResource Depend) (Node, []Depend) {
	visited := make(map[string]bool)
	return GenericDownstreamCollector(h, activeResource, visited)
}
