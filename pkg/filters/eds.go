package filters

import "go.mongodb.org/mongo-driver/bson"

func EdsDownstreamFilters(name string) MongoFilters {
	return MongoFilters{
		Collection: "clusters",
		Filter:     bson.D{{Key: "resource.resource.eds_cluster_config.service_name", Value: name}},
	}
}
