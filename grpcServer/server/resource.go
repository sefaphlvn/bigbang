package server

import (
	"math/rand"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/sefaphlvn/bigbang/grpcServer/server/resources"
)

const (
	ClusterName  = "example_proxy_cluster"
	RouteName    = "firstRoute"
	ListenerPort = 10000
	UpstreamHost = "lgsagem.free.fr"
	UpstreamPort = 80
)

func makeCluster(clusterName string) *cluster.Cluster {
	return &cluster.Cluster{
		Name:                 clusterName,
		ConnectTimeout:       durationpb.New(5 * time.Second),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_LOGICAL_DNS},
		LbPolicy:             cluster.Cluster_ROUND_ROBIN,
		LoadAssignment:       makeEndpoint(clusterName),
		DnsLookupFamily:      cluster.Cluster_V4_ONLY,
	}
}

func makeEndpoint(clusterName string) *endpoint.ClusterLoadAssignment {
	return &endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			LbEndpoints: []*endpoint.LbEndpoint{{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &core.Address{
							Address: &core.Address_SocketAddress{
								SocketAddress: &core.SocketAddress{
									Protocol: core.SocketAddress_TCP,
									Address:  UpstreamHost,
									PortSpecifier: &core.SocketAddress_PortValue{
										PortValue: UpstreamPort,
									},
								},
							},
						},
					},
				},
			}},
		}},
	}
}

func makeRoute(routeName string, clusterName string) *route.RouteConfiguration {
	return &route.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []*route.VirtualHost{{
			Name:    "local_service",
			Domains: []string{"*"},
			Routes: []*route.Route{{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: "/",
					},
				},
				Action: &route.Route_Route{
					Route: &route.RouteAction{
						ClusterSpecifier: &route.RouteAction_Cluster{
							Cluster: clusterName,
						},
						HostRewriteSpecifier: &route.RouteAction_HostRewriteLiteral{
							HostRewriteLiteral: UpstreamHost,
						},
					},
				},
			}},
		}},
	}
}

func makeHTTPListener(listenerName string, route string) *listener.Listener {
	/* var listen []*listener.Listener */
	routerConfig, _ := anypb.New(&router.Router{})
	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    makeConfigSource(),
				RouteConfigName: route,
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name:       wellknown.Router,
			ConfigType: &hcm.HttpFilter_TypedConfig{TypedConfig: routerConfig},
		}},
	}
	pbst, err := anypb.New(manager)
	if err != nil {
		panic(err)
	}
	aa := rand.Uint32() + 1900
	listen := &listener.Listener{
		Name: listenerName + strconv.FormatUint(uint64(aa), 10),
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: uint32(1000 + rand.Intn(2000-1000+1)),
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}},
	}

	return listen
}

func makeConfigSource() *core.ConfigSource {
	source := &core.ConfigSource{}
	source.ResourceApiVersion = resource.DefaultAPIVersion
	source.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			TransportApiVersion:       resource.DefaultAPIVersion,
			ApiType:                   core.ApiConfigSource_GRPC,
			SetNodeOnFirstMessageOnly: true,
			GrpcServices: []*core.GrpcService{{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "xds_cluster"},
				},
			}},
		},
	}
	return source
}

func GenerateSnapshot(resources *resources.AllResources) *cache.Snapshot {
	var listeners []types.Resource
	for _, listener := range resources.Listener {
		listeners = append(listeners, listener)
	}

	snap, _ := cache.NewSnapshot(resources.Version,
		map[resource.Type][]types.Resource{
			resource.ClusterType:  {makeCluster(ClusterName)},
			resource.RouteType:    {makeRoute(RouteName, ClusterName)},
			resource.ListenerType: listeners,
		},
	)
	return snap
}

func (h *Handler) GetAllResourcesFromListener(serviceName string) (*resources.AllResources, error) {
	rawListenerResource, err := resources.GetResource(h.DB, "listeners", serviceName)
	if err != nil {
		return nil, err
	}

	lis, err := resources.SetSnapshot(rawListenerResource, serviceName, h.DB)
	if err != nil {
		return nil, err
	}

	return lis, nil
}
