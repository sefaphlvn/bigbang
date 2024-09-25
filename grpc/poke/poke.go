package poke

import (
	"github.com/sefaphlvn/bigbang/grpc/server/resources/resource"
	"github.com/sefaphlvn/bigbang/grpc/server/snapshot"
	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sefaphlvn/bigbang/pkg/db"
)

type Poke struct {
	ctx    *snapshot.Context
	db     *db.AppContext
	logger *logrus.Logger
	config *config.AppConfig
}

type Nodes struct {
	Project string
	Service string
}

func NewPokeServer(ctx *snapshot.Context, db *db.AppContext, logger *logrus.Logger, config *config.AppConfig) *Poke {
	return &Poke{
		ctx:    ctx,
		db:     db,
		logger: logger,
		config: config,
	}
}

func (p *Poke) Run(pokeHandler *Poke) {
	p.initialSnapshots()
}

func (p *Poke) initialSnapshots() {
	nodes := p.getListenerList()
	for _, node := range nodes {
		p.logger.Debugf("start serve snapshot: (%s) Project: (%s)", node.Service, node.Project)

		allResource, err := p.getAllResourcesFromListener(node.Service, node.Project)
		if err != nil {
			p.logger.Errorf("BULK GetConfigurationFromListener(%v): %v", node.Service, err)
		}

		err = p.ctx.SetSnapshot(allResource, p.logger)
		if err != nil {
			p.logger.Errorf("%s", err)
		}
	}
	p.logger.Infof("All snapshots are loaded")
}

func (p *Poke) getListenerList() []Nodes {
	var serviceNamesWithProject []Nodes
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

		serviceNamesWithProject = append(serviceNamesWithProject, Nodes{Project: general.Project, Service: general.Name})
	}
	return serviceNamesWithProject
}

func (p *Poke) getAllResourcesFromListener(serviceName string, project string) (*resource.AllResources, error) {
	rawListenerResource, err := resources.GetResourceNGeneral(p.db, "listeners", serviceName, project)
	if err != nil {
		return nil, err
	}

	lis, err := resource.GenerateSnapshot(rawListenerResource, serviceName, p.db, p.logger, project)
	if err != nil {
		return nil, err
	}

	return lis, nil
}
