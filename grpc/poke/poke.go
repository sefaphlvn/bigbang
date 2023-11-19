package poke

import (
	"encoding/json"
	"fmt"
	"github.com/sefaphlvn/bigbang/grpc/models"
	"github.com/sefaphlvn/bigbang/grpc/server"
	"github.com/sefaphlvn/bigbang/grpc/server/resources"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"

	"github.com/sefaphlvn/bigbang/pkg/db"
)

type Poke struct {
	ctx    *server.Context
	db     *db.MongoDB
	logger *logrus.Logger
}

func NewPokeServer(ctx *server.Context, db *db.MongoDB, logger *logrus.Logger) *Poke {
	return &Poke{
		ctx:    ctx,
		db:     db,
		logger: logger,
	}
}

func (p *Poke) Run(pokeHandler *Poke) {
	p.initialSnapshots()
	p.logger.Infof("Poke server listening on :8080")
	if err := http.ListenAndServe(":8080", pokeHandler); err != nil {
		p.logger.Fatalf("failed to start HTTP server: %v", err)
	}
}

func (p *Poke) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/ping":
		p.handlePing(w, r)
	case "/poke":
		p.handlePoke(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (p *Poke) handlePing(w http.ResponseWriter, r *http.Request) {
	p.logger.Info(fmt.Fprint(w, "OK"))
}

func (p *Poke) handlePoke(w http.ResponseWriter, r *http.Request) {
	allResources, err := p.getAllResourcesFromListener("newListener")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		p.logger.Error(err)
		return
	}

	err = p.ctx.SetSnapshot(allResources, p.logger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		p.logger.Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(allResources)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		p.logger.Error(err)
		return
	}
}

func (p *Poke) getAllResourcesFromListener(serviceName string) (*resources.AllResources, error) {
	rawListenerResource, err := resources.GetResource(p.db, "listeners", serviceName)
	if err != nil {
		return nil, err
	}

	lis, err := resources.SetSnapshot(rawListenerResource, serviceName, p.db, p.logger)
	if err != nil {
		return nil, err
	}

	return lis, nil
}

func (p *Poke) initialSnapshots() {
	serviceNames := p.getListenerList()
	for _, serviceName := range serviceNames {
		allResource, err := p.getAllResourcesFromListener(serviceName)
		if err != nil {
			p.logger.Errorf("BULK GetConfigurationFromListener(%v): %v", serviceName, err)
		}
		err = p.ctx.SetSnapshot(allResource, p.logger)
		if err != nil {
			p.logger.Errorf("%s", err)
		}
	}
	p.logger.Infof("All snapshots are loaded")
}

func (p *Poke) getListenerList() []string {
	var serviceNames []string
	cur, err := p.db.GetGenerals("listeners")
	if err != nil {
		p.logger.Fatal(err)
	}

	for cur.Next(p.db.Ctx) {
		var result bson.M
		err = cur.Decode(&result)
		if err != nil {
			p.logger.Fatal(err)
		}

		var general models.General
		bsonBytes, _ := bson.Marshal(result["general"])

		err = bson.Unmarshal(bsonBytes, &general)
		if err != nil {
			p.logger.Errorf("%s", err)
			return nil
		}

		serviceNames = append(serviceNames, general.Name)
	}
	return serviceNames
}
