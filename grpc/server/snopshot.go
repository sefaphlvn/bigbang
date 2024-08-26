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

/* func (c *Context) SetSnapshot(resources *xdsResource.AllResources, logger *logrus.Logger) error {
	snapshot := GenerateSnapshot(resources)
	if err := snapshot.Consistent(); err != nil {
		logger.Fatalf("snapshot inconsistency: %+v\n%+v", snapshot, err)
	}

	logger.Debugf("end serve snapshot: (%s)", resources.NodeID)
	if err := c.Cache.Cache.SetSnapshot(context.Background(), resources.GetNodeID(), snapshot); err != nil {
		logger.Fatalf("snapshot error %q for %+v", err, snapshot)
	}
	aa := c.Cache.Cache.GetStatusKeys()
	fmt.Println(aa)

	return nil
} */

func (c *Context) SetSnapshot(resources *xdsResource.AllResources, logger *logrus.Logger) error {
	snapshot := GenerateSnapshot(resources)
	if err := snapshot.Consistent(); err != nil {
		logger.Fatalf("snapshot inconsistency: %+v\n%+v", snapshot, err)
	}

	logger.Debugf("end serve snapshot: (%s)", resources.NodeID)
	if err := c.Cache.Cache.SetSnapshot(context.Background(), resources.NodeID, snapshot); err != nil {
		logger.Fatalf("snapshot error %q for %+v", err, snapshot)
	}

	// Debug mesajları ekleyerek snapshot'ın gerçekten ayarlandığını doğrulayın
	logger.Infof("Successfully set snapshot for nodeID: %s", resources.NodeID)

	// Snapshot ayarlandıktan sonra durumu kontrol edin
	keys := c.Cache.Cache.GetStatusInfo(resources.NodeID)
	logger.Infof("Current snapshot keys: %v", keys)

	return nil
}
