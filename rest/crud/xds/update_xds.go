package xds

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sefaphlvn/bigbang/pkg/errstr"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
	"github.com/sefaphlvn/bigbang/rest/crud/typedconfigs"
)

func (xds *AppHandler) UpdateResource(ctx context.Context, resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	filter := bson.M{"general.name": requestDetails.Name}
	filterWithRestriction := common.AddUserFilter(requestDetails, filter)
	result := xds.Context.Client.Collection(requestDetails.Collection).FindOne(ctx, filterWithRestriction)

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, errstr.ErrNoDocumentsUpdate
		}
		return nil, errstr.ErrUnknownDBError
	}

	newResource := resource.GetResource()
	version, _ := strconv.Atoi(resource.GetVersion().(string))
	validateErr, isErr, err := crud.Validate(resource.GetGeneral().GType, newResource)
	//	validateErr, err, isErr := crud.Validate(models.GTypes(resource.GetGeneral().GType), newResource)
	if isErr {
		return validateErr, err
	}

	resource.SetVersion(strconv.Itoa(version + 1))
	resource.SetTypedConfig(typedconfigs.DecodeSetTypedConfigs(resource, xds.Context.Logger))

	update := bson.M{
		"$set": bson.M{
			"resource.resource":        newResource,
			"resource.version":         resource.GetVersion(),
			"general.config_discovery": resource.GetConfigDiscovery(),
			"general.updated_at":       primitive.NewDateTimeFromTime(time.Now()),
			"general.typed_config":     resource.GetTypedConfig(),
		},
	}

	collection := xds.Context.Client.Collection(requestDetails.Collection)
	_, err = collection.UpdateOne(ctx, filterWithRestriction, update)
	if err != nil {
		return nil, err
	}

	project := resource.GetGeneral().Project
	changedResources := crud.HandleResourceChange(ctx, resource, requestDetails, xds.Context, project, xds.Poke)

	return gin.H{"message": "Success", "data": changedResources}, nil
}
