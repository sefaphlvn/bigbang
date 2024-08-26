package server

import (
	"context"
	"encoding/json"

	resource "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	snapshotStats "github.com/sefaphlvn/bigbang/pkg/stats"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/anypb"
)

// SnapshotKeyServiceServer servisi
type SnapshotKeyServiceServer struct {
	snapshotStats.UnimplementedSnapshotKeyServiceServer
	context *Context
}

func NewSnapshotKeyServiceServer(context *Context) *SnapshotKeyServiceServer {
	return &SnapshotKeyServiceServer{context: context}
}

func (s *SnapshotKeyServiceServer) GetSnapshotKeys(ctx context.Context, req *snapshotStats.Empty) (*snapshotStats.SnapshotKeyList, error) {
	snapshotKeys := s.context.Cache.Cache.GetStatusKeys()
	return &snapshotStats.SnapshotKeyList{Keys: snapshotKeys}, nil
}

// SnapshotResourceServiceServer servisi
type SnapshotResourceServiceServer struct {
	snapshotStats.UnimplementedSnapshotResourceServiceServer
	context *Context
}

func NewSnapshotResourceServiceServer(context *Context) *SnapshotResourceServiceServer {
	return &SnapshotResourceServiceServer{context: context}
}

func (s *SnapshotResourceServiceServer) GetSnapshotResources(ctx context.Context, req *snapshotStats.SnapshotKey) (*snapshotStats.SnapshotResourceList, error) {
	snapshot, err := s.context.Cache.Cache.GetSnapshot(req.Key)
	status := s.context.Cache.Cache.GetStatusInfo(req.Key)
	if err != nil {
		logrus.Errorf("Error getting snapshot for key %s: %v", req.Key, err)
		return nil, err
	}

	resources := map[string]interface{}{
		"listeners":      snapshot.GetResources(resource.ListenerType),
		"clusters":       snapshot.GetResources(resource.ClusterType),
		"endpoints":      snapshot.GetResources(resource.EndpointType),
		"routes":         snapshot.GetResources(resource.RouteType),
		"secrets":        snapshot.GetResources(resource.SecretType),
		"runtimes":       snapshot.GetResources(resource.RuntimeType),
		"extensions":     snapshot.GetResources(resource.ExtensionConfigType),
		"virtualHosts":   snapshot.GetResources(resource.VirtualHostType),
		"thrift_routers": snapshot.GetResources(resource.ThriftRouteType),
		"scoped_routes":  snapshot.GetResources(resource.ScopedRouteType),
		"num_watches":    status.GetNumDeltaWatches(),
		"last_watch":     status.GetLastDeltaWatchRequestTime(),
	}
	var resourceList snapshotStats.SnapshotResourceList
	for resourceType, resourceData := range resources {
		data, err := json.Marshal(resourceData)
		if err != nil {
			return nil, err
		}
		anyData := &anypb.Any{
			TypeUrl: "type.googleapis.com/google.protobuf.Struct",
			Value:   data,
		}
		resourceList.Resources = append(resourceList.Resources, &snapshotStats.SnapshotResource{
			Type: resourceType,
			Data: anyData,
		})
	}

	return &resourceList, nil
}
