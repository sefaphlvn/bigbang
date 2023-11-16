package server

import (
	"context"
	"github.com/sefaphlvn/bigbang/grpc/server/resources"
	"github.com/sirupsen/logrus"
	"sync"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
)

var (
	once sync.Once
	ctx  *Context
)

type Cash struct {
	Cache cache.SnapshotCache
}

type Context struct {
	Cash *Cash
}

func NewCash(logger *logrus.Logger) *Cash {
	return &Cash{
		Cache: cache.NewSnapshotCache(true, cache.IDHash{}, logger),
	}
}

func GetContext(logger *logrus.Logger) *Context {
	once.Do(func() {
		ctx = &Context{
			Cash: NewCash(logger),
		}
	})
	return ctx
}

func (c *Context) SetSnapshot(resources *resources.AllResources, logger *logrus.Logger) error {
	snapshot := GenerateSnapshot(resources)

	if err := snapshot.Consistent(); err != nil {
		logger.Fatalf("snapshot inconsistency: %+v\n%+v", snapshot, err)
	}

	logger.Debugf("Will serve snapshot %+v", snapshot)

	if err := c.Cash.Cache.SetSnapshot(context.Background(), resources.NodeID, snapshot); err != nil {
		logger.Fatalf("snapshot error %q for %+v", err, snapshot)
	}

	return nil
}
