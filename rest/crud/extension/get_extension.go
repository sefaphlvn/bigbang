package extension

import (
	"errors"
	"fmt"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (extension *AppHandler) GetExtensions(resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	var records []bson.M
	filter := bson.M{"general.type": requestDetails.Type, "general.project": requestDetails.Project}
	filterWithRestriction := common.AddUserFilter(requestDetails, filter)
	collection := extension.Context.Client.Collection(requestDetails.Collection)

	opts := options.Find().SetProjection(bson.M{"resource": 0})
	cursor, err := collection.Find(extension.Context.Ctx, filterWithRestriction, opts)
	if err != nil {
		return nil, fmt.Errorf("db find error: %w", err)
	}
	defer cursor.Close(extension.Context.Ctx)

	if err = cursor.All(extension.Context.Ctx, &records); err != nil {
		return nil, fmt.Errorf("cursor all error: %w", err)
	}

	generals := common.TransformGenerals(records)
	return generals, nil
}

func (extension *AppHandler) GetExtension(resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	return getExtensionByFilter(resource, extension, requestDetails, bson.M{
		"general.name":           requestDetails.Name,
		"general.canonical_name": requestDetails.CanonicalName,
		"general.project":        requestDetails.Project,
	})
}

func (extension *AppHandler) GetOtherExtension(resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	return getExtensionByFilter(resource, extension, requestDetails, bson.M{
		"general.name":    requestDetails.Name,
		"general.project": requestDetails.Project,
	})
}

func getExtensionByFilter(resource models.DBResourceClass, extension *AppHandler, requestDetails models.RequestDetails, filter bson.M) (interface{}, error) {
	collection := extension.Context.Client.Collection(requestDetails.Collection)
	filterWithRestriction := common.AddUserFilter(requestDetails, filter)
	result := collection.FindOne(extension.Context.Ctx, filterWithRestriction)

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("not found: (%s)", requestDetails.Name)
		}
		return nil, fmt.Errorf("db find one error: %w", result.Err())
	}

	if err := result.Decode(resource); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return resource, nil
}
