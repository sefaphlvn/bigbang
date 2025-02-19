package extension

import (
	"log"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/rest/crud"
)

type AppHandler crud.Application

func NewExtensionHandler(context *db.AppContext) *AppHandler {
	conn, err := bridge.NewGRPCClient(context)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	PokeClient := bridge.NewPokeServiceClient(conn)

	return &AppHandler{
		Context: context,
		Poke:    &PokeClient,
	}
}
