package common

import (
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/sefaphlvn/bigbang/pkg/db"
)

func (ar *Resources) GetClustersFromRequestMirrorPolicies(rmps []*routev3.RouteAction_RequestMirrorPolicy, context *db.AppContext) {
	for _, rmp := range rmps {
		ar.GetClusters([]string{rmp.GetCluster()}, context)
	}
}
