package extension

import (
	"errors"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (extension *AppHandler) ListExtensions(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	var records []bson.M
	collection := extension.Context.Client.Collection("extensions")

	filter := bson.M{"general.canonical_name": resourceDetails.CanonicalName}
	filterWithRestriction := common.AddUserFilter(resourceDetails, filter)

	opts := options.Find().SetProjection(bson.M{"resource": 0})

	cursor, err := collection.Find(extension.Context.Ctx, filterWithRestriction, opts)
	if err != nil {
		return nil, errors.New("unknown db error")
	}

	if err = cursor.All(extension.Context.Ctx, &records); err != nil {
		return nil, errors.New("unknown db error")
	}

	generals := common.TransformGenerals(records)

	return generals, nil
}
