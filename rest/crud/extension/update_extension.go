package extension

import (
	"context"
	"errors"
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

func (extension *AppHandler) UpdateFilters(ctx context.Context, resource models.DBResourceClass, requestDetails models.RequestDetails) (any, error) {
	filter, err := common.AddResourceIDFilter(requestDetails, bson.M{"general.canonical_name": requestDetails.CanonicalName, "general.project": requestDetails.Project})
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	return updateResource(ctx, extension, resource, requestDetails, filter)
}

func (extension *AppHandler) UpdateExtensions(ctx context.Context, resource models.DBResourceClass, requestDetails models.RequestDetails) (any, error) {
	filter, err := common.AddResourceIDFilter(requestDetails, bson.M{"general.project": requestDetails.Project})
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	return updateResource(ctx, extension, resource, requestDetails, filter)
}

func updateResource(ctx context.Context, extension *AppHandler, resource models.DBResourceClass, requestDetails models.RequestDetails, filter bson.M) (any, error) {
	isDefault, err := common.IsDefaultResource(ctx, extension.Context, requestDetails.Name, requestDetails.Collection, requestDetails.Project)
	if err != nil {
		extension.Context.Logger.Errorf("An error occurred while checking if the resource is default: %v", err)
	} else if isDefault {
		if requestDetails.User.Role != models.RoleOwner {
			return nil, errors.New("this resource is a default resource and cannot be changed")
		}
	}

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
	nodeid := fmt.Sprintf("%s:%s", requestDetails.Name, requestDetails.Project)

	if err := resources.ValidateResourceWithClient(context.Background(), resource.GetGeneral().GType, resource.GetGeneral().Version, nodeid, newResource, extension.ResourceService); err != nil {
		return nil, fmt.Errorf("%v", err)
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
	changedResources := crud.HandleResourceChange(ctx, resource, requestDetails, extension.Context, project, extension.PokeService)

	return gin.H{"message": "Success", "data": changedResources}, nil
}
