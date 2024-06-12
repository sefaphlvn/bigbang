package extension

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
	"github.com/sefaphlvn/bigbang/rest/poker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (extension *DBHandler) UpdateExtensions(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	filter := bson.M{"general.name": resourceDetails.Name, "general.canonical_name": resourceDetails.CanonicalName}
	filterWithRestriction := common.AddUserFilter(resourceDetails, filter)

	version, _ := strconv.Atoi(resource.GetVersion().(string))
	resource.SetVersion(strconv.Itoa(version + 1))
	update := bson.M{
		"$set": bson.M{
			"resource.resource":        resource.GetResource(),
			"resource.version":         resource.GetVersion(),
			"general.config_discovery": resource.GetConfigDiscovery(),
			"general.updated_at":       primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	collection := extension.DB.Client.Collection("extensions")
	_, err := collection.UpdateOne(extension.DB.Ctx, filterWithRestriction, update)

	if err != nil {
		return nil, err
	}

	if resourceDetails.SaveOrPublish == "publish" {
		poker.DetectChangedResource(resource.GetGeneral().GType, resourceDetails.Name, extension.DB)
		poker.ResetProcessedResources()
	}

	return gin.H{"message": "Success"}, nil
}
