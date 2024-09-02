package dependency

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

func NewDependencyHandler(context *db.AppContext) *AppHandler {
	return &AppHandler{
		Context:      context,
		Dependencies: &DependencyGraph{},
		Cache:        make(map[string]CacheEntry),
	}
}

func (h *AppHandler) GetResourceDependencies(requestDetails models.RequestDetails) (*DependencyGraph, error) {
	var activeResource = Depend{
		Collection: requestDetails.Collection,
		Name:       requestDetails.Name,
		Gtype:      requestDetails.GType,
		Project:    requestDetails.Project,
		First:      true,
	}

	h.Dependencies = &DependencyGraph{}
	h.ProcessResource(activeResource)

	return h.Dependencies, nil
}

func GenericUpstreamHandler(ctx *AppHandler, activeResource Depend) (Node, []Depend) {
	// Upstream koleksiyoncusu çağrılarak upstream bileşenler bulunur
	return GenericUpstreamCollector(ctx, activeResource)
}

func (h *AppHandler) CallUpstreamFunction(activeResource Depend) (Node, []Depend) {
	// Upstream işleyicisi çağrılarak sonuçlar döndürülür
	return GenericUpstreamHandler(h, activeResource)
}

func GenericDownstreamHandler(ctx *AppHandler, activeResource Depend) (Node, []Depend) {
	visited := make(map[string]bool)
	return GenericDownstreamCollector(ctx, activeResource, visited)
}

func (h *AppHandler) CallDownstreamFunction(activeResource Depend) (Node, []Depend) {
	// Downstream işleyicisi çağrılarak sonuçlar döndürülür
	return GenericDownstreamHandler(h, activeResource)
}
