package db

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/sefaphlvn/bigbang/restServer/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Database
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewMongoDB(uri string) (*MongoDB, error) {
	tM := reflect.TypeOf(bson.M{})
	reg := bson.NewRegistryBuilder().RegisterTypeMapEntry(bsontype.EmbeddedDocument, tM).Build()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetRegistry(reg))
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database("navigazer")
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"general.name": 1}, Options: opt}

	for _, collectionName := range helper.Collections {
		collection := database.Collection(collectionName)
		if _, err := collection.Indexes().CreateOne(ctx, index); err != nil {
			log.Fatalf("could not create index for name on collection %s: %v", collectionName, err)
		}
	}

	userIndex := mongo.IndexModel{Keys: bson.M{"username": 1}, Options: opt}
	collection := database.Collection("user")
	if _, err := collection.Indexes().CreateOne(ctx, userIndex); err != nil {
		log.Fatalf("could not create index for username on collection %s: %v", "user", err)
	}

	return &MongoDB{
		Client: database,
		Ctx:    context.Background(),
		Cancel: cancel,
	}, err
}
