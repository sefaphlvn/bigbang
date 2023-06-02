package xds

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/restapi/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (xds *DBHandler) DelResource(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	collection := xds.DB.Client.Collection(resourceDetails.Type)
	filter := bson.M{"general.name": resourceDetails.Name}
	res, err := collection.DeleteOne(xds.DB.Ctx, filter)

	if err != nil {
		return nil, errors.New("unknown db error")
	}

	if res.DeletedCount == 0 {
		return nil, errors.New("document not found")
	}

	return gin.H{"message": "Success"}, nil
}
