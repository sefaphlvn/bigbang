package filters

import "go.mongodb.org/mongo-driver/bson"

func RouteDownstreamFilters(name string) MongoFilters {
	return MongoFilters{
		Collection: "extensions",
		Filter:     bson.D{{Key: "resource.resource.rds.route_config_name", Value: name}},
	}
}

func VirtualHostDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "routes",
			Filter:     bson.D{{Key: "general.config_discovery.name", Value: name}},
		},
		{
			Collection: "extensions",
			Filter:     bson.D{{Key: "general.config_discovery.name", Value: name}},
		},
	}
}
