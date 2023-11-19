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

type Cache struct {
	Cache cache.SnapshotCache
}

type Context struct {
	Cache *Cache
}

func NewCache(logger *logrus.Logger) *Cache {
	return &Cache{
		Cache: cache.NewSnapshotCache(true, cache.IDHash{}, logger),
	}
}

func GetContext(logger *logrus.Logger) *Context {
	once.Do(func() {
		ctx = &Context{
			Cache: NewCache(logger),
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

	if err := c.Cache.Cache.SetSnapshot(context.Background(), resources.NodeID, snapshot); err != nil {
		logger.Fatalf("snapshot error %q for %+v", err, snapshot)
	}

	return nil
}
