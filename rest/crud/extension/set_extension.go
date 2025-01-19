package extension

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sefaphlvn/bigbang/pkg/errstr"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

func (extension *AppHandler) SetExtension(ctx context.Context, resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	general := resource.GetGeneral()

	resources.PrepareResource(resource, requestDetails, extension.Context.Logger)
	collection := extension.Context.Client.Collection(general.Collection)
	_, err := collection.InsertOne(ctx, resource)
	if err != nil {
		if er := new(mongo.WriteException); errors.As(err, &er) && er.WriteErrors[0].Code == 11000 {
			return nil, errstr.ErrNameAlreadyExists
		}
		return nil, err
	}

	return map[string]interface{}{"message": "Success", "data": nil}, nil
}
