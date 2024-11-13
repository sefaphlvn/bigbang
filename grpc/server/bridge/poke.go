package bridge

import (
	"context"

	"github.com/sefaphlvn/bigbang/grpc/server/resources/resource"
	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

func (s *PokeServiceServer) Poke(ctx context.Context, req *bridge.PokeRequest) (*bridge.PokeResponse, error) {
	serviceValue := req.NodeID
	projectValue := req.Project

	rawListenerResource, err := resources.GetResourceNGeneral(ctx, s.AppContext, "listeners", serviceValue, projectValue)
	if err != nil {
		return nil, err
	}

	lis, err := resource.GenerateSnapshot(ctx, rawListenerResource, serviceValue, s.AppContext, s.AppContext.Logger, projectValue)
	if err != nil {
		return nil, err
	}

	err = s.context.SetSnapshot(ctx, lis, s.AppContext.Logger)
	if err != nil {
		return nil, err
	}
	response := &bridge.PokeResponse{}

	return response, nil
}
