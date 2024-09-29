package filters

import "go.mongodb.org/mongo-driver/bson"

func HcmDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "listeners",
			Filter:     bson.D{{Key: "general.config_discovery.name", Value: name}},
		},
	}
}

func RouterDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "extensions",
			Filter:     bson.D{{Key: "general.config_discovery.name", Value: name}},
		},
	}
}

func TcpProxyDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "listeners",
			Filter:     bson.D{{Key: "general.config_discovery.name", Value: name}},
		},
	}
}

func BasicAuthDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "extensions",
			Filter:     bson.D{{Key: "general.config_discovery.name", Value: name}},
		},
		{
			Collection: "routes",
			Filter:     bson.D{{Key: "general.typed_config.name", Value: name}},
		},
	}
}
