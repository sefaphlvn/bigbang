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

func SetSnapshot(rawListenerResource *models.DBResource, nodeID string, db *db.AppContext, logger *logrus.Logger) (*AllResources, error) {
	resourceAll := NewResources()
	resourceAll.SetNodeID(nodeID)
	resourceAll.DecodeListener(rawListenerResource, db, logger)
	return resourceAll, nil
}
