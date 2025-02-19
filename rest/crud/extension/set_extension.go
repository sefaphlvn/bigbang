package extension

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sefaphlvn/bigbang/pkg/errstr"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

func (extension *AppHandler) SetExtension(ctx context.Context, resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	general := resource.GetGeneral()
	resourceID := ""
	resources.PrepareResource(resource, requestDetails, extension.Context.Logger)
	collection := extension.Context.Client.Collection(general.Collection)
	inserResult, err := collection.InsertOne(ctx, resource)
	if err != nil {
		if er := new(mongo.WriteException); errors.As(err, &er) && er.WriteErrors[0].Code == 11000 {
			return nil, errstr.ErrNameAlreadyExists
		}
		return nil, err
	}

	if oid, ok := inserResult.InsertedID.(primitive.ObjectID); ok {
		resourceID = oid.Hex()
	}

	data := map[string]interface{}{"resource_id": resourceID}

	return map[string]interface{}{"message": "Success", "data": data}, nil
}
