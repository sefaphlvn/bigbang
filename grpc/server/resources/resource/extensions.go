package resource

/*
func (ar *AllResources) CollectExtensions(resource []*models.ConfigDiscovery, db *db.AppContext, logger *logrus.Logger) {
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

func (ar *AllResources) CreateDynamicFilter(typeUrl models.GTypes, resourceName string, context *db.AppContext) (*anypb.Any, []*models.ConfigDiscovery, error) {
	switch typeUrl {
	case models.HTTPConnectionManager:
		return ar.DecodeHTTPConnectionManager(resourceName, context)
	case models.Router:
		return ar.DecodeRouter(resourceName, context)
	case models.TcpProxy:
		return ar.DecodeTcpProxy(resourceName, context)
	default:
		return nil, nil, fmt.Errorf("unknown type URL: %s", typeUrl)
	}
}
*/
