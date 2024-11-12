package server

import (
	"context"
	"fmt"
	"sync"

	"github.com/sefaphlvn/bigbang/grpc/server/bridge"
	"github.com/sirupsen/logrus"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
)

type Callbacks struct {
	signal         chan struct{}
	logger         *logrus.Logger
	errorContext   *bridge.ErrorContext
	deltaRequests  int
	deltaResponses int
	mu             sync.Mutex
}

func NewCallbacks(logger *logrus.Logger, errorContext *bridge.ErrorContext) *Callbacks {
	return &Callbacks{
		logger:       logger,
		errorContext: errorContext,
	}
}

func (c *Callbacks) OnFetchResponse(*discovery.DiscoveryRequest, *discovery.DiscoveryResponse) {}

func (c *Callbacks) OnStreamResponse(_ context.Context, _ int64, _ *discovery.DiscoveryRequest, _ *discovery.DiscoveryResponse) {
}

func (c *Callbacks) OnFetchRequest(_ context.Context, _ *discovery.DiscoveryRequest) error {
	return nil
}

func (c *Callbacks) OnStreamOpen(_ context.Context, id int64, typ string) error {
	fmt.Println("stream open", id, typ)
	return nil
}

func (c *Callbacks) OnStreamClosed(id int64, node *core.Node) {
	fmt.Println("stream closed", id, node)
}

func (c *Callbacks) OnDeltaStreamOpen(_ context.Context, id int64, typ string) error {
	c.logger.Debugf("delta stream %d open for %s\n", id, typ)
	return nil
}

func (c *Callbacks) OnDeltaStreamClosed(id int64, node *core.Node) {
	c.logger.Debugf("delta stream %d of node %s closed\n", id, node.Id)
}

func (c *Callbacks) OnStreamRequest(_ int64, _ *discovery.DiscoveryRequest) error { return nil }

func (c *Callbacks) OnStreamDeltaResponse(_ int64, req *discovery.DeltaDiscoveryRequest, resp *discovery.DeltaDiscoveryResponse) {
	//Testit(nil, resp)
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Println(req.ErrorDetail)
	c.deltaResponses++

	nodeID := req.GetNode().GetId()
	typeURL := req.GetTypeUrl()
	responseNonce := req.GetResponseNonce()
	nonce := resp.GetNonce()
	c.logger.Warnf("respnonce: %s\n nonce: %s\n typeurl: %s\n", responseNonce, nonce, typeURL)
	if errEntry, found := c.errorContext.ErrorCache.GetErrorEntry(nodeID, typeURL, responseNonce); found {
		errEntry.ResponseNonce = nonce
		c.errorContext.ErrorCache.UpdateErrorEntry(nodeID, typeURL, responseNonce, *errEntry)
	}
}

func (c *Callbacks) OnStreamDeltaRequest(_ int64, req *discovery.DeltaDiscoveryRequest) error {
	//Testit(req, nil)

	nodeID := req.GetNode().GetId()
	typeURL := req.GetTypeUrl()
	responseNonce := req.GetResponseNonce()

	if errDetail := req.GetErrorDetail(); errDetail != nil {
		errorMsg := errDetail.Message
		c.logger.Errorf("Delta Discovery Request Error (Node %s, Resource %s): %s", nodeID, typeURL, errorMsg)
		c.errorContext.ErrorCache.AddOrUpdateError(nodeID, typeURL, errorMsg, responseNonce)
	} else {
		c.errorContext.ErrorCache.ResolveErrorsForResource(nodeID, typeURL, responseNonce)
		// c.errorContext.ErrorCache.ClearResolvedErrors(nodeID)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.deltaRequests++
	if c.signal != nil {
		close(c.signal)
		c.signal = nil
	}

	return nil
}
