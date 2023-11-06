package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/sefaphlvn/bigbang/grpcServer/db"
	"github.com/sefaphlvn/bigbang/grpcServer/poke"
	grpcserver "github.com/sefaphlvn/bigbang/grpcServer/server"
)

var (
	l      grpcserver.Logger
	port   uint
	nodeID string
)

func init() {
	l = grpcserver.Logger{}
	flag.BoolVar(&l.Debug, "debug", true, "Enable xDS server debug logging")
	flag.UintVar(&port, "port", 18000, "xDS management server port")
	flag.StringVar(&nodeID, "nodeID", "test", "Node ID")
}

func main() {
	// connect to database
	db, err := db.NewMongoDB("mongodb+srv://navigazer:3s7qObVXRpt2wWUT@navigazer.lfh5hlh.mongodb.net")
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// create cache
	ctxCache := grpcserver.GetContext(l)
	grpcServerHandler := &grpcserver.Handler{Ctx: ctxCache, DB: db, L: &l}
	pokeHandler := &poke.Handler{Ctx: ctxCache, DB: db, L: &l, Func: grpcServerHandler}

	// start http server
	go func() {
		err := http.ListenAndServe(":8080", pokeHandler)
		if err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	// set initial snapshots
	grpcServerHandler.InitialSnapshots()
	l.Infof("all snapshots are loaded")

	// start grpc server
	ctx := context.Background()
	cb := &grpcserver.Callbacks{Debug: l.Debug}
	srv := server.NewServer(ctx, ctxCache.Cash.Cache, cb)
	grpcserver.RunServer(srv, port)
}
