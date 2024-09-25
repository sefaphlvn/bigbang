package bridge

import (
	"context"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"google.golang.org/grpc/metadata"
)

func (brg *AppHandler) PokeNode(resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	return PokeNode(brg.Poke, requestDetails.Name, requestDetails.Project)
}

func PokeNode(Poke bridge.PokeServiceClient, nodeID string, project string) (interface{}, error) {
	md := metadata.Pairs("bigbang-controller", "1")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resp, err := Poke.Poke(ctx, &bridge.PokeRequest{NodeID: nodeID, Project: project})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
