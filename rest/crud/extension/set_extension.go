package extension

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sefaphlvn/bigbang/pkg/errstr"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
	"github.com/sefaphlvn/bigbang/rest/crud/typedconfigs"
)

func (extension *AppHandler) SetExtension(ctx context.Context, resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	general := resource.GetGeneral()
	now := time.Now()
	general.CreatedAt = primitive.NewDateTimeFromTime(now)
	general.UpdatedAt = primitive.NewDateTimeFromTime(now)
	resource.SetGeneral(&general)
	validateErr, isErr, err := crud.Validate(resource.GetGeneral().GType, resource.GetResource())
	if isErr {
		return validateErr, err
	}

	resource.SetTypedConfig(typedconfigs.DecodeSetTypedConfigs(resource, extension.Context.Logger))
	common.DetectSetPermissions(resource, requestDetails)

	collection := extension.Context.Client.Collection(requestDetails.Collection)
	_, err = collection.InsertOne(ctx, resource)
	if err != nil {
		if er := new(mongo.WriteException); errors.As(err, &er) && er.WriteErrors[0].Code == 11000 {
			return nil, errstr.ErrNameAlreadyExists
		}
		return nil, err
	}

	return gin.H{"message": "Success", "data": nil}, nil
}
