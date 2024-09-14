package server

import (
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	extensionservice "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeservice "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
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
	context   *Context
}

func NewServer(xdsServer server.Server, port uint, logger *logrus.Logger, context *Context) *Server {
	return &Server{
		xdsServer: xdsServer,
		port:      port,
		logger:    logger,
		context:   context,
	}
}

// Run starts an xDS server at the given port.
func (s *Server) Run() {
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

	s.registerServer(grpcServer)

	s.logger.Infof("Management server listening on :%d\n", s.port)
	if err = grpcServer.Serve(lis); err != nil {
		s.logger.Fatal(err)
	}
}

func (s *Server) registerServer(grpcServer *grpc.Server) {
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, s.xdsServer)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, s.xdsServer)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, s.xdsServer)
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, s.xdsServer)
	routeservice.RegisterVirtualHostDiscoveryServiceServer(grpcServer, s.xdsServer)
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, s.xdsServer)
	secretservice.RegisterSecretDiscoveryServiceServer(grpcServer, s.xdsServer)
	runtimeservice.RegisterRuntimeDiscoveryServiceServer(grpcServer, s.xdsServer)
	extensionservice.RegisterExtensionConfigDiscoveryServiceServer(grpcServer, s.xdsServer)

	// custom grpc services
	/* 	snapshotStats.RegisterSnapshotKeyServiceServer(grpcServer, NewSnapshotKeyServiceServer(s.context))
	   	snapshotStats.RegisterSnapshotResourceServiceServer(grpcServer, NewSnapshotResourceServiceServer(s.context)) */
}
