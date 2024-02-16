package resource

import (
	tcpProxy "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/tcp_proxy/v3"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"google.golang.org/protobuf/types/known/anypb"
)

func (ar *AllResources) DecodeTcpProxy(resourceName string, wtf *db.WTF) (*anypb.Any, []*models.AdditionalResource, error) {
	var message *anypb.Any
	resource, err := resources.GetResource(wtf, "extensions", resourceName)
	if err != nil {
		return nil, nil, err
	}

	singleResource := &tcpProxy.TcpProxy{}
	err = resources.GetResourceWithType(resource, singleResource)
	if err != nil {
		return nil, nil, err
	}

	ar.GetClustersFromClusterOrWeightedCluster(singleResource.GetClusterSpecifier(), wtf)
	message, _ = anypb.New(singleResource)

	return message, nil, nil
}

func (ar *AllResources) GetClustersFromClusterOrWeightedCluster(clusterType interface{}, wtf *db.WTF) {
	var clusters []string
	switch clusterType := clusterType.(type) {
	case *tcpProxy.TcpProxy_Cluster:
		c := clusterType.Cluster

		if c != "" {
			clusters = append(clusters, c)
		}
	case *tcpProxy.TcpProxy_WeightedClusters:
		wc := clusterType.WeightedClusters.GetClusters()

		for _, cw := range wc {
			clusters = append(clusters, cw.GetName())
		}
	}

	ar.GetClusters(clusters, wtf)
}
