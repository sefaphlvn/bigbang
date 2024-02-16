package server

import (
	"context"
	"sync"

	xdsResource "github.com/sefaphlvn/bigbang/grpc/server/resources/resource"
	"github.com/sirupsen/logrus"

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

func (c *Context) SetSnapshot(resources *xdsResource.AllResources, logger *logrus.Logger) error {
	snapshot := GenerateSnapshot(resources)

	// helper.PrettyPrinter(snapshot)
	if err := snapshot.Consistent(); err != nil {
		logger.Fatalf("snapshot inconsistency: %+v\n%+v", snapshot, err)
	}

	logger.Debugf("Will serve snapshot %+v", snapshot)

	if err := c.Cache.Cache.SetSnapshot(context.Background(), resources.GetNodeID(), snapshot); err != nil {
		logger.Fatalf("snapshot error %q for %+v", err, snapshot)
	}

	return nil
}
