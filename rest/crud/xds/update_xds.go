package xds

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
	"github.com/sefaphlvn/bigbang/rest/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (xds *DBHandler) UpdateResource(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	filter := bson.M{"general.name": resourceDetails.Name}
	filterWithRestriction := common.AddUserFilter(resourceDetails, filter)
	result := xds.DB.Client.Collection(resourceDetails.Type).FindOne(xds.DB.Ctx, filterWithRestriction)

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("document not found or no permission to update")
		} else {
			return nil, errors.New("unknown db error")
		}
	}

	version, _ := strconv.Atoi(resource.GetVersion().(string))
	resource.SetVersion(strconv.Itoa(version + 1))
	update := bson.M{
		"$set": bson.M{
			"resource.resource":            resource.GetResource(),
			"resource.version":             resource.GetVersion(),
			"general.additional_resources": resource.GetAdditionalResources(),
			"general.updated_at":           primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	collection := xds.DB.Client.Collection(resourceDetails.Type)
	_, err := collection.UpdateOne(xds.DB.Ctx, filterWithRestriction, update)

	if err != nil {
		return nil, err
	}
	return gin.H{"message": "Success"}, nil
}
