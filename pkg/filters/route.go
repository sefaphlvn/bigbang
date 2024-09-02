package filters

import "go.mongodb.org/mongo-driver/bson"

func RouteDownstreamFilters(name string) MongoFilters {
	return MongoFilters{
		Collection: "extensions",
		Filter:     bson.D{{Key: "resource.resource.rds.route_config_name", Value: name}},
	}
}
