package bridge

import (
	"github.com/sirupsen/logrus"

	"github.com/sefaphlvn/bigbang/grpc/server/snapshot"
	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/db"
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

// ErrorContext yapısı, hataların merkezi yönetimi için.
type ErrorServiceServer struct {
	bridge.UnimplementedErrorServiceServer
	errorContext *ErrorContext
	logger       *logrus.Logger
}

// Yeni bir ErrorServiceServer oluşturur.
func NewErrorServiceServer(ctx *ErrorContext, logger *logrus.Logger) *ErrorServiceServer {
	return &ErrorServiceServer{
		errorContext: ctx,
		logger:       logger,
	}
}
