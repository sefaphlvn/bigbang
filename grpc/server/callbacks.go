package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"log"
	"sync"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
)

type Callbacks struct {
	Signal         chan struct{}
	Logger         *logrus.Logger
	Fetches        int
	Requests       int
	DeltaRequests  int
	DeltaResponses int
	mu             sync.Mutex
}

var _ server.Callbacks = &Callbacks{}

func (c *Callbacks) Report() {
	c.mu.Lock()
	defer c.mu.Unlock()
	log.Printf("server callbacks fetches=%d requests=%d\n", c.Fetches, c.Requests)
}

func (c *Callbacks) OnStreamOpen(_ context.Context, id int64, typ string) error {
	c.Logger.Debugf("stream %d open for %s\n", id, typ)
	return nil
}

func (c *Callbacks) OnStreamClosed(id int64, node *core.Node) {
	c.Logger.Debugf("stream %d of node %s closed\n", id, node.Id)

}

func (c *Callbacks) OnDeltaStreamOpen(_ context.Context, id int64, typ string) error {
	c.Logger.Debugf("delta stream %d open for %s\n", id, typ)
	return nil
}

func (c *Callbacks) OnDeltaStreamClosed(id int64, node *core.Node) {
	c.Logger.Debugf("delta stream %d of node %s closed\n", id, node.Id)
}

func (c *Callbacks) OnStreamRequest(_ int64, req *discovery.DiscoveryRequest) error {
	log.Printf("DiscoveryRequest: %v\n", req.TypeUrl)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Requests++
	if c.Signal != nil {
		close(c.Signal)
		c.Signal = nil
	}
	return nil
}

func (c *Callbacks) OnStreamResponse(context.Context, int64, *discovery.DiscoveryRequest, *discovery.DiscoveryResponse) {
}

func (c *Callbacks) OnStreamDeltaResponse(int64, *discovery.DeltaDiscoveryRequest, *discovery.DeltaDiscoveryResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.DeltaResponses++
}

func (c *Callbacks) OnStreamDeltaRequest(_ int64, req *discovery.DeltaDiscoveryRequest) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.DeltaRequests++
	if c.Signal != nil {
		close(c.Signal)
		c.Signal = nil
	}

	return nil
}

func (c *Callbacks) OnFetchRequest(context.Context, *discovery.DiscoveryRequest) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Fetches++
	if c.Signal != nil {
		close(c.Signal)
		c.Signal = nil
	}
	return nil
}

func (c *Callbacks) OnFetchResponse(*discovery.DiscoveryRequest, *discovery.DiscoveryResponse) {}
