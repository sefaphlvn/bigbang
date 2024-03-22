package typed_configs

import (
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sirupsen/logrus"
)

func DecodeSetTypedConfigs(resource models.DBResourceClass, logger *logrus.Logger) []*models.TypedConfig {
	switch resource.GetGeneral().GType {
	case models.Listener:
		return getListenerTypedConfigs(resource, logger)
	}

	return []*models.TypedConfig{}
}
