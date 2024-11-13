package xds

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

func (xds *AppHandler) SetResource(ctx context.Context, resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	general := resource.GetGeneral()
	now := time.Now()
	general.CreatedAt = primitive.NewDateTimeFromTime(now)
	general.UpdatedAt = primitive.NewDateTimeFromTime(now)

	resource.SetGeneral(&general)
	validateErr, isErr, err := crud.Validate(resource.GetGeneral().GType, resource.GetResource())
	if isErr {
		return validateErr, err
	}

	resource.SetTypedConfig(typedconfigs.DecodeSetTypedConfigs(resource, xds.Context.Logger))
	common.DetectSetPermissions(resource, requestDetails)

	collection := xds.Context.Client.Collection(requestDetails.Collection)
	_, err = collection.InsertOne(ctx, resource)
	if err != nil {
		if er := new(mongo.WriteException); errors.As(err, &er) && er.WriteErrors[0].Code == 11000 {
			return nil, errstr.ErrNameAlreadyExists
		}
		return nil, err
	}

	if general.GType == models.Listener {
		if err := xds.createBootstrap(ctx, general); err != nil {
			return nil, err
		}

		if general.Managed {
			if err := xds.createService(ctx, general.Name); err != nil {
				return nil, err
			}
		}
	}

	return gin.H{"message": "Success", "data": nil}, nil
}

func (xds *AppHandler) createService(ctx context.Context, serviceName string) error {
	var service models.Service
	collection := xds.Context.Client.Collection("service")
	service.Name = serviceName
	_, err := collection.InsertOne(ctx, service)
	if err != nil {
		if er := new(mongo.WriteException); errors.As(err, &er) && er.WriteErrors[0].Code == 11000 {
			return errstr.ErrNameAlreadyExists
		}
		return err
	}

	return nil
}

func (xds *AppHandler) createBootstrap(ctx context.Context, listenerGeneral models.General) error {
	collection := xds.Context.Client.Collection("bootstrap")
	bootstrap := crud.GetBootstrap(listenerGeneral, xds.Context.Config)
	_, err := collection.InsertOne(ctx, bootstrap)
	if err != nil {
		if er := new(mongo.WriteException); errors.As(err, &er) && er.WriteErrors[0].Code == 11000 {
			return errstr.ErrNameAlreadyExists
		}
		return err
	}

	return nil
}
