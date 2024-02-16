package resource

import (
	"github.com/sefaphlvn/bigbang/grpc/server/resources/common"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sirupsen/logrus"
)

type AllResources struct {
	*common.Resources
}

func NewResources() *AllResources {
	return &AllResources{
		&common.Resources{},
	}
}

func SetSnapshot(rawListenerResource *models.DBResource, nodeID string, db *db.WTF, logger *logrus.Logger) (*AllResources, error) {
	resourceAll := NewResources()
	resourceAll.SetNodeID(nodeID)
	resourceAll.DecodeListener(rawListenerResource, db, logger)
	return resourceAll, nil
}
