package typed_configs

import (
	"encoding/json"
	"fmt"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func getListenerTypedConfigs(listener models.DBResourceClass, logger *logrus.Logger) []*models.TypedConfig {
	var typedConfigs []*models.TypedConfig
	resource := listener.GetResource()

	listeners, _ := resource.([]interface{})

	for _, lr := range listeners {
		jsonString, err := json.Marshal(lr)
		if err != nil {
			logger.Debugf("Error marshalling JSON: %v", err)
			return typedConfigs
		}

		jsonStringStr := string(jsonString)
		for i := range gjson.Get(jsonStringStr, "filter_chains").Array() {
			path := fmt.Sprintf(models.TransportSocketPath.String(), i)
			singleTypedConfig := resources.GetTypedConfigValue(jsonStringStr, path+".value", logger)

			if singleTypedConfig == nil {
				continue
			}

			typedConfigs = append(typedConfigs, singleTypedConfig)
		}
	}

	return typedConfigs
}
