package resources

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sirupsen/logrus"
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

func GetResourceNGeneral(db *db.AppContext, collectionName string, name string, project string) (*models.DBResource, error) {
	var doc models.DBResource

	collection := db.Client.Collection(collectionName)
	findOptions := options.FindOne()
	findOptions.SetProjection(bson.D{{Key: "resource", Value: 1}, {Key: "_id", Value: 1}, {Key: "general", Value: 1}})

	filter := bson.D{{Key: "general.name", Value: name}, {Key: "general.project", Value: project}}

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

func IncrementResourceVersion(db *db.AppContext, name string, project string) (string, error) {
	collection := db.Client.Collection("listeners")

	var doc models.DBResource
	filter := bson.D{{Key: "general.name", Value: name}, {Key: "general.project", Value: project}}
	findOptions := options.FindOne()
	findOptions.SetProjection(bson.D{{Key: "resource.version", Value: 1}})

	err := collection.FindOne(db.Ctx, filter, findOptions).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errors.New("not found: (" + name + ")")
		}
		return "", errors.New("unknown db error")
	}

	// Mevcut version değerini int'e çevir ve artır
	versionInt, err := strconv.Atoi(doc.Resource.Version)
	if err != nil {
		return "", errors.New("invalid version format")
	}

	// Version'u 1 artır
	versionInt++

	// Yeni version'u string'e çevir
	newVersion := strconv.Itoa(versionInt)

	// MongoDB'de version değerini güncelle
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "resource.version", Value: newVersion}}}}

	_, err = collection.UpdateOne(db.Ctx, filter, update)
	if err != nil {
		return "", errors.New("failed to update resource version")
	}

	return newVersion, nil
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

func ConvertToJSON(v interface{}, log *logrus.Logger) string {
	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Infof("JSON convert err: %v", err)
	}
	return string(jsonData)
}
