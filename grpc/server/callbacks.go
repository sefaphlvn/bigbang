package server

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	core "github.com/sefaphlvn/versioned-go-control-plane/envoy/config/core/v3"
	discovery "github.com/sefaphlvn/versioned-go-control-plane/envoy/service/discovery/v3"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/sefaphlvn/bigbang/grpc/server/bridge"
	"github.com/sefaphlvn/bigbang/grpc/server/snapshot"
	"github.com/sefaphlvn/bigbang/pkg/db"
)

type Callbacks struct {
	signal               chan struct{}
	poke                 *bridge.PokeService
	mu                   sync.Mutex
	cache                *snapshot.Context
	activeClientsService *bridge.ActiveClientsService
	appContext           *db.AppContext
}

func NewCallbacks(poke *bridge.PokeService, cache *snapshot.Context, activeClientsService *bridge.ActiveClientsService, appContext *db.AppContext) *Callbacks {
	return &Callbacks{
		poke:                 poke,
		cache:                cache,
		activeClientsService: activeClientsService,
		appContext:           appContext,
	}
}

func (c *Callbacks) OnFetchResponse(*discovery.DiscoveryRequest, *discovery.DiscoveryResponse) {}

func (c *Callbacks) OnStreamRequest(_ int64, _ *discovery.DiscoveryRequest) error {
	return nil
}

func (c *Callbacks) OnStreamResponse(_ context.Context, _ int64, _ *discovery.DiscoveryRequest, _ *discovery.DiscoveryResponse) {
}

func (c *Callbacks) OnFetchRequest(_ context.Context, _ *discovery.DiscoveryRequest) error {
	return nil
}

func (c *Callbacks) OnStreamOpen(_ context.Context, _ int64, _ string) error {
	return nil
}

func (c *Callbacks) OnStreamClosed(_ int64, _ *core.Node) {
}

func (c *Callbacks) OnDeltaStreamOpen(ctx context.Context, id int64, typ string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	clientAddr, localAddr, nodeID, version := GetMetadata(ctx)
	if nodeID == "" {
		c.appContext.Logger.Warn("NodeID missing from metadata, skipping tracking")
		return nil
	}

	err := c.CheckSetSnapshot(nodeID, version)
	if err != nil {
		c.appContext.Logger.Warnf("Error checking snapshot: %v", err)
		return err
	}
	c.activeClientsService.TrackClient(c.appContext.Client, nodeID, clientAddr, localAddr, id)
	c.appContext.Logger.Infof("Delta stream %d opened for NodeID %s", id, nodeID)

	return nil
}

func (c *Callbacks) OnDeltaStreamClosed(id int64, node *core.Node) {
	c.mu.Lock()
	defer c.mu.Unlock()
	nodeID := node.GetId()
	if nodeID == "" {
		c.appContext.Logger.Warn("NodeID missing, skipping client cleanup")
		return
	}

	c.activeClientsService.CloseClientConnection(c.appContext.Client, c.cache.Cache.Cache, nodeID, id)
	c.appContext.Logger.Infof("Delta stream %d closed for NodeID %s", id, nodeID)
}

func (c *Callbacks) OnStreamDeltaRequest(id int64, req *discovery.DeltaDiscoveryRequest) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	nodeID := req.GetNode().GetId()

	c.activeClientsService.UpdateClientActivity(nodeID)

	typeURL := req.GetTypeUrl()
	responseNonce := req.GetResponseNonce()

	if errDetail := req.GetErrorDetail(); errDetail != nil {
		errorMsg := errDetail.Message
		c.appContext.Logger.Errorf("Delta Discovery Request Error (Node %s, Resource %s): %s", nodeID, typeURL, errorMsg)
		c.activeClientsService.AddOrUpdateError(nodeID, typeURL, errorMsg, responseNonce)
	} else {
		c.activeClientsService.ResolveErrorsForResource(nodeID, typeURL, responseNonce)
	}

	return nil
}

func (c *Callbacks) OnStreamDeltaResponse(_ int64, req *discovery.DeltaDiscoveryRequest, resp *discovery.DeltaDiscoveryResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()

	nodeID := req.GetNode().GetId()
	typeURL := req.GetTypeUrl()
	responseNonce := req.GetResponseNonce()
	nonce := resp.GetNonce()
	errMsg := req.GetErrorDetail()

	if errEntry, found := c.activeClientsService.GetErrorEntry(nodeID, typeURL, responseNonce, errMsg); found {
		errEntry.ResponseNonce = nonce
		c.activeClientsService.UpdateErrorEntry(nodeID, typeURL, responseNonce, *errEntry)
	}
}

func (c *Callbacks) CheckSetSnapshot(nodeID, version string) error {
	if nodeID == "" {
		return nil
	}

	parts := strings.Split(nodeID, ":")
	if len(parts) != 2 {
		c.appContext.Logger.Errorf("Invalid nodeID format: %s", nodeID)
		return errors.New("invalid nodeID format")
	}
	node, project := parts[0], parts[1]

	if c.poke.CheckSnapshot(nodeID) {
		return c.poke.GetResourceSetSnapshot(context.Background(), node, project, version)
	}
	return nil
}

func GetMetadata(ctx context.Context) (string, string, string, string) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", "", "", ""
	}

	clientAddr := p.Addr.String()
	localAddr := ""
	nodeID := ""
	version := ""
	if p.LocalAddr != nil {
		localAddr = p.LocalAddr.String()
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get("nodeid"); len(vals) > 0 {
			nodeID = vals[0]
			fmt.Printf("Stream opened with NodeID: %s", nodeID)
		} else {
			fmt.Println("Stream opened without NodeID in metadata")
		}

		if vals := md.Get("version"); len(vals) > 0 {
			version = vals[0]
			fmt.Printf("Stream opened with Version: %s", version)
		} else {
			fmt.Println("Stream opened without Version in metadata")
		}
	}

	return clientAddr, localAddr, nodeID, version
}
