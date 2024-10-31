package downstream_filters

import "go.mongodb.org/mongo-driver/bson"

var (
	general_config_discovery_name = "general.config_discovery.name"
	general_typed_config_name     = "general.typed_config.name"
)

func ConfigDiscoveryListenerDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "listeners",
			Filter:     bson.D{{Key: general_config_discovery_name, Value: name}},
		},
	}
}

func ConfigDiscoveryHttpFilterDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "filters",
			Filter:     bson.D{{Key: general_config_discovery_name, Value: name}},
		},
	}
}

func TypedHttpFilterDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "routes",
			Filter:     bson.D{{Key: general_typed_config_name, Value: name}},
		},
	}
}

func DiscoverAndTypedHttpFilterDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "filters",
			Filter:     bson.D{{Key: general_config_discovery_name, Value: name}},
		},
		{
			Collection: "routes",
			Filter:     bson.D{{Key: general_typed_config_name, Value: name}},
		},
	}
}
