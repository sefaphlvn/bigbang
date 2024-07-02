package extension

import (
	"errors"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (extension *AppHandler) GetExtension(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	collection := extension.Context.Client.Collection("extensions")
	filter := bson.M{"general.name": resourceDetails.Name, "general.canonical_name": resourceDetails.CanonicalName}
	filterWithRestriction := common.AddUserFilter(resourceDetails, filter)
	result := collection.FindOne(extension.Context.Ctx, filterWithRestriction)

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("not found: (" + resourceDetails.Name + ")")
		} else {
			return nil, errors.New("unknown db error")
		}
	}

	err := result.Decode(resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (extension *AppHandler) GetExtensions(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	collection := extension.Context.Client.Collection("extensions")
	filter := bson.M{"general.name": resourceDetails.Name}
	result := collection.FindOne(extension.Context.Ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("not found: (" + resourceDetails.Name + ")")
		} else {
			return nil, errors.New("unknown db error")
		}
	}
	err := result.Decode(resource)
	if err != nil {
		return nil, err
	}
	return resource, nil
}
