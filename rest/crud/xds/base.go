package xds

import (
	"log"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/rest/crud"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AppHandler crud.Application

func NewXDSHandler(context *db.AppContext) *AppHandler {
	conn, err := grpc.NewClient("localhost:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	PokeClient := bridge.NewPokeServiceClient(conn)
	return &AppHandler{
		Context: context,
		Poke:    &PokeClient,
	}
}
