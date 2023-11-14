package xds

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/rest/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (xds *DBHandler) DelResource(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
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
			return nil, errors.New("document not found or no permission to delete")
		} else {
			return nil, errors.New("unknown db error")
		}
	}

	res, err := collection.DeleteOne(xds.DB.Ctx, filter)
	if err != nil {
		return nil, errors.New("unknown db error")
	}

	if res.DeletedCount == 0 {
		return nil, errors.New("document not found")
	}

	return gin.H{"message": "Success"}, nil
}
