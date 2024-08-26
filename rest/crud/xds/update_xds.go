package xds

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
	"github.com/sefaphlvn/bigbang/rest/crud/typed_configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (xds *AppHandler) UpdateResource(resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	filter := bson.M{"general.name": requestDetails.Name}
	filterWithRestriction := common.AddUserFilter(requestDetails, filter)
	result := xds.Context.Client.Collection(requestDetails.Collection).FindOne(xds.Context.Ctx, filterWithRestriction)

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("document not found or no permission to update")
		} else {
			return nil, errors.New("unknown db error")
		}
	}

	version, _ := strconv.Atoi(resource.GetVersion().(string))
	resource.SetVersion(strconv.Itoa(version + 1))
	resource.SetTypedConfig(typed_configs.DecodeSetTypedConfigs(resource, xds.Context.Logger))

	update := bson.M{
		"$set": bson.M{
			"resource.resource":        resource.GetResource(),
			"resource.version":         resource.GetVersion(),
			"general.config_discovery": resource.GetConfigDiscovery(),
			"general.updated_at":       primitive.NewDateTimeFromTime(time.Now()),
			"general.typed_config":     resource.GetTypedConfig(),
		},
	}

	collection := xds.Context.Client.Collection(requestDetails.Collection)
	_, err := collection.UpdateOne(xds.Context.Ctx, filterWithRestriction, update)
	if err != nil {
		return nil, err
	}

	changedResources := crud.HandleResourceChange(resource, requestDetails, xds.Context)

	return gin.H{"message": "Success", "data": changedResources}, nil
}
