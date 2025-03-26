package xds

import (
	"log"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/rest/crud"
)

type AppHandler crud.Application

func NewXDSHandler(context *db.AppContext) *AppHandler {
	conn, err := bridge.NewGRPCClient(context)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	PokeClient := bridge.NewPokeServiceClient(conn)
	ResourceServiceClient := bridge.NewResourceServiceClient(conn)
	return &AppHandler{
		Context:         context,
		PokeService:            &PokeClient,
		ResourceService: &ResourceServiceClient,
	}
}
