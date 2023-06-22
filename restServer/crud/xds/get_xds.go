package xds

import (
	"errors"

	"github.com/sefaphlvn/bigbang/restServer/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (xds *DBHandler) GetResource(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	collection := xds.DB.Client.Collection(resourceDetails.Type)
	var filter bson.M
	if resourceDetails.User.IsAdmin {
		filter = bson.M{"general.name": resourceDetails.Name}
	} else {
		filter = bson.M{
			"general.name": resourceDetails.Name,
			"general.groups": bson.M{
				"$in": resourceDetails.User.Groups,
			},
		}
	}
	result := collection.FindOne(xds.DB.Ctx, filter)
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
