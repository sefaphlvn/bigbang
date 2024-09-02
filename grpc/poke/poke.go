package poke

import (
	"fmt"
	"net/http"

	"github.com/sefaphlvn/bigbang/grpc/server"
	"github.com/sefaphlvn/bigbang/grpc/server/resources/resource"
	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sefaphlvn/bigbang/pkg/db"
)

type Poke struct {
	ctx    *server.Context
	db     *db.AppContext
	logger *logrus.Logger
	config *config.AppConfig
}

type ProjectServices struct {
	Project string
	Service string
}

func NewPokeServer(ctx *server.Context, db *db.AppContext, logger *logrus.Logger, config *config.AppConfig) *Poke {
	return &Poke{
		ctx:    ctx,
		db:     db,
		logger: logger,
		config: config,
	}
}

func (p *Poke) Run(pokeHandler *Poke) {
	p.initialSnapshots()
	p.logger.Infof("Poke server listening on :%s", p.config.BIGBANG_GRPC_POKE_PORT)
	address := fmt.Sprintf(":%s", p.config.BIGBANG_GRPC_POKE_PORT)
	if err := http.ListenAndServe(address, pokeHandler); err != nil {
		p.logger.Fatalf("failed to start HTTP server: %v", err)
	}
}

func (p *Poke) initialSnapshots() {
	serviceNames := p.getListenerList()
	for _, serviceName := range serviceNames {
		p.logger.Debugf("start serve snapshot: (%s) Project: (%s)", serviceName.Service, serviceName.Project)

		allResource, err := p.getAllResourcesFromListener(serviceName.Service, serviceName.Project)
		if err != nil {
			p.logger.Errorf("BULK GetConfigurationFromListener(%v): %v", serviceName.Service, err)
		}

		err = p.ctx.SetSnapshot(allResource, p.logger)
		if err != nil {
			p.logger.Errorf("%s", err)
		}
	}
	p.logger.Infof("All snapshots are loaded")
}

func (p *Poke) getListenerList() []ProjectServices {
	var serviceNamesWithProject []ProjectServices
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

		serviceNamesWithProject = append(serviceNamesWithProject, ProjectServices{Project: general.Project, Service: general.Name})
	}
	return serviceNamesWithProject
}

func (p *Poke) getAllResourcesFromListener(serviceName string, project string) (*resource.AllResources, error) {
	rawListenerResource, err := resources.GetResourceNGeneral(p.db, "listeners", serviceName, project)
	if err != nil {
		return nil, err
	}

	lis, err := resource.SetSnapshot(rawListenerResource, serviceName, p.db, p.logger)
	if err != nil {
		return nil, err
	}

	return lis, nil
}
