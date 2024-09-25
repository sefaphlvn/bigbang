package snapshot

import (
	"context"
	"sync"

	xdsResource "github.com/sefaphlvn/bigbang/grpc/server/resources/resource"
	"github.com/sirupsen/logrus"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
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
	/* if err := snapshot.Consistent(); err != nil {
		logger.Fatalf("snapshot inconsistency: %+v\n%+v", snapshot, err)
	} */

	logger.Debugf("end serve snapshot: (%s)", resources.NodeID)
	if err := c.Cache.Cache.SetSnapshot(context.Background(), resources.NodeID, snapshot); err != nil {
		logger.Fatalf("snapshot error %q for %+v", err, snapshot)
	}

	logger.Infof("Successfully set snapshot for nodeID: %s", resources.NodeID)

	return nil
}

func GenerateSnapshot(r *xdsResource.AllResources) *cache.Snapshot {
	snap, _ := cache.NewSnapshot(r.GetVersion(),
		map[resource.Type][]types.Resource{
			resource.ClusterType:         r.GetClusterT(),
			resource.RouteType:           r.GetRouteT(),
			resource.VirtualHostType:     r.GetVirtualHostT(),
			resource.EndpointType:        r.GetEndpointT(),
			resource.ListenerType:        r.GetListenerT(),
			resource.ExtensionConfigType: r.GetExtensionsT(),
			resource.SecretType:          r.GetSecretT(),
		},
	)

	return snap
}
