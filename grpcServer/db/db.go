package db

import (
	"context"
	"errors"
	"log"
	"reflect"
	"time"

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
	return &MongoDB{
		Client: database,
		Ctx:    context.Background(),
		Cancel: cancel,
	}, err
}

func (db *MongoDB) GetGenerals(collectionName string) (*mongo.Cursor, error) {
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
