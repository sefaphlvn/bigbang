package xds

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
	"github.com/sefaphlvn/bigbang/rest/crud/typed_configs"
	"github.com/sefaphlvn/bigbang/rest/poker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (xds *DBHandler) UpdateResource(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	filter := bson.M{"general.name": resourceDetails.Name}
	filterWithRestriction := common.AddUserFilter(resourceDetails, filter)
	result := xds.DB.Client.Collection(resourceDetails.Type.String()).FindOne(xds.DB.Ctx, filterWithRestriction)

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("document not found or no permission to update")
		} else {
			return nil, errors.New("unknown db error")
		}
	}

	version, _ := strconv.Atoi(resource.GetVersion().(string))
	resource.SetVersion(strconv.Itoa(version + 1))
	resource.SetTypedConfig(typed_configs.DecodeSetTypedConfigs(resource, xds.DB.Logger))

	update := bson.M{
		"$set": bson.M{
			"resource.resource":        resource.GetResource(),
			"resource.version":         resource.GetVersion(),
			"general.config_discovery": resource.GetConfigDiscovery(),
			"general.updated_at":       primitive.NewDateTimeFromTime(time.Now()),
			"general.typed_config":     resource.GetTypedConfig(),
		},
	}

	collection := xds.DB.Client.Collection(resourceDetails.Type.String())
	_, err := collection.UpdateOne(xds.DB.Ctx, filterWithRestriction, update)
	if err != nil {
		return nil, err
	}

	if resourceDetails.SaveOrPublish == "publish" {
		poker.DetectChangedResource(resource.GetGeneral().GType, resourceDetails.Name, xds.DB)
		poker.ResetProcessedResources()
	}

	return gin.H{"message": "Success"}, nil
}
