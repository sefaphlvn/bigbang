package extension

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sefaphlvn/bigbang/pkg/errstr"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
)

func (extension *AppHandler) ListExtensions(ctx context.Context, _ models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	var records []bson.M
	collection := extension.Context.Client.Collection(requestDetails.Collection)
	filter := bson.M{"general.canonical_name": requestDetails.CanonicalName, "general.project": requestDetails.Project}
	filterWithRestriction := common.AddUserFilter(requestDetails, filter)

	opts := options.Find().SetProjection(bson.M{"resource": 0})

	cursor, err := collection.Find(ctx, filterWithRestriction, opts)
	if err != nil {
		return nil, errstr.ErrUnknownDBError
	}

	if err = cursor.All(ctx, &records); err != nil {
		return nil, errstr.ErrUnknownDBError
	}

	generals := common.TransformGenerals(records)

	return generals, nil
}
