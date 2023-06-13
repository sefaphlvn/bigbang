package server

import (
	"context"
	"os"
	"sync"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/sefaphlvn/bigbang/grpcServer/server/resources"
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

func NewCash(l Logger) *Cash {
	return &Cash{
		Cache: cache.NewSnapshotCache(true, cache.IDHash{}, l),
	}
}

func GetContext(l Logger) *Context {
	once.Do(func() {
		ctx = &Context{
			Cash: NewCash(l),
		}
	})
	return ctx
}

func (cash *Context) SetSnapshot(aa *resources.AllResources, l Logger) error {
	snapshot := GenerateSnapshot(aa)

	if err := snapshot.Consistent(); err != nil {
		l.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
		os.Exit(1)
	}

	//l.Debugf("will serve snapshot %+v", snapshot)
	if err := cash.Cash.Cache.SetSnapshot(context.Background(), "test", snapshot); err != nil {
		l.Errorf("snapshot error %q for %+v", err, snapshot)
		os.Exit(1)
	}

	return nil
}
