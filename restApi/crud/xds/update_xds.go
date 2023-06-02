package xds

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/restApi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (xds *DBHandler) UpdateResource(resource models.DBResourceClass, collectionName models.ResourceDetails) (interface{}, error) {
	filter := bson.M{"general.name": collectionName.Name}
	version, _ := strconv.Atoi(resource.GetVersion().(string))
	resource.SetVersion(strconv.Itoa(version + 1))
	update := bson.M{
		"$set": bson.M{
			"resource.resource":  resource.GetResource(),
			"resource.version":   resource.GetVersion(),
			"general.updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	collection := xds.DB.Client.Collection(collectionName.Type)
	_, err := collection.UpdateOne(xds.DB.Ctx, filter, update)

	if err != nil {
		return nil, err
	}
	return gin.H{"message": "Success"}, nil
}
