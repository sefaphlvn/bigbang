package xds

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sefaphlvn/bigbang/pkg/errstr"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
)

func (xds *AppHandler) GetResource(resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	collection := xds.Context.Client.Collection(requestDetails.Collection)
	filter := bson.M{"general.name": requestDetails.Name}
	filterWithRestriction := common.AddUserFilter(requestDetails, filter)
	result := collection.FindOne(xds.Context.Ctx, filterWithRestriction)

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("not found: (" + requestDetails.Name + ")")
		}
		return nil, errstr.ErrUnknownDBError
	}

	// GetSnapshotsFromServer("localhost:18000")

	err := result.Decode(resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}
