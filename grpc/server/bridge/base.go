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

// SnapshotKeyServiceServer service.
type SnapshotKeyServiceServer struct {
	bridge.UnimplementedSnapshotKeyServiceServer
	*BaseServiceServer
}

func NewSnapshotKeyServiceServer(context *snapshot.Context) *SnapshotKeyServiceServer {
	return &SnapshotKeyServiceServer{
		BaseServiceServer: &BaseServiceServer{context: context},
	}
}

// SnapshotResourceServiceServer service.
type SnapshotResourceServiceServer struct {
	bridge.UnimplementedSnapshotResourceServiceServer
	*BaseServiceServer
}

func NewSnapshotResourceServiceServer(context *snapshot.Context) *SnapshotResourceServiceServer {
	return &SnapshotResourceServiceServer{
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

// ActiveClientsServiceServer service
type ActiveClientsServiceServer struct {
	bridge.UnimplementedActiveClientsServiceServer
	*BaseServiceServer
	activeClients *models.ActiveClients
}

// Yeni bir ActiveClientsServiceServer oluşturur.
func NewActiveClientsServiceServer(context *snapshot.Context, activeClients *models.ActiveClients) *ActiveClientsServiceServer {
	return &ActiveClientsServiceServer{
		BaseServiceServer: &BaseServiceServer{context: context},
		activeClients:     activeClients,
	}
}
