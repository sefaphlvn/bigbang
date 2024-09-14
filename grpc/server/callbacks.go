package server

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
)

type Callbacks struct {
	signal         chan struct{}
	logger         *logrus.Logger
	fetches        int
	requests       int
	deltaRequests  int
	deltaResponses int
	mu             sync.Mutex
}

func NewCallbacks(logger *logrus.Logger) *Callbacks {
	return &Callbacks{
		logger: logger,
	}
}

func (c *Callbacks) Report() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger.Infof("server callbacks fetches=%d requests=%d\n", c.fetches, c.requests)
}

func (c *Callbacks) OnStreamOpen(_ context.Context, id int64, typ string) error {
	c.logger.Debugf("stream %d open for %s\n", id, typ)
	return nil
}

func (c *Callbacks) OnStreamClosed(id int64, node *core.Node) {
	c.logger.Debugf("stream %d of node %s closed\n", id, node.Id)

}

func (c *Callbacks) OnDeltaStreamOpen(_ context.Context, id int64, typ string) error {
	c.logger.Debugf("delta stream %d open for %s\n", id, typ)
	return nil
}

func (c *Callbacks) OnDeltaStreamClosed(id int64, node *core.Node) {
	c.logger.Debugf("delta stream %d of node %s closed\n", id, node.Id)
}

func (c *Callbacks) OnStreamRequest(_ int64, req *discovery.DiscoveryRequest) error {

	/* 	helper.PrettyPrinter(req.ResourceNames)
	   	helper.PrettyPrinter(req.TypeUrl)
	   	helper.PrettyPrinter(req.ResourceLocators)
	   	helper.PrettyPrinter(req.ResponseNonce) */
	c.logger.Infof("DiscoveryRequest: %v\n", req.TypeUrl)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.requests++
	if c.signal != nil {
		close(c.signal)
		c.signal = nil
	}
	return nil
}

func (c *Callbacks) OnStreamResponse(context.Context, int64, *discovery.DiscoveryRequest, *discovery.DiscoveryResponse) {
	fmt.Println("OnStreamResponse")
}

func (c *Callbacks) OnStreamDeltaResponse(id int64, req *discovery.DeltaDiscoveryRequest, resp *discovery.DeltaDiscoveryResponse) {
	fmt.Println()
	if req.TypeUrl == "type.googleapis.com/envoy.config.endpoint.v3.ClusterLoadAssignment" {
		c.logger.Warnf("Sending Delta EDS response: %v", resp)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.deltaResponses++

	// DeltaDiscoveryResponse nesnesini JSON formatına dönüştür
	/* respJson, err := json.Marshal(resp)
	if err != nil {
		c.logger.Errorf("JSON marshalling error: %v", err)
		return
	}

	helper.PrettyPrinter(respJson) */
	//c.logger.Debugf("DeltaDiscoveryResponse: %s", helper.PrettyPrinter(string(respJson)))
}

func (c *Callbacks) OnStreamDeltaRequest(_ int64, req *discovery.DeltaDiscoveryRequest) error {
	/* fmt.Println("-----------------")
	helper.PrettyPrint(req)
	fmt.Println("-----------------") */
	c.logger.Errorf("Delta Discovery Request: TypeUrl=%v\n", req.TypeUrl)

	/* if errDetail := req.GetErrorDetail(); errDetail != nil {
	    c.logger.Errorf("Delta Discovery Request Error: Code=%v, Message=%v\n", errDetail, errDetail.Message)
	} */

	c.mu.Lock()
	defer c.mu.Unlock()
	c.deltaRequests++
	if c.signal != nil {
		close(c.signal)
		c.signal = nil
	}

	return nil
}

func (c *Callbacks) OnFetchRequest(context.Context, *discovery.DiscoveryRequest) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.fetches++
	if c.signal != nil {
		close(c.signal)
		c.signal = nil
	}
	return nil
}

func (c *Callbacks) OnFetchResponse(*discovery.DiscoveryRequest, *discovery.DiscoveryResponse) {}
