package xds

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (xds *AppHandler) DelResource(resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	resourceType := requestDetails.Collection
	collection := xds.Context.Client.Collection(resourceType)

	filter := buildFilter(requestDetails)

	if err := checkDocumentExists(xds, collection, filter); err != nil {
		return nil, err
	}

	if err := deleteDocument(xds, collection, filter); err != nil {
		return nil, err
	}

	if resourceType == "listeners" {
		if err := xds.delBootstrap(filter); err != nil {
			return nil, err
		}
	}

	return gin.H{"message": "Success"}, nil
}

func (xds *AppHandler) delBootstrap(filter primitive.M) error {
	collection := xds.Context.Client.Collection("bootstrap")

	if err := checkDocumentExists(xds, collection, filter); err != nil {
		return err
	}

	if err := deleteDocument(xds, collection, filter); err != nil {
		return err
	}

	return nil
}

func buildFilter(requestDetails models.RequestDetails) bson.M {
	if requestDetails.User.IsOwner {
		return bson.M{"general.name": requestDetails.Name}
	}
	return bson.M{
		"general.name": requestDetails.Name,
		"general.groups": bson.M{
			"$in": requestDetails.User.Groups,
		},
	}
}

func checkDocumentExists(xds *AppHandler, collection *mongo.Collection, filter bson.M) error {
	result := collection.FindOne(xds.Context.Ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return errors.New("document not found or no permission to delete")
		}
		return errors.New("unknown db error")
	}
	return nil
}

func deleteDocument(xds *AppHandler, collection *mongo.Collection, filter bson.M) error {
	res, err := collection.DeleteOne(xds.Context.Ctx, filter)
	if err != nil {
		return errors.New("unknown db error")
	}

	if res.DeletedCount == 0 {
		return errors.New("document not found")
	}

	return nil
}
