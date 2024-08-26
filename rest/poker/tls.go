package poker

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"go.mongodb.org/mongo-driver/bson"
)

func PokerTLS(context *db.AppContext, name string, gType models.GTypes, processed *Processed) {
	switch gType {
	case models.DownstreamTlsContext:
		pStreamTLS(context, name, processed)
	case models.UpstreamTlsContext:
		pStreamTLS(context, name, processed)
	case models.TlsCertificate:
		pTlsCertificate(context, name, processed)
	case models.CertificateValidationContext:
		pCertValidContext(context, name, processed)
	}
}

func pStreamTLS(context *db.AppContext, name string, processed *Processed) {
	filter := bson.D{{Key: "general.typed_config.name", Value: name}}

	rGeneral, err := resources.GetGenerals(context, "listeners", filter)
	if err != nil {
		context.Logger.Debug(err)
	}

	for _, general := range rGeneral {
		DetectChangedResource(general.GType, general.Name, context, processed)
	}
}

func pTlsCertificate(context *db.AppContext, name string, processed *Processed) {
	filter := bson.D{{Key: "resource.resource.common_tls_context.tls_certificate_sds_secret_configs.name", Value: name}}

	rGeneral, err := resources.GetGenerals(context, "secrets", filter)
	if err != nil {
		context.Logger.Debug(err)
	}

	for _, general := range rGeneral {
		DetectChangedResource(general.GType, general.Name, context, processed)
	}
}

func pCertValidContext(context *db.AppContext, name string, processed *Processed) {
	filter := bson.D{{Key: "resource.resource.common_tls_context.validation_context_sds_secret_config.name", Value: name}}

	rGeneral, err := resources.GetGenerals(context, "secrets", filter)
	if err != nil {
		context.Logger.Debug(err)
	}

	for _, general := range rGeneral {
		DetectChangedResource(general.GType, general.Name, context, processed)
	}
}
