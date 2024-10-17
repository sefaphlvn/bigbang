package downstream_filters

import (
	"go.mongodb.org/mongo-driver/bson"
)

var (
	resource_resource_request_mirror_policies_cluster = "resource.resource.request_mirror_policies.cluster"
)

func ClusterDownstreamFilters(clusterName string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "routes",
			Filter: bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "resource.resource.virtual_hosts.routes.route.cluster", Value: clusterName}},
					bson.D{{Key: "resource.resource.virtual_hosts.routes.route.weighted_clusters.clusters.name", Value: clusterName}},
					bson.D{{Key: "resource.resource.virtual_hosts.request_mirror_policies.cluster", Value: clusterName}},
					bson.D{{Key: resource_resource_request_mirror_policies_cluster, Value: clusterName}},
				}},
			},
		},
		{
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
		{
			Collection: "others",
			Filter: bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "resource.resource.cluster", Value: clusterName}},
				}},
			},
		},
		{
			Collection: "virtual_host",
			Filter: bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "resource.resource.routes.route.cluster", Value: clusterName}},
					bson.D{{Key: "resource.resource.routes.route.weighted_clusters.clusters.name", Value: clusterName}},
					bson.D{{Key: resource_resource_request_mirror_policies_cluster, Value: clusterName}},
					bson.D{{Key: resource_resource_request_mirror_policies_cluster, Value: clusterName}},
				}},
			},
		},
	}
}
