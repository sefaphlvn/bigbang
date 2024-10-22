package downstream_filters

import "go.mongodb.org/mongo-driver/bson"

func ALSDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "extensions",
			Filter:     bson.D{{Key: "general.typed_config.name", Value: name}},
		},
		{
			Collection: "listeners",
			Filter:     bson.D{{Key: "general.typed_config.name", Value: name}},
		},
	}
}

func HCEFSDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "clusters",
			Filter:     bson.D{{Key: "general.typed_config.name", Value: name}},
		},
	}
}

func CompressorLibraryDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "extensions",
			Filter:     bson.D{{Key: "general.typed_config.name", Value: name}},
		},
	}
}

func TypedHttpProtocolDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "clusters",
			Filter:     bson.D{{Key: "general.typed_config.name", Value: name}},
		},
	}
}
