package models

import (
	"time"
)

type Client struct {
	ServerAddress   string
	NodeID          string
	FirstConnected  int64
	LastSeen        int64
	LastActivity    int64
	ConnectionCount int64
	RequestCount    int64
	ClientAddr      string
	LocalAddr       string
	StreamIDs       []int64
	Errors          *BoundedCache
}

type ActiveClients struct {
	Clients map[string]*Client
}

type ErrorEntry struct {
	Message       string
	ResourceID    string
	ResponseNonce string
	Timestamp     time.Time
	Count         int
	Resolved      bool
}

type BoundedCache struct {
	Errors []ErrorEntry
	Limit  int
}

func NewBoundedCache(limit int) *BoundedCache {
	return &BoundedCache{
		Errors: []ErrorEntry{},
		Limit:  limit,
	}
}
