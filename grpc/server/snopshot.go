package server

import (
	"context"
	"github.com/sefaphlvn/bigbang/grpc/server/resources"
	"github.com/sirupsen/logrus"
	"os"
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

func NewCash(l *logrus.Logger) *Cash {
	return &Cash{
		Cache: cache.NewSnapshotCache(true, cache.IDHash{}, l),
	}
}

func GetContext(l *logrus.Logger) *Context {
	once.Do(func() {
		ctx = &Context{
			Cash: NewCash(l),
		}
	})
	return ctx
}

func (cash *Context) SetSnapshot(resources *resources.AllResources, l *logrus.Logger) error {
	snapshot := GenerateSnapshot(resources)

	if err := snapshot.Consistent(); err != nil {
		l.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
		os.Exit(1)
	}

	//l.Debugf("will serve snapshot %+v", snapshot)
	if err := cash.Cash.Cache.SetSnapshot(context.Background(), resources.NodeID, snapshot); err != nil {
		l.Errorf("snapshot error %q for %+v", err, snapshot)
		os.Exit(1)
	}

	return nil
}
