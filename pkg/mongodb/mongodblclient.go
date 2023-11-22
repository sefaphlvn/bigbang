package mongodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewClient(appConfig *config.AppConfig, log *logrus.Logger) DBClient {

	hosts := strings.Join(appConfig.MongoDB.Hosts, fmt.Sprintf("%s,", appConfig.MongoDB.Port))
	connectionString := fmt.Sprintf("%s://%s:%s@%s%s", appConfig.MongoDB.Scheme, appConfig.MongoDB.Username, appConfig.MongoDB.Password, hosts, appConfig.MongoDB.Port)

	clientOptions := options.Client().ApplyURI(connectionString)

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), appConfig.MongoDB.TimeoutSeconds)
	defer cancel()

	client, err := mongo.Connect(ctxWithTimeout, clientOptions)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return &dbClient{client: client, logger: log}
}

type DBClient interface {
	NewSession(dbName string) *mongo.Database
	NewSessionWithSecondaryPreferred(dbName string) *mongo.Database
	EnsureIndex(indexKeys []string, isUnique bool, indexName, dbName, collection string) error
	Ping() error
}

type dbClient struct {
	client *mongo.Client
	logger *logrus.Logger
}

func (c *dbClient) NewSessionWithSecondaryPreferred(dbName string) *mongo.Database {

	secondary := readpref.SecondaryPreferred()
	dbOpts := options.Database().SetReadPreference(secondary)

	return c.client.Database(dbName, dbOpts)
}

func (c *dbClient) NewSession(dbName string) *mongo.Database {
	return c.client.Database(dbName)
}

func (c *dbClient) EnsureIndex(indexKeys []string, isUnique bool, indexName, dbName, collection string) error {

	serviceCollection := c.client.Database(dbName).Collection(collection)

	_, err := serviceCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    generateIndexKeys(indexKeys),
		Options: options.Index().SetName(indexName).SetUnique(isUnique)})

	if err != nil {
		return err
	}

	return nil
}

func (c *dbClient) Ping() error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := c.client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	return nil
}

func generateIndexKeys(arr []string) bson.D {

	var keys bson.D

	for _, s := range arr {
		keys = append(keys, bson.E{
			Key:   s,
			Value: int32(1),
		})
	}

	return keys
}
