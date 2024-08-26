package resource

import (
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

func (ar *AllResources) GetRoutes(rdsName string, context *db.AppContext) error {
	route, err := resources.GetResource(context, "routes", rdsName)
	if err != nil {
		return err
	}

	singleRoute := &routev3.RouteConfiguration{}
	err = resources.MarshalUnmarshalWithType(route.GetResource(), singleRoute)
	if err != nil {
		return err
	}

	ar.AppendRoute(singleRoute)
	ar.GetClustersFromRequestMirrorPolicies(singleRoute.RequestMirrorPolicies, context)
	ar.SetClustersFromVirtualHosts(singleRoute.VirtualHosts, context)

	return nil
}

func (ar *AllResources) SetClustersFromVirtualHosts(virtualHosts []*routev3.VirtualHost, context *db.AppContext) {
	var clusters []string
	for _, vh := range virtualHosts {
		ar.GetClustersFromRequestMirrorPolicies(vh.RequestMirrorPolicies, context)
		for _, r := range vh.Routes {
			clusters = ar.GetClustersFromAction(r.GetAction(), context)
			ar.GetClusters(clusters, context)
		}
	}
}

func (ar *AllResources) GetClustersFromAction(action interface{}, db *db.AppContext) []string {
	var clusters []string
	switch action := action.(type) {
	case *routev3.Route_Route:
		wc := action.Route.GetWeightedClusters().GetClusters()
		c := action.Route.GetCluster()

		for _, cw := range wc {
			clusters = append(clusters, cw.GetName())
		}

		if c != "" {
			clusters = append(clusters, c)
		}
	}
	return clusters
}
