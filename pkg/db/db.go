package db

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WTF struct {
	Client *mongo.Database
	Ctx    context.Context
	Logger *logrus.Logger
	Config *config.AppConfig
}

func NewMongoDB(config *config.AppConfig, logger *logrus.Logger) *WTF {
	hosts := strings.Join(config.MongoDB.Hosts, fmt.Sprintf("%s,", config.MongoDB.Port))
	connectionString := fmt.Sprintf("%s://%s%s", config.MongoDB.Scheme, hosts, config.MongoDB.Port)

	tM := reflect.TypeOf(bson.M{})
	reg := bson.NewRegistryBuilder().RegisterTypeMapEntry(bsontype.EmbeddedDocument, tM).Build()
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString).SetRegistry(reg))
	if err != nil {
		logger.Fatalf("%s", err)
	}

	database := client.Database(config.MongoDB.Database)
	_, err = collectCreateIndex(database, ctx, logger)
	if err != nil {
		logger.Fatalf("%s", err)
	}

	return &WTF{
		Client: database,
		Ctx:    ctx,
		Logger: logger,
		Config: config,
	}
}

func (db *WTF) GetGenerals(collectionName string) (*mongo.Cursor, error) {
	collection := db.Client.Collection(collectionName)
	findOptions := options.Find()
	findOptions.SetProjection(bson.D{{Key: "general", Value: 1}})

	cur, err := collection.Find(db.Ctx, bson.D{{}}, findOptions)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("not found")
		} else {
			return nil, errors.New("unknown db error")
		}
	}

	return cur, err
}

func indexExists(ctx context.Context, collection *mongo.Collection, indexName string) (bool, error) {
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return false, fmt.Errorf("indexExists: %w", err)
	}
	defer cursor.Close(ctx)
	var indexes []bson.M
	if err = cursor.All(ctx, &indexes); err != nil {
		return false, err
	}
	for _, index := range indexes {
		if index["name"] == indexName {
			return true, nil
		}
	}
	return false, nil
}

func collectCreateIndex(database *mongo.Database, ctx context.Context, logger *logrus.Logger) (interface{}, error) {
	indices := map[string]mongo.IndexModel{
		"user":       {Keys: bson.M{"username": 1}, Options: options.Index().SetUnique(true).SetName("username_1")},
		"service":    {Keys: bson.M{"name": 1}, Options: options.Index().SetUnique(true).SetName("name_1")},
		"clusters":   {Keys: bson.M{"general.name": 1}, Options: options.Index().SetUnique(true).SetName("general_name_1")},
		"listeners":  {Keys: bson.M{"general.name": 1}, Options: options.Index().SetUnique(true).SetName("general_name_1")},
		"endpoints":  {Keys: bson.M{"general.name": 1}, Options: options.Index().SetUnique(true).SetName("general_name_1")},
		"routes":     {Keys: bson.M{"general.name": 1}, Options: options.Index().SetUnique(true).SetName("general_name_1")},
		"extensions": {Keys: bson.M{"general.name": 1}, Options: options.Index().SetUnique(true).SetName("general_name_1")},
		"vhds":       {Keys: bson.M{"general.name": 1}, Options: options.Index().SetUnique(true).SetName("general_name_1")},
		"bootstrap":  {Keys: bson.M{"general.name": 1}, Options: options.Index().SetUnique(true).SetName("general_name_1")},
	}

	for collectionName, index := range indices {
		collection := database.Collection(collectionName)
		indexName := getIndexName(index)
		if err := createIndex(ctx, collection, index, indexName); err != nil {
			logger.Fatalf("Failed to create index for %s: %v\n", collectionName, err)
		}
	}

	return nil, nil
}

func createIndex(ctx context.Context, collection *mongo.Collection, index mongo.IndexModel, indexName string) error {
	if indexName == "" {
		return errors.New("invalid index name")
	}
	exists, err := indexExists(ctx, collection, indexName)
	if err != nil {
		return fmt.Errorf("could not check for index existence: %v", err)
	}

	if !exists {
		_, err = collection.Indexes().CreateOne(ctx, index)
		if err != nil {
			return fmt.Errorf("could not create index for %v on collection %v: %w", index.Keys, collection.Name(), err)
		}
	}
	return nil
}

func getIndexName(index mongo.IndexModel) string {
	keys, ok := index.Keys.(bson.M)
	if !ok {
		return ""
	}

	nameParts := make([]string, 0, len(keys))

	for key, val := range keys {
		if nestedKeys, ok := val.(bson.M); ok {
			for nestedKey := range nestedKeys {
				nameParts = append(nameParts, key+"."+nestedKey+"_1")
			}
		} else {
			nameParts = append(nameParts, key+"_1")
		}
	}

	return strings.Join(nameParts, "_")
}
