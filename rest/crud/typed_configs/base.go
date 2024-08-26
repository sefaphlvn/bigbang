package typed_configs

import (
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"github.com/sirupsen/logrus"
)

var configGetters = map[models.GTypes][]models.TypedConfigPath{
	models.Listener:              models.ListenerTypedConfigPaths,
	models.HTTPConnectionManager: models.GeneralAccessLogTypedConfigPaths,
	models.TcpProxy:              models.GeneralAccessLogTypedConfigPaths,
	models.BootStrap:             models.BootstrapTypedConfigPaths,
}

func DecodeSetTypedConfigs(resource models.DBResourceClass, logger *logrus.Logger) []*models.TypedConfig {
	var typedConfigs []*models.TypedConfig

	if paths, exists := configGetters[resource.GetGeneral().GType]; exists {
		typedConfigs = getTypedConfigs(resource, logger, paths)
	} else {
		logger.Debugf("Unsupported general type: %v", resource.GetGeneral().GType)
	}

	return typedConfigs
}

func getTypedConfigs(resource models.DBResourceClass, logger *logrus.Logger, paths []models.TypedConfigPath) []*models.TypedConfig {
	var typedConfigs []*models.TypedConfig
	resourceValue := resource.GetResource()

	processResource := func(value interface{}) {
		jsonStringStr, err := helper.MarshalJSON(value, logger)
		if err != nil {
			logger.Errorf("Error marshaling JSON: %v", err)
			return
		}
		for _, path := range paths {
			typedConfigsPart, _ := resources.ProcessTypedConfigs(jsonStringStr, path.JsonPath, path.PathTemplate, logger)
			typedConfigs = append(typedConfigs, typedConfigsPart...)
		}
	}

	switch v := resourceValue.(type) {
	case []interface{}:
		for _, r := range v {
			processResource(r)
		}
	case interface{}:
		processResource(v)
	default:
		logger.Errorf("Unsupported resource type")
	}

	return typedConfigs
}
