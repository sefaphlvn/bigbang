package server

import (
	"fmt"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	xdsResource "github.com/sefaphlvn/bigbang/grpc/server/resources/resource"
)

func GenerateSnapshot(r *xdsResource.AllResources) *cache.Snapshot {
	fmt.Printf("GenerateSnapshotID: %v\n", r.GetVersion())
	snap, _ := cache.NewSnapshot(r.GetVersion(),
		map[resource.Type][]types.Resource{
			resource.ClusterType:         r.GetClusterT(),
			resource.RouteType:           r.GetRouteT(),
			resource.VirtualHostType:     r.GetVirtualHostT(),
			resource.EndpointType:        r.GetEndpointT(),
			resource.ListenerType:        r.GetListenerT(),
			resource.ExtensionConfigType: r.GetExtensionsT(),
			resource.SecretType:          r.GetSecretT(),
		},
	)

	//helper.PrettyPrinter(r.GetExtensionsT())
	return snap
}
