package bridge

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

func (brg *AppHandler) GetSnapshotResources(_ models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	md := metadata.Pairs("bigbang-controller", "1")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resp, err := brg.SnapshotResource.GetSnapshotResources(ctx, &bridge.SnapshotKey{Key: requestDetails.Name})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (brg *AppHandler) GetSnapshotKeys(_ models.DBResourceClass, _ models.RequestDetails) (interface{}, error) {
	md := metadata.Pairs("bigbang-controller", "1")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resp, err := brg.SnapshotKeys.GetSnapshotKeys(ctx, &bridge.Empty{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
