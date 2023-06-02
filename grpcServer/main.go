package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v2"
	"github.com/envoyproxy/go-control-plane/pkg/server/v2"
)

var (
	l      server.Logger
	port   uint
	nodeID string
)

func init() {
	l = server.Logger{}

	flag.BoolVar(&l.Debug, "debug", true, "Enable xDS server debug logging")

	// The port that this xDS server listens on
	flag.UintVar(&port, "port", 18000, "xDS management server port")

	// Tell Envoy to use this Node ID
	flag.StringVar(&nodeID, "nodeID", "test", "Node ID")
}

func main() {
	// GRPC
	cache := cache.NewSnapshotCache(false, cache.IDHash{}, l)
	snapshot := server.GenerateSnapshot()
	if err := snapshot.Consistent(); err != nil {
		l.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
		os.Exit(1)
	}
	l.Debugf("will serve snapshot %+v", snapshot)
	if err := cache.SetSnapshot(context.Background(), nodeID, snapshot); err != nil {
		l.Errorf("snapshot error %q for %+v", err, snapshot)
		os.Exit(1)
	}

	// Run the xDS server
	ctx := context.Background()
	cb := &server.Callbacks{Debug: l.Debug}
	srv := server.NewServer(ctx, cache, cb)
	server.RunServer(srv, port)
	fmt.Println("222")
}
