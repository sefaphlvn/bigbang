package downstreamfilters

import "go.mongodb.org/mongo-driver/bson"

func ALSDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "filters",
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

func UTMDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "virtual_host",
			Filter:     bson.D{{Key: "general.typed_config.name", Value: name}},
		},
		{
			Collection: "routes",
			Filter:     bson.D{{Key: "general.typed_config.name", Value: name}},
		},
		{
			Collection: "filters",
			Filter:     bson.D{{Key: "general.typed_config.name", Value: name}},
		},
	}
}

func CompressorLibraryDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "filters",
			Filter:     bson.D{{Key: "general.typed_config.name", Value: name}},
		},
	}
}

func TypedHTTPProtocolDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "clusters",
			Filter:     bson.D{{Key: "general.typed_config.name", Value: name}},
		},
	}
}
