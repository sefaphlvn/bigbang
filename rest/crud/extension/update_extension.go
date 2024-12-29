package extension

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"github.com/sefaphlvn/bigbang/rest/crud"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
)

func (extension *AppHandler) UpdateExtensions(ctx context.Context, resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	filter := bson.M{"general.name": requestDetails.Name, "general.canonical_name": requestDetails.CanonicalName, "general.project": requestDetails.Project}
	return updateResource(ctx, extension, resource, requestDetails, filter)
}

func (extension *AppHandler) UpdateOtherExtensions(ctx context.Context, resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	filter := bson.M{"general.name": requestDetails.Name, "general.project": requestDetails.Project}
	return updateResource(ctx, extension, resource, requestDetails, filter)
}

func updateResource(ctx context.Context, extension *AppHandler, resource models.DBResourceClass, requestDetails models.RequestDetails, filter bson.M) (interface{}, error) {
	filterWithRestriction := common.AddUserFilter(requestDetails, filter)
	versionStr, ok := resource.GetVersion().(string)
	if !ok {
		extension.Context.Logger.Warnf("expected string type for version, got %v", resource.GetVersion())
		return nil, fmt.Errorf("invalid version format: %v", resource.GetVersion())
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return nil, fmt.Errorf("invalid version format: %s", versionStr)
	}
	resource.SetVersion(strconv.Itoa(version + 1))
	newResource := resource.GetResource()
	validateErr, isErr, err := resources.Validate(resource.GetGeneral().GType, newResource)
	if isErr {
		return validateErr, err
	}

	resource.SetTypedConfig(resources.DecodeSetTypedConfigs(resource, extension.Context.Logger))

	update := bson.M{
		"$set": bson.M{
			"resource.resource":        newResource,
			"resource.version":         resource.GetVersion(),
			"general.config_discovery": resource.GetConfigDiscovery(),
			"general.updated_at":       primitive.NewDateTimeFromTime(time.Now()),
			"general.typed_config":     resource.GetTypedConfig(),
		},
	}

	collection := extension.Context.Client.Collection(requestDetails.Collection)
	_, err = collection.UpdateOne(ctx, filterWithRestriction, update)
	if err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}

	project := resource.GetGeneral().Project
	changedResources := crud.HandleResourceChange(ctx, resource, requestDetails, extension.Context, project, extension.Poke)

	return gin.H{"message": "Success", "data": changedResources}, nil
}
