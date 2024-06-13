package xds

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud/typed_configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (xds *DBHandler) SetResource(resource models.DBResourceClass, collectionName models.ResourceDetails) (interface{}, error) {
	general := resource.GetGeneral()
	now := time.Now()
	general.CreatedAt = primitive.NewDateTimeFromTime(now)
	general.UpdatedAt = primitive.NewDateTimeFromTime(now)

	helper.PrettyPrinter(general)
	resource.SetGeneral(&general)
	resource.SetTypedConfig(typed_configs.DecodeSetTypedConfigs(resource, xds.DB.Logger))

	collection := xds.DB.Client.Collection(collectionName.Type.String())
	_, err := collection.InsertOne(xds.DB.Ctx, resource)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("name already exists")
		}
		return nil, err
	}

	if general.Type == "listeners" {
		if general.Service.Enabled {
			xds.createService(general.Name)
		}
	}

	return gin.H{"message": "Success"}, nil
}

func (xds *DBHandler) createService(serviceName string) (interface{}, error) {
	var service models.Service
	collection := xds.DB.Client.Collection("service")
	service.Name = serviceName
	_, err := collection.InsertOne(xds.DB.Ctx, service)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("name already exists")
		}
		return nil, err
	}

	return nil, nil
}
