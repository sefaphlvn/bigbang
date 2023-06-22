package db

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Database
	Ctx    context.Context
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

func collectCreateIndex(database *mongo.Database, ctx context.Context) (interface{}, error) {
	opt := options.Index()
	opt.SetUnique(true)

	indices := map[string]mongo.IndexModel{
		"user":         {Keys: bson.M{"username": 1}, Options: opt},
		"service":      {Keys: bson.M{"name": 1}, Options: opt},
		"clusters":     {Keys: bson.M{"general.name": 1}, Options: opt},
		"listeners":    {Keys: bson.M{"general.name": 1}, Options: opt},
		"endpoints":    {Keys: bson.M{"general.name": 1}, Options: opt},
		"routes":       {Keys: bson.M{"general.name": 1}, Options: opt},
		"lb_endpoints": {Keys: bson.M{"general.name": 1}, Options: opt},
		"extensions":   {Keys: bson.M{"general.name": 1}, Options: opt},
	}

	for collectionName, index := range indices {
		collection := database.Collection(collectionName)
		indexName := getIndexName(index)
		if err := createIndex(ctx, collection, index, indexName); err != nil {
			fmt.Printf("Failed to create index for %s: %v\n", collectionName, err)
		}
	}

	return nil, nil
}

func NewMongoDB(uri string) (*MongoDB, error) {
	tM := reflect.TypeOf(bson.M{})
	reg := bson.NewRegistryBuilder().RegisterTypeMapEntry(bsontype.EmbeddedDocument, tM).Build()
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetRegistry(reg))
	if err != nil {
		return nil, err
	}

	database := client.Database("navigazer")
	_, err = collectCreateIndex(database, ctx)
	if err != nil {
		fmt.Println(err)
	}

	return &MongoDB{
		Client: database,
		Ctx:    ctx,
	}, nil
}
