package resources

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type GeneralResponse struct {
	General models.General `bson:"general"`
}

var unmarshaler = protojson.UnmarshalOptions{
	AllowPartial:   true,
	DiscardUnknown: true,
}

func GetResource(db *db.AppContext, collectionName string, name string) (*models.DBResource, error) {
	var doc models.DBResource

	collection := db.Client.Collection(collectionName)
	findOptions := options.FindOne()
	findOptions.SetProjection(bson.D{{Key: "resource", Value: 1}, {Key: "_id", Value: 0}, {Key: "general", Value: 1}})

	filter := bson.D{{Key: "general.name", Value: name}}

	err := collection.FindOne(db.Ctx, filter, findOptions).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("not found: (" + name + ")")
		} else {
			return nil, errors.New("unknown db error")
		}
	}

	return &doc, nil
}

func GetGenerals(context *db.AppContext, collectionName string, filter primitive.D) ([]*models.General, error) {
	collection := context.Client.Collection(collectionName)

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$project", Value: bson.D{{Key: "general", Value: 1}, {Key: "_id", Value: 0}}}},
	}

	cursor, err := collection.Aggregate(context.Ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Ctx)

	var results []*models.General
	for cursor.Next(context.Ctx) {
		var resp GeneralResponse
		if err = cursor.Decode(&resp); err != nil {
			context.Logger.Debug(err)
			return nil, err
		}
		results = append(results, &resp.General)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func MarshalUnmarshalWithType(data interface{}, msg proto.Message) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = unmarshaler.Unmarshal(jsonData, msg)
	if err != nil {
		fmt.Println("proto unmarshall error: ", err)
		return err
	}

	return nil
}
