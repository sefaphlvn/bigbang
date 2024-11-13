package extension

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
	"github.com/sefaphlvn/bigbang/rest/crud/typedconfigs"
)

func (extension *AppHandler) UpdateExtensions(resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	filter := bson.M{"general.name": requestDetails.Name, "general.canonical_name": requestDetails.CanonicalName, "general.project": requestDetails.Project}
	return updateResource(extension, resource, requestDetails, filter)
}

func (extension *AppHandler) UpdateOtherExtensions(resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	filter := bson.M{"general.name": requestDetails.Name, "general.project": requestDetails.Project}
	return updateResource(extension, resource, requestDetails, filter)
}

func updateResource(extension *AppHandler, resource models.DBResourceClass, requestDetails models.RequestDetails, filter bson.M) (interface{}, error) {
	filterWithRestriction := common.AddUserFilter(requestDetails, filter)
	versionStr := resource.GetVersion().(string)
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return nil, fmt.Errorf("invalid version format: %s", versionStr)
	}
	resource.SetVersion(strconv.Itoa(version + 1))
	newResource := resource.GetResource()
	validateErr, isErr, err := crud.Validate(resource.GetGeneral().GType, newResource)
	if isErr {
		return validateErr, err
	}

	resource.SetTypedConfig(typedconfigs.DecodeSetTypedConfigs(resource, extension.Context.Logger))

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
	_, err = collection.UpdateOne(extension.Context.Ctx, filterWithRestriction, update)
	if err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}

	project := resource.GetGeneral().Project
	changedResources := crud.HandleResourceChange(resource, requestDetails, extension.Context, project, extension.Poke)

	return gin.H{"message": "Success", "data": changedResources}, nil
}
