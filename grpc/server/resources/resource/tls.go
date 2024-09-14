package resource

import (
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ar *AllResources) DecodeDownstreamTLS(data *models.DBResource, context *db.AppContext) {
	dtc := &tls.DownstreamTlsContext{}
	err := resources.MarshalUnmarshalWithType(data.GetResource(), dtc)
	if err != nil {
		context.Logger.Debug(err)
	}

	if vcName := dtc.CommonTlsContext.GetValidationContextSdsSecretConfig().GetName(); vcName != "" {
		ar.getValiDationContext(vcName, context, ar.Project)
	}

	ar.getTlsCertificate(dtc.CommonTlsContext.TlsCertificateSdsSecretConfigs, context)
}

func (ar *AllResources) getTlsCertificate(sdsSecretConfig []*tls.SdsSecretConfig, context *db.AppContext) {
	for _, secretConf := range sdsSecretConfig {
		resource, err := resources.GetResourceNGeneral(context, "secrets", secretConf.GetName(), ar.Project)
		if err != nil {
			context.Logger.Debugf("tls certificate empty resource err: %v", err)
		}

		certResources, _ := resource.Resource.Resource.(primitive.A)
		for _, certResource := range certResources {
			tlsCert := &tls.TlsCertificate{}
			err = resources.MarshalUnmarshalWithType(certResource, tlsCert)
			if err != nil {
				context.Logger.Debugf("tls certificate decode err: %v", err)
			}

			singleResource := GetSecret(secretConf.GetName(), &tls.Secret_TlsCertificate{TlsCertificate: tlsCert})
			ar.AppendSecret(singleResource)
		}
	}
}

func (ar *AllResources) getValiDationContext(vcName string, context *db.AppContext, project string) {
	validationContext, err := resources.GetResourceNGeneral(context, "secrets", vcName, project)
	if err != nil {
		context.Logger.Debugf("validation context empty resource err: %v", err)
	}

	cvc := &tls.CertificateValidationContext{}
	err = resources.MarshalUnmarshalWithType(validationContext.Resource.Resource, cvc)
	if err != nil {
		context.Logger.Debugf("validation context decode err: %v", err)
	}

	singleResource := GetSecret(vcName, &tls.Secret_ValidationContext{ValidationContext: cvc})

	ar.AppendSecret(singleResource)
}

func GetSecret(name string, typ interface{}) *tls.Secret {
	singleResource := &tls.Secret{}
	singleResource.Name = name

	switch v := typ.(type) {
	case *tls.Secret_TlsCertificate:
		singleResource.Type = &tls.Secret_TlsCertificate{TlsCertificate: v.TlsCertificate}
	case *tls.Secret_ValidationContext:
		singleResource.Type = &tls.Secret_ValidationContext{ValidationContext: v.ValidationContext}
	case *tls.Secret_SessionTicketKeys:
		singleResource.Type = &tls.Secret_SessionTicketKeys{SessionTicketKeys: v.SessionTicketKeys}
	case *tls.Secret_GenericSecret:
		singleResource.Type = &tls.Secret_GenericSecret{GenericSecret: v.GenericSecret}
	default:
		return nil
	}

	return singleResource
}
