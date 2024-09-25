package filters

import "go.mongodb.org/mongo-driver/bson"

func TlsCertificateDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "secrets",
			Filter:     bson.D{{Key: "resource.resource.common_tls_context.tls_certificate_sds_secret_configs.name", Value: name}},
		},
	}
}

func ContextValidateDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "secrets",
			Filter:     bson.D{{Key: "resource.resource.common_tls_context.validation_context_sds_secret_config.name", Value: name}},
		},
	}
}

func DownstreamTlsDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "listeners",
			Filter:     bson.D{{Key: "general.typed_config.name", Value: name}},
		},
	}
}

func UpstreamTlsDownstreamFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "clusters",
			Filter:     bson.D{{Key: "general.typed_config.name", Value: name}},
		},
	}
}
