package extension

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
	"github.com/sefaphlvn/bigbang/rest/crud/typed_configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (extension *AppHandler) SetExtension(resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	general := resource.GetGeneral()
	now := time.Now()
	general.CreatedAt = primitive.NewDateTimeFromTime(now)
	general.UpdatedAt = primitive.NewDateTimeFromTime(now)
	resource.SetGeneral(&general)
	resource.SetTypedConfig(typed_configs.DecodeSetTypedConfigs(resource, extension.Context.Logger))
	common.DetectSetPermissions(resource, requestDetails)
	
	collection := extension.Context.Client.Collection(requestDetails.Collection)
	_, err := collection.InsertOne(extension.Context.Ctx, resource)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("name already exists")
		}
		return nil, err
	}

	return gin.H{"message": "Success", "data": nil}, nil
}
