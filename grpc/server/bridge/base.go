package bridge

import (
	"github.com/sefaphlvn/bigbang/grpc/server/snapshot"
	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

type BaseServiceServer struct {
	context *snapshot.Context
}

// ResourceServiceServer service.
type ResourceServiceServer struct {
	bridge.UnimplementedResourceServiceServer
	*BaseServiceServer
}

func NewResourceServiceServer(context *snapshot.Context) *ResourceServiceServer {
	return &ResourceServiceServer{
		BaseServiceServer: &BaseServiceServer{context: context},
	}
}

// SnapshotServiceServer service.
type SnapshotServiceServer struct {
	bridge.UnimplementedSnapshotServiceServer
	*BaseServiceServer
}

func NewSnapshotServiceServer(context *snapshot.Context) *SnapshotServiceServer {
	return &SnapshotServiceServer{
		BaseServiceServer: &BaseServiceServer{context: context},
	}
}

// PokeServiceServer service.
type PokeServiceServer struct {
	bridge.UnimplementedPokeServiceServer
	*BaseServiceServer
	AppContext *db.AppContext
}

func NewPokeServiceServer(context *snapshot.Context, db *db.AppContext) *PokeServiceServer {
	return &PokeServiceServer{
		BaseServiceServer: &BaseServiceServer{context: context},
		AppContext:        db,
	}
}

// ActiveClientsServiceServer service.
type ActiveClientsServiceServer struct {
	bridge.UnimplementedActiveClientsServiceServer
	*BaseServiceServer
	activeClients *models.ActiveClients
}

func NewActiveClientsServiceServer(context *snapshot.Context, activeClients *models.ActiveClients) *ActiveClientsServiceServer {
	return &ActiveClientsServiceServer{
		BaseServiceServer: &BaseServiceServer{context: context},
		activeClients:     activeClients,
	}
}
