package resources

import (
	"context"
	"errors"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/errstr"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

type GeneralResponse struct {
	General models.General `bson:"general"`
}

func GetResourceNGeneral(ctx context.Context, db *db.AppContext, collectionName, name, project string) (*models.DBResource, error) {
	var doc models.DBResource

	collection := db.Client.Collection(collectionName)
	findOptions := options.FindOne()
	findOptions.SetProjection(bson.D{{Key: "resource", Value: 1}, {Key: "_id", Value: 1}, {Key: "general", Value: 1}})

	filter := bson.D{{Key: "general.name", Value: name}, {Key: "general.project", Value: project}}

	err := collection.FindOne(ctx, filter, findOptions).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("not found: (" + name + ")")
		}
		return nil, errstr.ErrUnknownDBError
	}

	return &doc, nil
}

func IncrementResourceVersion(ctx context.Context, db *db.AppContext, name, project string) (string, error) {
	collection := db.Client.Collection("listeners")

	var doc models.DBResource
	filter := bson.D{{Key: "general.name", Value: name}, {Key: "general.project", Value: project}}
	findOptions := options.FindOne()
	findOptions.SetProjection(bson.D{{Key: "resource.version", Value: 1}})

	err := collection.FindOne(ctx, filter, findOptions).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errors.New("not found: (" + name + ")")
		}
		return "", errstr.ErrUnknownDBError
	}

	// Mevcut version değerini int'e çevir ve artır
	versionInt, err := strconv.Atoi(doc.Resource.Version)
	if err != nil {
		return "", errstr.ErrInvalidVersion
	}

	// Version'u 1 artır
	versionInt++

	// Yeni version'u string'e çevir
	newVersion := strconv.Itoa(versionInt)

	// MongoDB'de version değerini güncelle
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "resource.version", Value: newVersion}}}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", errstr.ErrFailedToUpdateVersion
	}

	return newVersion, nil
}

func GetGenerals(ctx context.Context, context *db.AppContext, collectionName string, filter primitive.D) ([]*models.General, error) {
	collection := context.Client.Collection(collectionName)

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$project", Value: bson.D{{Key: "general", Value: 1}, {Key: "_id", Value: 0}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*models.General
	for cursor.Next(ctx) {
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
