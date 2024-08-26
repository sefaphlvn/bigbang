package resources

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type TempTypedConfig struct {
	Name          string `json:"name"`
	CanonicalName string `json:"canonical_name"`
	Gtype         string `json:"gtype"`
	Type          string `json:"type"`
	Category      string `json:"category"`
}

func DecodeBase64Config(encodedConfig string) (*models.TypedConfig, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedConfig)
	if err != nil {
		return nil, err
	}

	var configData models.TypedConfig
	err = json.Unmarshal(decodedBytes, &configData)
	if err != nil {
		return nil, err
	}

	return &configData, nil
}

func GetTypedConfigValue(jsonStringStr string, path string, logger *logrus.Logger) *models.TypedConfig {
	value := gjson.Get(jsonStringStr, path).String()

	if value == "" {
		logger.Debugf("typed_config value empty for path: %s", path)
		return nil
	}

	typedConfig, err := DecodeBase64Config(value)
	if err != nil {
		logger.Debugf("Error decoding base64 config: %v", err)
		return nil
	}

	return typedConfig
}

func ProcessTypedConfigs(jsonStringStr string, jsonPath string, pathTemplate string, logger *logrus.Logger) ([]*models.TypedConfig, map[string]*models.TypedConfig) {
	var typedConfigs []*models.TypedConfig
	typedConfigsMap := make(map[string]*models.TypedConfig)

	for i := range gjson.Get(jsonStringStr, jsonPath).Array() {
		path := fmt.Sprintf(pathTemplate, i)

		singleTypedConfig := GetTypedConfigValue(jsonStringStr, path+".value", logger)

		if singleTypedConfig == nil {
			continue
		}

		typedConfigs = append(typedConfigs, singleTypedConfig)
		typedConfigsMap[path] = singleTypedConfig
	}

	return typedConfigs, typedConfigsMap
}
