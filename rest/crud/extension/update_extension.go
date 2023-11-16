package extension

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/rest/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (extension *DBHandler) UpdateExtensions(resource models.DBResourceClass, collectionName models.ResourceDetails) (interface{}, error) {
	filter := bson.M{"general.name": collectionName.Name, "general.canonical_name": collectionName.CanonicalName}
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

	collection := extension.DB.Client.Collection("extensions")
	_, err := collection.UpdateOne(extension.DB.Ctx, filter, update)

	if err != nil {
		return nil, err
	}
	return gin.H{"message": "Success"}, nil
}
