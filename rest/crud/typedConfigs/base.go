package typedConfigs

import (
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"github.com/sirupsen/logrus"
)

func DecodeSetTypedConfigs(resource models.DBResourceClass, logger *logrus.Logger) []*models.TypedConfig {
	var typedConfigs []*models.TypedConfig

	if paths := resource.GetGtype().TypedConfigPaths(); paths != nil {
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
			typedConfigsPart, _ := resources.ProcessTypedConfigs(jsonStringStr, path, logger)
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
