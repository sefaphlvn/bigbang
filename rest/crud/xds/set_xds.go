package xds

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud"
	"github.com/sefaphlvn/bigbang/rest/crud/typed_configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (xds *AppHandler) SetResource(resource models.DBResourceClass, collectionName models.ResourceDetails) (interface{}, error) {
	general := resource.GetGeneral()
	now := time.Now()
	general.CreatedAt = primitive.NewDateTimeFromTime(now)
	general.UpdatedAt = primitive.NewDateTimeFromTime(now)

	helper.PrettyPrinter(general)
	resource.SetGeneral(&general)
	resource.SetTypedConfig(typed_configs.DecodeSetTypedConfigs(resource, xds.Context.Logger))

	collection := xds.Context.Client.Collection(collectionName.Type.String())
	_, err := collection.InsertOne(xds.Context.Ctx, resource)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("name already exists")
		}
		return nil, err
	}

	if general.Type == "listeners" {
		xds.createBootstrap(general.Name)
		if general.Service.Enabled {
			xds.createService(general.Name)
		}
	}

	return gin.H{"message": "Success"}, nil
}

func (xds *AppHandler) createService(serviceName string) (interface{}, error) {
	var service models.Service
	collection := xds.Context.Client.Collection("service")
	service.Name = serviceName
	_, err := collection.InsertOne(xds.Context.Ctx, service)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("name already exists")
		}
		return nil, err
	}

	return nil, nil
}

func (xds *AppHandler) createBootstrap(listenerName string) (interface{}, error) {
	collection := xds.Context.Client.Collection("bootstrap")
	bootstrap := crud.GetBootstrap(listenerName)
	_, err := collection.InsertOne(xds.Context.Ctx, bootstrap)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("name already exists")
		}
		return nil, err
	}

	return nil, nil
}
