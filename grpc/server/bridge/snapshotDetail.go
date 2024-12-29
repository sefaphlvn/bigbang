package bridge

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/sefaphlvn/bigbang/grpc/grpcserver"
	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/genproto/googleapis/rpc/status"
)

type ActiveClientsService struct {
	ActiveClients *models.ActiveClients
	logger        *logrus.Logger
	mu            sync.Mutex
}

func NewActiveClientsService(logger *logrus.Logger) *ActiveClientsService {
	return &ActiveClientsService{
		ActiveClients: &models.ActiveClients{
			Clients: make(map[string]*models.Client),
		},
		logger: logger,
	}
}

func (acss *ActiveClientsServiceServer) GetActiveClient(ctx context.Context, req *bridge.NodeRequest) (*bridge.ActiveClientResponse, error) {
	client, exists := acss.activeClients.Clients[req.GetNodeId()]
	if !exists {
		return nil, fmt.Errorf("client with NodeID %s not found", req.GetNodeId())
	}

	return &bridge.ActiveClientResponse{
		Client: &bridge.Client{
			ServerAddress:   client.ServerAddress,
			NodeId:          client.NodeID,
			FirstConnected:  client.FirstConnected,
			LastSeen:        client.LastSeen,
			ConnectionCount: client.ConnectionCount,
			ClientAddr:      client.ClientAddr,
			LocalAddr:       client.LocalAddr,
			StreamIds:       client.StreamIDs,
			RequestCount:    client.RequestCount,
			Errors:          convertErrorsToProto(client.Errors.Errors),
		},
	}, nil
}

func (acs *ActiveClientsService) UpdateClientActivity(nodeID string) {
	acs.mu.Lock()
	defer acs.mu.Unlock()

	client, exists := acs.ActiveClients.Clients[nodeID]
	if exists {
		client.LastSeen = time.Now().Unix()
		client.RequestCount++
		acs.logger.Debugf("Activity updated for client %s, request count: %d", nodeID, client.RequestCount)
	} else {
		acs.logger.Warnf("Client with NodeID %s not found while updating activity", nodeID)
	}
}

func (acs *ActiveClientsService) UpdateErrorEntry(nodeID, resourceID, nonce string, updatedEntry models.ErrorEntry) {
	acs.mu.Lock()
	defer acs.mu.Unlock()

	for i, entry := range acs.ActiveClients.Clients[nodeID].Errors.Errors {
		if entry.ResourceID == resourceID && entry.ResponseNonce == nonce {
			acs.ActiveClients.Clients[nodeID].Errors.Errors[i] = updatedEntry
			return
		}
	}
}

func (acs *ActiveClientsService) TrackClient(dbClient *mongo.Database, nodeID, clientAddr, localAddr string, streamID int64) {
	go grpcserver.AddOrUpdateGrpcServer(dbClient, localAddr, nodeID)
	acs.mu.Lock()
	defer acs.mu.Unlock()

	client, exists := acs.ActiveClients.Clients[nodeID]
	if exists {
		client.ConnectionCount++
		client.LastSeen = time.Now().Unix()
		client.StreamIDs = append(client.StreamIDs, streamID)
		acs.logger.Debugf("Reconnected to node %s, connection count: %d", nodeID, client.ConnectionCount)
	} else {
		acs.ActiveClients.Clients[nodeID] = &models.Client{
			ServerAddress:   grpcserver.GetHostname(),
			NodeID:          nodeID,
			FirstConnected:  time.Now().Unix(),
			LastSeen:        time.Now().Unix(),
			ConnectionCount: 1,
			ClientAddr:      clientAddr,
			LocalAddr:       localAddr,
			StreamIDs:       []int64{streamID},
			RequestCount:    0,
			Errors:          models.NewBoundedCache(25),
		}
		acs.logger.Debugf("New client tracked: %s", nodeID)
	}
}

func (acs *ActiveClientsService) CloseClientConnection(dbClient *mongo.Database, cache cache.SnapshotCache, nodeID string, streamID int64) {
	acs.mu.Lock()
	defer acs.mu.Unlock()

	client, exists := acs.ActiveClients.Clients[nodeID]
	if !exists {
		acs.logger.Warnf("Client with NodeID %s not found", nodeID)
		return
	}

	for i, id := range client.StreamIDs {
		if id == streamID {
			client.StreamIDs = append(client.StreamIDs[:i], client.StreamIDs[i+1:]...)
			break
		}
	}

	client.ConnectionCount--
	if client.ConnectionCount <= 0 {
		delete(acs.ActiveClients.Clients, nodeID)
		go grpcserver.RemoveNodeID(dbClient, nodeID)
		cache.ClearSnapshot(nodeID)
		acs.logger.Infof("Client with NodeID %s removed", nodeID)
	}
}

func (acs *ActiveClientsService) GetClients() map[string]*models.Client {
	acs.mu.Lock()
	defer acs.mu.Unlock()

	return acs.ActiveClients.Clients
}

func (acs *ActiveClientsService) AddOrUpdateError(nodeID, resourceID, errorMsg, nonce string) {
	acs.mu.Lock()
	defer acs.mu.Unlock()

	client, exists := acs.ActiveClients.Clients[nodeID]
	if !exists {
		acs.logger.Warnf("Client with NodeID %s not found when adding/updating error", nodeID)
		return
	}

	if client.Errors == nil {
		client.Errors = models.NewBoundedCache(10)
	}

	for i, entry := range client.Errors.Errors {
		if entry.ResourceID == resourceID && entry.Message == errorMsg {
			entry.ResponseNonce = nonce
			entry.Timestamp = time.Now()
			entry.Count++
			entry.Resolved = false
			client.Errors.Errors[i] = entry
			return
		}
	}

	client.Errors.Errors = append(client.Errors.Errors, models.ErrorEntry{
		Message:       errorMsg,
		ResourceID:    resourceID,
		ResponseNonce: nonce,
		Timestamp:     time.Now(),
		Count:         1,
		Resolved:      false,
	})

	if len(client.Errors.Errors) > client.Errors.Limit {
		client.Errors.Errors = client.Errors.Errors[1:]
	}
}

func convertErrorsToProto(errors []models.ErrorEntry) []*bridge.ErrorEntry {
	var protoErrors []*bridge.ErrorEntry
	for _, err := range errors {
		protoErrors = append(protoErrors, &bridge.ErrorEntry{
			Message:       err.Message,
			ResourceId:    err.ResourceID,
			ResponseNonce: err.ResponseNonce,
			Count:         int32(err.Count),
			Resolved:      err.Resolved,
			Timestamp:     err.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return protoErrors
}

func (acs *ActiveClientsService) GetErrorEntry(nodeID, resourceID, nonce string, errMsg *status.Status) (*models.ErrorEntry, bool) {
	acs.mu.Lock()
	defer acs.mu.Unlock()

	client, exists := acs.ActiveClients.Clients[nodeID]
	if !exists {
		return nil, false
	}

	if errMsg == nil {
		return nil, false
	}

	for _, entry := range client.Errors.Errors {
		if entry.ResourceID == resourceID && entry.ResponseNonce == nonce {
			//entry.Resolved = true
			return &entry, true
		}
	}
	return nil, false
}

func (acs *ActiveClientsService) ResolveErrorsForResource(nodeID, resourceID, nonce string) {
	acs.mu.Lock()
	defer acs.mu.Unlock()

	if acs.ActiveClients == nil {
		acs.logger.Debugf("ActiveClients is nil for nodeID: %s", nodeID)
		return
	}

	client, ok := acs.ActiveClients.Clients[nodeID]
	if !ok || client == nil {
		acs.logger.Debugf("Client not found or nil for nodeID: %s", nodeID)
		return
	}

	if client.Errors == nil {
		acs.logger.Debugf("Errors is nil for nodeID: %s", nodeID)
		return
	}

	for i, entry := range client.Errors.Errors {
		if entry.ResourceID == resourceID && entry.ResponseNonce == nonce {
			entry.Resolved = true
			client.Errors.Errors[i] = entry
			acs.logger.Debugf("Error resolved for resourceID: %s, nonce: %s, nodeID: %s", resourceID, nonce, nodeID)
			return
		}
	}

	acs.logger.Debugf("No matching error found for resourceID: %s, nonce: %s, nodeID: %s", resourceID, nonce, nodeID)
}
