package xds

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
	"github.com/sefaphlvn/bigbang/rest/crud/typedConfigs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (xds *AppHandler) SetResource(resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	general := resource.GetGeneral()
	now := time.Now()
	general.CreatedAt = primitive.NewDateTimeFromTime(now)
	general.UpdatedAt = primitive.NewDateTimeFromTime(now)

	resource.SetGeneral(&general)
	validateErr, err, isErr := crud.Validate(models.GTypes(resource.GetGeneral().GType), resource.GetResource())
	if isErr {
		return validateErr, err
	}

	resource.SetTypedConfig(typedConfigs.DecodeSetTypedConfigs(resource, xds.Context.Logger))
	common.DetectSetPermissions(resource, requestDetails)

	collection := xds.Context.Client.Collection(requestDetails.Collection)
	_, err = collection.InsertOne(xds.Context.Ctx, resource)
	if err != nil {
		if er := new(mongo.WriteException); errors.As(err, &er) && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("name already exists")
		}
		return nil, err
	}

	if general.GType == models.Listener {
		if err := xds.createBootstrap(general); err != nil {
			return nil, err
		}

		if general.Managed {
			if err := xds.createService(general.Name); err != nil {
				return nil, err
			}
		}
	}

	return gin.H{"message": "Success", "data": nil}, nil
}

func (xds *AppHandler) createService(serviceName string) error {
	var service models.Service
	collection := xds.Context.Client.Collection("service")
	service.Name = serviceName
	_, err := collection.InsertOne(xds.Context.Ctx, service)
	if err != nil {
		if er := new(mongo.WriteException); errors.As(err, &er) && er.WriteErrors[0].Code == 11000 {
			return errors.New("name already exists")
		}
		return err
	}

	return nil
}

func (xds *AppHandler) createBootstrap(listenerGeneral models.General) error {
	collection := xds.Context.Client.Collection("bootstrap")
	bootstrap := crud.GetBootstrap(listenerGeneral, xds.Context.Config)
	_, err := collection.InsertOne(xds.Context.Ctx, bootstrap)
	if err != nil {
		if er := new(mongo.WriteException); errors.As(err, &er) && er.WriteErrors[0].Code == 11000 {
			return errors.New("name already exists")
		}
		return err
	}

	return nil
}
