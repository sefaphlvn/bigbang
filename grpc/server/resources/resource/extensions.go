package resource

import (
	"fmt"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/anypb"
)

func (ar *AllResources) CollectExtensions(resource []*models.AdditionalResource, db *db.WTF, logger *logrus.Logger) {
	var typedExtensionConfig types.Resource
	for _, additionalResource := range resource {
		for _, extension := range additionalResource.Extensions {
			anyResource, additionalResources, err := ar.CreateDynamicFilter(extension.GType, extension.Name, db)
			if err != nil {
				logger.Error(err)
			}

			typedExtensionConfig = &core.TypedExtensionConfig{
				Name:        additionalResource.ParentName,
				TypedConfig: anyResource,
			}

			ar.Extensions = append(ar.Extensions, typedExtensionConfig)
			if additionalResources != nil {
				ar.CollectExtensions(additionalResources, db, logger)
			}
		}
	}

}

func (ar *AllResources) CreateDynamicFilter(typeUrl models.GTypes, resourceName string, wtf *db.WTF) (*anypb.Any, []*models.AdditionalResource, error) {
	switch typeUrl {
	case models.HTTPConnectionManager:
		return ar.DecodeHTTPConnectionManager(ar.Resources, resourceName, wtf)
	case models.Router:
		return ar.DecodeRouter(resourceName, wtf)
	case models.TcpProxy:
		return ar.DecodeTcpProxy(resourceName, wtf)
	default:
		return nil, nil, fmt.Errorf("unknown type URL: %s", typeUrl)
	}
}
