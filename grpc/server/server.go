package server

import (
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"

	/* 	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	   	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	   	extensionservice "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	   	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	   	runtimeservice "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	   	secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3" */
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	serverBridge "github.com/sefaphlvn/bigbang/grpc/server/bridge"
	"github.com/sefaphlvn/bigbang/grpc/server/snapshot"
	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/db"
)

const (
	grpcKeepaliveTime        = 30 * time.Second
	grpcKeepaliveTimeout     = 10 * time.Second
	grpcKeepaliveMinTime     = 15 * time.Second
	grpcMaxConcurrentStreams = 1000000
	grpcMaxRecvMsgSize       = 1024 * 1024 * 50 // 50MB
	grpcMaxSendMsgSize       = 1024 * 1024 * 50 // 50MB
)

type Server struct {
	xdsServer server.Server
	port      uint
	logger    *logrus.Logger
	context   *snapshot.Context
}

func NewServer(xdsServer server.Server, port uint, logger *logrus.Logger, context *snapshot.Context) *Server {
	return &Server{
		xdsServer: xdsServer,
		port:      port,
		logger:    logger,
		context:   context,
	}
}

// Run starts an xDS server at the given port.
func (s *Server) Run(db *db.AppContext, errorContext *serverBridge.ErrorContext) {
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions,
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepaliveTime,
			Timeout: grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcKeepaliveMinTime,
			PermitWithoutStream: true,
		}),
		grpc.MaxRecvMsgSize(grpcMaxRecvMsgSize),
		grpc.MaxSendMsgSize(grpcMaxSendMsgSize),
	)
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.logger.Fatal(err)
	}

	s.registerServer(grpcServer, db, errorContext)

	s.logger.Infof("Management server listening on :%d\n", s.port)
	if err = grpcServer.Serve(lis); err != nil {
		s.logger.Fatal(err)
	}
}

func (s *Server) registerServer(grpcServer *grpc.Server, db *db.AppContext, errorContext *serverBridge.ErrorContext) {
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, s.xdsServer)
	routeservice.RegisterVirtualHostDiscoveryServiceServer(grpcServer, s.xdsServer)

	// bridge grpc services
	bridge.RegisterSnapshotKeyServiceServer(grpcServer, serverBridge.NewSnapshotKeyServiceServer(s.context))
	bridge.RegisterSnapshotResourceServiceServer(grpcServer, serverBridge.NewSnapshotResourceServiceServer(s.context))
	bridge.RegisterPokeServiceServer(grpcServer, serverBridge.NewPokeServiceServer(s.context, db))

	errorService := serverBridge.NewErrorServiceServer(errorContext, s.logger)
	bridge.RegisterErrorServiceServer(grpcServer, errorService)
}
