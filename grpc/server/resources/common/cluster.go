package common

import (
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

func (ar *Resources) GetClusters(clusters []string, wtf *db.WTF) {
	for _, cls := range clusters {
		doc, _ := resources.GetResource(wtf, "clusters", cls)
		singleCluster := &cluster.Cluster{}
		err := resources.GetResourceWithType(doc.GetResource(), singleCluster)
		if err != nil {
			wtf.Logger.Debug(err)
		}

		cc := singleCluster.GetEdsClusterConfig()
		ar.GetEndpoints(cc.ServiceName, wtf)
		ar.Cluster = append(ar.Cluster, singleCluster)
	}
}
