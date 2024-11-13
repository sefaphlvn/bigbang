package bridge

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
)

func PokeNode(ctx context.Context, poke bridge.PokeServiceClient, nodeID, project string) (interface{}, error) {
	md := metadata.Pairs("bigbang-controller", "1")
	ctxOut := metadata.NewOutgoingContext(ctx, md)
	resp, err := poke.Poke(ctxOut, &bridge.PokeRequest{NodeID: nodeID, Project: project})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
