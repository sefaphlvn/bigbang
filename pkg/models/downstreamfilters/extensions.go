package downstreamfilters

import "go.mongodb.org/mongo-driver/bson"

var (
	generalConfigDiscoveryName = "general.config_discovery.name"
	generalTypedConfigName     = "general.typed_config.name"
)

func ConfigDiscoveryListenerDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "listeners",
			Filter:     bson.D{{Key: generalConfigDiscoveryName, Value: name}},
		},
	}
}

func ConfigDiscoveryHTTPFilterDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "filters",
			Filter:     bson.D{{Key: generalConfigDiscoveryName, Value: name}},
		},
	}
}

func TypedHTTPFilterDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "routes",
			Filter:     bson.D{{Key: generalTypedConfigName, Value: name}},
		},
		{
			Collection: "virtual_host",
			Filter:     bson.D{{Key: generalTypedConfigName, Value: name}},
		},
		{
			Collection: "filters",
			Filter:     bson.D{{Key: generalTypedConfigName, Value: name}},
		},
	}
}

func DiscoverAndTypedHTTPFilterDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "filters",
			Filter:     bson.D{{Key: generalConfigDiscoveryName, Value: name}},
		},
		{
			Collection: "routes",
			Filter:     bson.D{{Key: generalTypedConfigName, Value: name}},
		},
		{
			Collection: "virtual_host",
			Filter:     bson.D{{Key: generalTypedConfigName, Value: name}},
		},
	}
}
