package xds

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sefaphlvn/bigbang/pkg/errstr"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
)

func (xds *AppHandler) DelResource(_ models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	resourceType := requestDetails.Collection
	collection := xds.Context.Client.Collection(resourceType)

	dependList := common.IsDeletable(xds.Context, requestDetails.GType, requestDetails.Name)
	if len(dependList) > 0 {
		message := "Resource has dependencies: \n " + strings.Join(dependList, ", ")
		return nil, errors.New(message)
	}

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
		return bson.M{"general.name": requestDetails.Name, "general.project": requestDetails.Project}
	}
	return bson.M{
		"general.name":    requestDetails.Name,
		"general.project": requestDetails.Project,
		"general.groups": bson.M{
			"$in": requestDetails.User.Groups,
		},
	}
}

func checkDocumentExists(xds *AppHandler, collection *mongo.Collection, filter bson.M) error {
	result := collection.FindOne(xds.Context.Ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return errstr.ErrNoDocumentsDelete
		}
		return errstr.ErrUnknownDBError
	}
	return nil
}

func deleteDocument(xds *AppHandler, collection *mongo.Collection, filter bson.M) error {
	res, err := collection.DeleteOne(xds.Context.Ctx, filter)
	if err != nil {
		return errstr.ErrUnknownDBError
	}

	if res.DeletedCount == 0 {
		return errstr.ErrNoDocuments
	}

	return nil
}
