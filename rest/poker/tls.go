package poker

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"go.mongodb.org/mongo-driver/bson"
)

func PokerTLS(wtf *db.WTF, name string, gType models.GTypes) {
	switch gType {
	case models.DownstreamTlsContext:
		pStreamTLS(wtf, name)
	case models.UpstreamTlsContext:
		pStreamTLS(wtf, name)
	case models.TlsCertificate:
		pTlsCertificate(wtf, name)
	case models.CertificateValidationContext:
		pCertValidContext(wtf, name)
	}
}

func pStreamTLS(wtf *db.WTF, name string) {
	filter := bson.D{{Key: "general.typed_config.name", Value: name}}

	rGeneral, err := resources.GetGenerals(wtf, "listeners", filter)
	if err != nil {
		wtf.Logger.Debug(err)
	}

	for _, general := range rGeneral {
		DetectChangedResource(general.GType, general.Name, wtf)
	}
}

func pTlsCertificate(wtf *db.WTF, name string) {
	filter := bson.D{{Key: "resource.resource.common_tls_context.tls_certificate_sds_secret_configs.name", Value: name}}

	rGeneral, err := resources.GetGenerals(wtf, "secrets", filter)
	if err != nil {
		wtf.Logger.Debug(err)
	}

	for _, general := range rGeneral {
		DetectChangedResource(general.GType, general.Name, wtf)
	}
}

func pCertValidContext(wtf *db.WTF, name string) {
	filter := bson.D{{Key: "resource.resource.common_tls_context.validation_context_sds_secret_config.name", Value: name}}

	rGeneral, err := resources.GetGenerals(wtf, "secrets", filter)
	if err != nil {
		wtf.Logger.Debug(err)
	}

	for _, general := range rGeneral {
		DetectChangedResource(general.GType, general.Name, wtf)
	}
}
