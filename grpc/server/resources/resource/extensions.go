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

func (ar *AllResources) CollectExtensions(resource []*models.ConfigDiscovery, db *db.WTF, logger *logrus.Logger) {
	var typedExtensionConfig types.Resource
	for _, configDiscovery := range resource {
		for _, extension := range configDiscovery.Extensions {
			anyResource, additionalResources, err := ar.CreateDynamicFilter(extension.GType, extension.Name, db)
			if err != nil {
				logger.Error(err)
			}

			typedExtensionConfig = &core.TypedExtensionConfig{
				Name:        configDiscovery.ParentName,
				TypedConfig: anyResource,
			}

			ar.Extensions = append(ar.Extensions, typedExtensionConfig)
			if additionalResources != nil {
				ar.CollectExtensions(additionalResources, db, logger)
			}
		}
	}

}

func (ar *AllResources) CreateDynamicFilter(typeUrl models.GTypes, resourceName string, wtf *db.WTF) (*anypb.Any, []*models.ConfigDiscovery, error) {
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
