package poke

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sefaphlvn/bigbang/grpc/server/resources/resource"
	"github.com/sefaphlvn/bigbang/grpc/server/snapshot"
	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
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

// Run starts the main execution of the program.
// It initializes necessary components and begins processing based on the provided configuration.
// Parameters:
// - ctx: context for controlling the request lifetime
// Returns:
// - error: an error if any occurred during the execution process
func (p *Poke) Run() {
	ctx := context.Background()
	p.initialSnapshots(ctx)
}

// initialSnapshots initializes the snapshots for the given context and database.
// It retrieves the initial state of the resources and prepares them for further processing.
// Parameters:
// - ctx: context for controlling the request lifetime
// - db: database connection to fetch the resources from
// Returns:
// - error: an error if any occurred during the initialization process
func (p *Poke) initialSnapshots(ctx context.Context) {
	nodes := p.getListenerList(ctx)
	for _, node := range nodes {
		p.logger.Debugf("start serve snapshot: (%s) Project: (%s)", node.Service, node.Project)

		allResource, err := p.getAllResourcesFromListener(ctx, node.Service, node.Project)
		if err != nil {
			p.logger.Errorf("BULK GetConfigurationFromListener(%v): %v", node.Service, err)
		}

		err = p.ctx.SetSnapshot(ctx, allResource, p.logger)
		if err != nil {
			p.logger.Errorf("%s", err)
		}
	}
	p.logger.Infof("All snapshots are loaded")
}

// getListenerList retrieves a list of listeners from the database.
// It fetches the general information of each listener and appends it to the serviceNamesWithProject slice.
// Parameters:
// - ctx: context for controlling the request lifetime
// Returns:
// - []Nodes: a slice containing the project and service names of the listeners
// - error: an error if any occurred during the process
func (p *Poke) getListenerList(ctx context.Context) []Nodes {
	var serviceNamesWithProject []Nodes
	cur, err := p.db.GetGenerals(ctx, "listeners")
	if err != nil {
		p.logger.Fatal(err)
	}

	for cur.Next(ctx) {
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

// getAllResourcesFromListener retrieves all resources from a listener based on the provided service name and project.
// It fetches the raw listener resource and generates a snapshot of the resource.
// Parameters:
// - ctx: context for controlling the request lifetime
// - serviceName: name of the service to fetch resources for
// - project: name of the project to fetch resources for
// Returns:
// - *resource.AllResources: a pointer to the structure containing all resources
// - error: an error if any occurred during the process
func (p *Poke) getAllResourcesFromListener(ctx context.Context, serviceName, project string) (*resource.AllResources, error) {
	rawListenerResource, err := resources.GetResourceNGeneral(ctx, p.db, "listeners", serviceName, project)
	if err != nil {
		return nil, err
	}

	lis, err := resource.GenerateSnapshot(ctx, rawListenerResource, serviceName, p.db, p.logger, project)
	if err != nil {
		return nil, err
	}

	return lis, nil
}
