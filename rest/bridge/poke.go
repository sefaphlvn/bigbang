package bridge

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
)

func PokeNode(ctx context.Context, poke bridge.PokeServiceClient, nodeID, project, version string) (interface{}, error) {
	nodeid := fmt.Sprintf("%s:%s", nodeID, project)
	md := metadata.Pairs("nodeid", nodeid)
	ctxOut := metadata.NewOutgoingContext(ctx, md)
	resp, err := poke.Poke(ctxOut, &bridge.PokeRequest{NodeID: nodeID, Project: project, Version: version})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
