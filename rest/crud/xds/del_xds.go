package xds

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sefaphlvn/bigbang/pkg/errstr"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/models/downstreamfilters"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
)

func (xds *AppHandler) DelResource(ctx context.Context, _ models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	resourceType := requestDetails.Collection
	collection := xds.Context.Client.Collection(resourceType)
	downstreamFilterModel := downstreamfilters.DownstreamFilter{
		Name:    requestDetails.Name,
		Project: requestDetails.Project,
		Version: requestDetails.Version,
	}

	dependList := common.IsDeletable(ctx, xds.Context, requestDetails.GType, downstreamFilterModel)
	if len(dependList) > 0 {
		message := "Resource has dependencies: \n " + strings.Join(dependList, ", ")
		return nil, errors.New(message)
	}

	filter, err := common.AddResourceIDFilter(requestDetails, buildFilter(requestDetails))
	if err != nil {
		return nil, errors.New("invalid id format")
	}

	if err := checkDocumentExists(ctx, xds, collection, filter); err != nil {
		return nil, err
	}

	if err := deleteDocument(ctx, xds, collection, filter); err != nil {
		return nil, err
	}

	if resourceType == "listeners" {
		if err := xds.delBootstrap(ctx, filter); err != nil {
			return nil, err
		}
	}

	return gin.H{"message": "Success"}, nil
}

func (xds *AppHandler) delBootstrap(ctx context.Context, filter primitive.M) error {
	collection := xds.Context.Client.Collection("bootstrap")

	if err := checkDocumentExists(ctx, xds, collection, filter); err != nil {
		return err
	}

	if err := deleteDocument(ctx, xds, collection, filter); err != nil {
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

func checkDocumentExists(ctx context.Context, _ *AppHandler, collection *mongo.Collection, filter bson.M) error {
	result := collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return errstr.ErrNoDocumentsDelete
		}
		return errstr.ErrUnknownDBError
	}
	return nil
}

func deleteDocument(ctx context.Context, _ *AppHandler, collection *mongo.Collection, filter bson.M) error {
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return errstr.ErrUnknownDBError
	}

	if res.DeletedCount == 0 {
		return errstr.ErrNoDocuments
	}

	return nil
}
