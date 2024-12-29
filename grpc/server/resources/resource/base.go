package resource

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/sefaphlvn/bigbang/grpc/server/resources/common"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

type AllResources struct {
	*common.Resources
	mutex sync.RWMutex
}

func NewResources() *AllResources {
	return &AllResources{
		Resources: &common.Resources{},
	}
}

func GenerateSnapshot(ctx context.Context, rawListenerResource *models.DBResource, listenerName string, db *db.AppContext, logger *logrus.Logger, project string) (*AllResources, error) {
	ar := NewResources()
	nodeID := fmt.Sprintf("%s:%s", listenerName, project)
	ar.mutex.Lock()
	ar.SetNodeID(nodeID)
	ar.mutex.Unlock()
	ar.DecodeListener(ctx, rawListenerResource, db, logger)
	return ar, nil
}
