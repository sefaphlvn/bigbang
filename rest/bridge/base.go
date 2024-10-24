package bridge

import (
	"log"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AppHandler struct {
	Context          *db.AppContext
	GRPCConn         *grpc.ClientConn
	SnapshotResource bridge.SnapshotResourceServiceClient
	SnapshotKeys     bridge.SnapshotKeyServiceClient
	Poke             bridge.PokeServiceClient
	Errors           bridge.ErrorServiceClient
}

func NewBridgeHandler(context *db.AppContext) *AppHandler {
	conn, err := grpc.NewClient("localhost:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	ResourcesClient := bridge.NewSnapshotResourceServiceClient(conn)
	ListenersClient := bridge.NewSnapshotKeyServiceClient(conn)
	PokeClient := bridge.NewPokeServiceClient(conn)
	ErrorsClient := bridge.NewErrorServiceClient(conn)

	return &AppHandler{
		Context:          context,
		GRPCConn:         conn,
		SnapshotResource: ResourcesClient,
		SnapshotKeys:     ListenersClient,
		Poke:             PokeClient,
		Errors:           ErrorsClient,
	}
}
