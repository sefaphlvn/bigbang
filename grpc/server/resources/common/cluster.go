package common

import (
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

func (ar *Resources) GetClusters(clusters []string, context *db.AppContext) {
	for _, cls := range clusters {
		doc, _ := resources.GetResource(context, "clusters", cls)
		singleCluster := &cluster.Cluster{}
		err := resources.MarshalUnmarshalWithType(doc.GetResource(), singleCluster)
		if err != nil {
			context.Logger.Debug(err)
		}

		cc := singleCluster.GetEdsClusterConfig()
		ar.GetEndpoints(cc.ServiceName, context)
		ar.Cluster = append(ar.Cluster, singleCluster)
	}
}
