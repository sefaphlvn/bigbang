package resources

import (
	"encoding/base64"
	"encoding/json"

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
