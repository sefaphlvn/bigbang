package poker

import (
	"fmt"

	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"go.mongodb.org/mongo-driver/bson"
)

type CdsFilters struct {
	Collection string
	Filter     bson.D
}

func CreateCdsFilters(clusterName string) []CdsFilters {
	return []CdsFilters{
		{
			Collection: "routes",
			Filter: bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "resource.resource.virtual_hosts.routes.route.cluster", Value: clusterName}},
					bson.D{{Key: "resource.resource.virtual_hosts.routes.route.weighted_clusters.clusters.name", Value: clusterName}},
					bson.D{{Key: "resource.resource.virtual_hosts.request_mirror_policies.cluster", Value: clusterName}},
					bson.D{{Key: "resource.resource.request_mirror_policies.cluster", Value: clusterName}},
				}},
			},
		}, {
			Collection: "extensions",
			Filter: bson.D{
				{Key: "$and", Value: bson.A{
					bson.D{{Key: "general.gtype", Value: "envoy.extensions.filters.network.tcp_proxy.v3.TcpProxy"}},
					bson.D{{Key: "$or", Value: bson.A{
						bson.D{{Key: "resource.resource.cluster", Value: clusterName}},
						bson.D{{Key: "resource.resource.weighted_clusters.clusters.name", Value: clusterName}},
					}}},
				}},
			},
		},
	}
}

func PokerCds(wtf *db.WTF, clusterName string) {
	cdsFilters := CreateCdsFilters(clusterName)

	for _, filter := range cdsFilters {
		resourceGeneral, err := resources.GetGenerals(wtf, filter.Collection, filter.Filter)
		if err != nil {
			wtf.Logger.Debug(err)
		}

		for _, general := range resourceGeneral {
			fmt.Println("-----------------------------", general.Name, general.GType)
			DetectChangedResource(general.GType, general.Name, wtf)
		}
	}
}
