package resource

import (
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ar *AllResources) DecodeDownstreamTLS(data *models.DBResource, wtf *db.WTF) {
	dtc := &tls.DownstreamTlsContext{}
	err := resources.GetResourceWithType(data.GetResource(), dtc)
	if err != nil {
		wtf.Logger.Debug(err)
	}

	ar.AppendSecret(getValiDationContext(dtc.CommonTlsContext.GetValidationContextSdsSecretConfig().GetName(), wtf))
	ar.getTlsCertificate(dtc.CommonTlsContext.TlsCertificateSdsSecretConfigs, wtf)
}

func (ar *AllResources) getTlsCertificate(sdsSecretConfig []*tls.SdsSecretConfig, wtf *db.WTF) {
	for _, secretConf := range sdsSecretConfig {
		resource, err := resources.GetResource(wtf, "secrets", secretConf.GetName())
		if err != nil {
			wtf.Logger.Debugf("tls certificate empty resource err: %v", err)
		}

		certResources, _ := resource.Resource.Resource.(primitive.A)
		for _, certResource := range certResources {
			tlsCert := &tls.TlsCertificate{}
			err = resources.GetResourceWithType(certResource, tlsCert)
			if err != nil {
				wtf.Logger.Debugf("tls certificate decode err: %v", err)
			}

			singleResource := GetSecret(secretConf.GetName(), &tls.Secret_TlsCertificate{TlsCertificate: tlsCert})
			ar.AppendSecret(singleResource)
		}
	}
}

func getValiDationContext(vcName string, wtf *db.WTF) *tls.Secret {
	if vcName == "" {
		return nil
	}
	validationContext, err := resources.GetResource(wtf, "secrets", vcName)
	if err != nil {
		wtf.Logger.Debugf("validation context empty resource err: %v", err)
	}

	cvc := &tls.CertificateValidationContext{}
	err = resources.GetResourceWithType(validationContext.Resource.Resource, cvc)
	if err != nil {
		wtf.Logger.Debugf("validation context decode err: %v", err)
	}

	singleResource := GetSecret(vcName, &tls.Secret_ValidationContext{ValidationContext: cvc})
	return singleResource
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
