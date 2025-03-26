package bridge

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

func (brg *AppHandler) GetSnapshotDetails(ctx context.Context, _ models.DBResourceClass, requestDetails models.RequestDetails) (any, error) {
	md := metadata.Pairs("nodeid", requestDetails.Metadata["node_id"], "envoy-version", requestDetails.Version)
	ctxOut := metadata.NewOutgoingContext(ctx, md)
	resp, err := brg.ActiveClients.GetActiveClient(ctxOut, &bridge.NodeRequest{NodeId: requestDetails.Metadata["node_id"]})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (brg *AppHandler) GetClients(ctx context.Context, _ models.DBResourceClass, requestDetails models.RequestDetails) (any, error) {
	md := metadata.Pairs("nodeid", requestDetails.Metadata["node_id"], "envoy-version", requestDetails.Version)
	ctxOut := metadata.NewOutgoingContext(ctx, md)
	resp, err := brg.ActiveClients.GetActiveClients(ctxOut, &bridge.Empty{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
