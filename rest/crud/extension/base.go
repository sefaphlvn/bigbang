package extension

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/rest/crud"
)

type AppHandler crud.Application

func NewExtensionHandler(context *db.AppContext) *AppHandler {
	conn, err := grpc.NewClient(
		context.Config.BigbangAddress+":"+context.Config.BigbangPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithAuthority(context.Config.BigbangAddress))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	PokeClient := bridge.NewPokeServiceClient(conn)

	return &AppHandler{
		Context: context,
		Poke:    &PokeClient,
	}
}
