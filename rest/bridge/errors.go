package bridge

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

func (brg *AppHandler) GetErrors(_ models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	md := metadata.Pairs("bigbang-controller", "1")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resp, err := brg.Errors.GetNodeErrors(ctx, &bridge.NodeErrorRequest{NodeId: requestDetails.Metadata["node_id"]})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
