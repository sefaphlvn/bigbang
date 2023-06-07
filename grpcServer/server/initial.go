package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/sefaphlvn/bigbang/grpcServer/db"
	"github.com/sefaphlvn/bigbang/restServer/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitialSnapshots(db *db.MongoDB) error {
	collection := db.Client.Collection("listeners")
	findOptions := options.Find()
	findOptions.SetProjection(bson.D{{Key: "general", Value: 1}})

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("not found")
		} else {
			return errors.New("unknown db error")
		}
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		var general models.General
		bsonBytes, _ := bson.Marshal(result["general"])
		bson.Unmarshal(bsonBytes, &general)

		fmt.Println(general)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return err
}
