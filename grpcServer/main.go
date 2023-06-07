package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
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

	db, err := db.NewMongoDB("mongodb://localhost:27017")

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	grpcserver.InitialSnapshots(db)

}

func main() {
	go func() {
		poke.Poke()
	}()

	cache := cache.NewSnapshotCache(false, cache.IDHash{}, l)
	snapshot := grpcserver.GenerateSnapshot()
	if err := snapshot.Consistent(); err != nil {
		l.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
		os.Exit(1)
	}

	l.Debugf("will serve snapshot %+v", snapshot)
	if err := cache.SetSnapshot(context.Background(), nodeID, snapshot); err != nil {
		l.Errorf("snapshot error %q for %+v", err, snapshot)
		os.Exit(1)
	}

	ctx := context.Background()
	cb := &grpcserver.Callbacks{Debug: l.Debug}
	srv := server.NewServer(ctx, cache, cb)
	grpcserver.RunServer(srv, port)
}
