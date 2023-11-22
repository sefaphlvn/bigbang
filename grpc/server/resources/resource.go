package resources

import (
	"errors"
	"fmt"

	"github.com/sefaphlvn/bigbang/grpc/models"

	"github.com/sefaphlvn/bigbang/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetResource(db *db.MongoDB, collectionName string, name string) (*models.DBResource, error) {
	var doc models.DBResource

	collection := db.Client.Collection(collectionName)
	findOptions := options.FindOne()
	findOptions.SetProjection(bson.D{{Key: "resource", Value: 1}, {Key: "_id", Value: 0}, {Key: "general", Value: 1}})

	filter := bson.D{{Key: "general.name", Value: name}}

	err := collection.FindOne(db.Ctx, filter, findOptions).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("not found2")
		} else {
			fmt.Println(err)
			return nil, errors.New("unknown db error")
		}
	}

	return &doc, nil
}
