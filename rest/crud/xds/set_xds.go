package xds

import (
	"context"
	"encoding/json"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sefaphlvn/bigbang/pkg/errstr"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"github.com/sefaphlvn/bigbang/rest/crud"
)

func (xds *AppHandler) SetResource(ctx context.Context, resource models.DBResourceClass, requestDetails models.RequestDetails) (any, error) {
	general := resource.GetGeneral()
	err := resources.PrepareResource(resource, requestDetails, xds.Context.Logger, xds.ResourceService)
	if err != nil {
		return nil, err
	}
	bootstrapID := ""
	resourceID := ""
	collection := xds.Context.Client.Collection(general.Collection)
	inserResult, err := collection.InsertOne(ctx, resource)
	if err != nil {
		if er := new(mongo.WriteException); errors.As(err, &er) && er.WriteErrors[0].Code == 11000 {
			return nil, errstr.ErrNameAlreadyExists
		}
		return nil, err
	}

	if general.GType == models.Listener {
		bootstrapID, err = xds.createBootstrap(ctx, general, requestDetails)
		if err != nil {
			return nil, err
		}

		if general.Managed {
			if err := xds.createService(ctx, general.Name); err != nil {
				return nil, err
			}
		}
	}

	if oid, ok := inserResult.InsertedID.(primitive.ObjectID); ok {
		resourceID = oid.Hex()
	}

	data := map[string]any{"bootstrap_id": bootstrapID, "resource_id": resourceID}

	return map[string]any{"message": "Success", "data": data}, nil
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

func (xds *AppHandler) createBootstrap(ctx context.Context, listenerGeneral models.General, requestDetails models.RequestDetails) (string, error) {
	collection := xds.Context.Client.Collection("bootstrap")
	bootstrap := crud.GetBootstrap(listenerGeneral, xds.Context.Config)
	resource, err := DecodeFromMap(bootstrap)
	if err != nil {
		return "", err
	}
	err = resources.PrepareResource(resource, requestDetails, xds.Context.Logger, xds.ResourceService)
	if err != nil {
		return "", err
	}

	inserResult, err := collection.InsertOne(ctx, resource)
	if err != nil {
		if er := new(mongo.WriteException); errors.As(err, &er) && er.WriteErrors[0].Code == 11000 {
			return "", errstr.ErrNameAlreadyExists
		}
		return "", err
	}

	if oid, ok := inserResult.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}

	return "", errors.New("inserted ID is not a valid ObjectID")
}

func DecodeFromMap(data map[string]any) (models.DBResourceClass, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var resource models.DBResource
	if err := json.Unmarshal(jsonData, &resource); err != nil {
		return nil, err
	}

	return &resource, nil
}
