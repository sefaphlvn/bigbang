package dependency

import (
	"context"

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

func (h *AppHandler) GetResourceDependencies(ctx context.Context, requestDetails models.RequestDetails) (*Graph, error) {
	activeResource := Depend{
		Collection: requestDetails.Collection,
		Name:       requestDetails.Name,
		Gtype:      requestDetails.GType,
		Project:    requestDetails.Project,
		First:      true,
	}

	h.Dependencies = &Graph{}
	h.ProcessResource(ctx, activeResource)

	return h.Dependencies, nil
}

func (h *AppHandler) CallUpstreamFunction(ctx context.Context, activeResource Depend) (Node, []Depend) {
	return GenericUpstreamCollector(ctx, h, activeResource)
}

func (h *AppHandler) CallDownstreamFunction(ctx context.Context, activeResource Depend) (Node, []Depend) {
	visited := make(map[string]bool)
	return GenericDownstreamCollector(ctx, h, activeResource, visited)
}
