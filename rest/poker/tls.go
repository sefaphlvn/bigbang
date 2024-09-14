package poker

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

func PokerTLS(context *db.AppContext, name string, project string, gType models.GTypes, processed *Processed) {
	switch gType {
	case models.DownstreamTlsContext:
		pStreamTLS(context, name, project, processed, "listeners")
	case models.UpstreamTlsContext:
		pStreamTLS(context, name, project, processed, "clusters")
	case models.TlsCertificate:
		pTlsCertificate(context, name, project, processed)
	case models.CertificateValidationContext:
		pCertValidContext(context, name, project, processed)
	}
}

func pStreamTLS(context *db.AppContext, name string, project string, processed *Processed, collection string) {
	filter := bson.D{{Key: "general.typed_config.name", Value: name}}
	CheckResource(context, filter, collection, project, processed)
}

func pTlsCertificate(context *db.AppContext, name string, project string, processed *Processed) {
	filter := bson.D{{Key: "resource.resource.common_tls_context.tls_certificate_sds_secret_configs.name", Value: name}}
	CheckResource(context, filter, "secrets", project, processed)
}

func pCertValidContext(context *db.AppContext, name string, project string, processed *Processed) {
	filter := bson.D{{Key: "resource.resource.common_tls_context.validation_context_sds_secret_config.name", Value: name}}
	CheckResource(context, filter, "secrets", project, processed)
}
