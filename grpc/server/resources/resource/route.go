package resource

import (
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

func (ar *AllResources) GetRoutes(rdsName string, wtf *db.WTF) error {
	route, err := resources.GetResource(wtf, "routes", rdsName)
	if err != nil {
		return err
	}

	singleRoute := &routev3.RouteConfiguration{}
	err = resources.GetResourceWithType(route.GetResource(), singleRoute)
	if err != nil {
		return err
	}

	ar.AppendRoute(singleRoute)
	ar.GetClustersFromRequestMirrorPolicies(singleRoute.RequestMirrorPolicies, wtf)
	ar.SetClustersFromVirtualHosts(singleRoute.VirtualHosts, wtf)

	return nil
}

func (ar *AllResources) SetClustersFromVirtualHosts(virtualHosts []*routev3.VirtualHost, wtf *db.WTF) {
	var clusters []string
	for _, vh := range virtualHosts {
		ar.GetClustersFromRequestMirrorPolicies(vh.RequestMirrorPolicies, wtf)
		for _, r := range vh.Routes {
			clusters = ar.GetClustersFromAction(r.GetAction(), wtf)
			ar.GetClusters(clusters, wtf)
		}
	}
}

func (ar *AllResources) GetClustersFromAction(action interface{}, db *db.WTF) []string {
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
