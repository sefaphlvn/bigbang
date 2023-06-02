package extension

import (
	"errors"

	"github.com/sefaphlvn/bigbang/restapi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (extension *DBHandler) GetExtension(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	collection := extension.DB.Client.Collection("extensions")
	filter := bson.M{"general.name": resourceDetails.Name}
	result := collection.FindOne(extension.DB.Ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("not found")
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

func (extension *DBHandler) GetExtensions(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	collection := extension.DB.Client.Collection("extensions")
	filter := bson.M{"general.name": resourceDetails.Name}
	result := collection.FindOne(extension.DB.Ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("not found")
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
