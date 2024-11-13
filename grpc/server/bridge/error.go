package bridge

import (
	"context"
	"sync"
	"time"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
)

type ErrorEntry struct {
	Message       string
	ResourceID    string
	ResponseNonce string
	Timestamp     time.Time
	Count         int
	Resolved      bool
}

type BoundedCache struct {
	mu     sync.Mutex
	errors map[string][]ErrorEntry
	limit  int
}

func NewBoundedCache(limit int) *BoundedCache {
	return &BoundedCache{
		errors: make(map[string][]ErrorEntry),
		limit:  limit,
	}
}

type ErrorContext struct {
	ErrorCache *BoundedCache
}

func NewErrorContext(cacheLimit int) *ErrorContext {
	return &ErrorContext{
		ErrorCache: NewBoundedCache(cacheLimit),
	}
}

func (b *BoundedCache) AddOrUpdateError(nodeID, resourceID, errorMsg, nonce string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.errors[nodeID] == nil {
		b.errors[nodeID] = []ErrorEntry{}
	}

	for i, entry := range b.errors[nodeID] {
		if entry.ResourceID == resourceID && entry.Message == errorMsg {
			entry.ResponseNonce = nonce
			entry.Timestamp = time.Now()
			entry.Count++
			entry.Resolved = false
			b.errors[nodeID][i] = entry
			return
		}
	}

	b.errors[nodeID] = append(b.errors[nodeID], ErrorEntry{
		Message:       errorMsg,
		ResourceID:    resourceID,
		ResponseNonce: nonce,
		Timestamp:     time.Now(),
		Count:         1,
		Resolved:      false,
	})

	if len(b.errors[nodeID]) > b.limit {
		b.evictOldestError(nodeID)
	}
}

func (b *BoundedCache) GetErrors(nodeID string) []ErrorEntry {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.errors[nodeID]
}

func (b *BoundedCache) evictOldestError(nodeID string) {
	var oldestIndex int
	var oldestTime time.Time

	for i, entry := range b.errors[nodeID] {
		if i == 0 || entry.Timestamp.Before(oldestTime) {
			oldestIndex = i
			oldestTime = entry.Timestamp
		}
	}
	b.errors[nodeID] = append(b.errors[nodeID][:oldestIndex], b.errors[nodeID][oldestIndex+1:]...)
}

// Hata çözülmüş olarak işaretler.
func (b *BoundedCache) ResolveErrorsForResource(nodeID, resourceID, nonce string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, entry := range b.errors[nodeID] {
		if entry.ResourceID == resourceID && entry.ResponseNonce == nonce {
			// Hata çözülmüş olarak işaretleniyor, ancak listeden silinmiyor
			entry.Resolved = true
			b.errors[nodeID][i] = entry
			return
		}
	}
}

// Çözülmüş hataları temizler.
func (b *BoundedCache) ClearResolvedErrors(nodeID string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	validErrors := []ErrorEntry{}
	for _, entry := range b.errors[nodeID] {
		if !entry.Resolved {
			validErrors = append(validErrors, entry)
		}
	}
	b.errors[nodeID] = validErrors
}

func (s *ErrorServiceServer) GetNodeErrors(_ context.Context, req *bridge.NodeErrorRequest) (*bridge.NodeErrorResponse, error) {
	nodeID := req.GetNodeId()
	errors := s.errorContext.ErrorCache.GetErrors(nodeID)

	response := &bridge.NodeErrorResponse{}
	for _, entry := range errors {
		response.Errors = append(response.Errors, &bridge.ErrorEntry{
			Message:       entry.Message,
			ResourceId:    entry.ResourceID,
			ResponseNonce: entry.ResponseNonce,
			Count:         int32(entry.Count),
			Resolved:      entry.Resolved,
			Timestamp:     entry.Timestamp.Format(time.RFC3339),
		})
	}

	return response, nil
}

func (b *BoundedCache) GetErrorEntry(nodeID, resourceID, nonce string) (*ErrorEntry, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, entry := range b.errors[nodeID] {
		if entry.ResourceID == resourceID && entry.ResponseNonce == nonce {
			return &entry, true
		}
	}
	return nil, false
}

func (b *BoundedCache) UpdateErrorEntry(nodeID, resourceID, nonce string, updatedEntry ErrorEntry) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, entry := range b.errors[nodeID] {
		if entry.ResourceID == resourceID && entry.ResponseNonce == nonce {
			b.errors[nodeID][i] = updatedEntry
			return
		}
	}
}
