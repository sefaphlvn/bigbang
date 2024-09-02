package dependency

import (
	"sync"
	"time"

	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

type CacheEntry struct {
	ID        string
	JSON      string
	Timestamp time.Time
	TTL       time.Duration
}

type AppHandler struct {
	Context      *db.AppContext
	Dependencies *DependencyGraph
	Cache        map[string]CacheEntry
	CacheMutex   sync.Mutex
}

type Dependency struct {
	Data struct {
		ID        string `json:"id"`
		Label     string `json:"label"`
		Category  string `json:"category"`
		Gtype     string `json:"gtype"`
		Link      string `json:"link"`
		First     bool   `json:"first"`
		Direction string `json:"direction"`
	} `json:"data"`
}

type DependencyGraph struct {
	Nodes []Dependency `json:"nodes"`
	Edges []Edge       `json:"edges"`
}

type Edge struct {
	Data struct {
		Source string `json:"source"`
		Target string `json:"target"`
		Label  string `json:"label"`
	} `json:"data"`
}

type Depend struct {
	ID         string        `json:"id"`
	Collection string        `json:"collection"`
	Name       string        `json:"name"`
	Gtype      models.GTypes `json:"gtype"`
	Project    string        `json:"project"`
	First      bool          `json:"first"`
	Direction  string        `json:"direction"`
	Source     string        `json:"source"`
}

type Node struct {
	Name       string        `json:"name"`
	Gtype      models.GTypes `json:"gtype"`
	Collection string        `json:"collection"`
	Link       string        `json:"link"`
	First      bool          `json:"first"`
	ID         string        `json:"id"`
	Direction  string        `json:"direction"`
}
