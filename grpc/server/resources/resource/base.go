package resource

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/sefaphlvn/bigbang/grpc/server/resources/common"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

type AllResources struct {
	*common.Resources
}

func NewResources() *AllResources {
	return &AllResources{
		&common.Resources{},
	}
}

func GenerateSnapshot(ctx context.Context, rawListenerResource *models.DBResource, listenerName string, db *db.AppContext, logger *logrus.Logger, project string) (*AllResources, error) {
	ar := NewResources()
	nodeID := fmt.Sprintf("%s:%s", listenerName, project)
	ar.SetNodeID(nodeID)
	ar.DecodeListener(ctx, rawListenerResource, db, logger)
	return ar, nil
}
