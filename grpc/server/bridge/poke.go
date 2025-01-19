package bridge

import (
	"context"
	"sync"

	"github.com/sefaphlvn/bigbang/grpc/server/resources/resource"
	"github.com/sefaphlvn/bigbang/grpc/server/snapshot"
	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

type PokeService struct {
	Snapshot   *snapshot.Context
	appContext *db.AppContext
	mu         sync.Mutex
}

func NewPokeService(ctxCache *snapshot.Context, appContext *db.AppContext) *PokeService {
	return &PokeService{
		Snapshot:   ctxCache,
		appContext: appContext,
	}
}

func (pss *PokeServiceServer) Poke(ctx context.Context, req *bridge.PokeRequest) (*bridge.PokeResponse, error) {
	rawListenerResource, err := resources.GetResourceNGeneral(ctx, pss.AppContext, "listeners", req.NodeID, req.Project, req.Version)
	if err != nil {
		return nil, err
	}

	lis, err := resource.GenerateSnapshot(ctx, rawListenerResource, req.NodeID, pss.AppContext, pss.AppContext.Logger, req.Project, req.Version)
	if err != nil {
		return nil, err
	}

	err = pss.context.SetSnapshot(ctx, lis, pss.AppContext.Logger)
	if err != nil {
		return nil, err
	}
	response := &bridge.PokeResponse{}

	return response, nil
}

func (ps *PokeService) CheckSnapshot(node string) bool {
	snapshot, err := ps.Snapshot.Cache.Cache.GetSnapshot(node)
	if err != nil {
		ps.appContext.Logger.Debugf("Error while fetching snapshot for node %s: %v", node, err)
		return true
	}

	if snapshot == nil {
		ps.appContext.Logger.Debugf("Snapshot is nil for node: %s", node)
		return true
	}

	ps.appContext.Logger.Debugf("Snapshot exists for node: %s", node)
	return false
}

func (ps *PokeService) getAllResourcesFromListener(ctx context.Context, listenerName, project, version string) (*resource.AllResources, error) {
	rawListenerResource, err := resources.GetResourceNGeneral(ctx, ps.appContext, "listeners", listenerName, project, version)
	if err != nil {
		return nil, err
	}

	lis, err := resource.GenerateSnapshot(ctx, rawListenerResource, listenerName, ps.appContext, ps.appContext.Logger, project, version)
	if err != nil {
		return nil, err
	}

	return lis, nil
}

func (ps *PokeService) GetResourceSetSnapshot(ctx context.Context, node, project, version string) error {
	allResource, err := ps.getAllResourcesFromListener(ctx, node, project, version)
	if err != nil {
		ps.appContext.Logger.Warnf("get resources err (%v:%v): %v", node, project, err)
		return err
	}

	err = ps.Snapshot.SetSnapshot(ctx, allResource, ps.appContext.Logger)
	if err != nil {
		ps.appContext.Logger.Warnf("%s", err)
	}
	return nil
}
