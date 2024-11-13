package bridge

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

func (brg *AppHandler) PokeNode(_ models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	return PokeNode(brg.Poke, requestDetails.Name, requestDetails.Project)
}

func PokeNode(poke bridge.PokeServiceClient, nodeID, project string) (interface{}, error) {
	md := metadata.Pairs("bigbang-controller", "1")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resp, err := poke.Poke(ctx, &bridge.PokeRequest{NodeID: nodeID, Project: project})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
