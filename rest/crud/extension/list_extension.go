package extension

import (
	"errors"

	"github.com/sefaphlvn/bigbang/rest/crud/common"
	"github.com/sefaphlvn/bigbang/rest/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (extension *DBHandler) ListExtensions(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	var records []bson.M
	collection := extension.DB.Client.Collection("extensions")

	filter := bson.M{"general.canonical_name": resourceDetails.CanonicalName}
	filterWithRestriction := common.AddUserFilter(resourceDetails, filter)

	opts := options.Find().SetProjection(bson.M{"resource": 0})

	cursor, err := collection.Find(extension.DB.Ctx, filterWithRestriction, opts)
	if err != nil {
		return nil, errors.New("unknown db error")
	}

	if err = cursor.All(extension.DB.Ctx, &records); err != nil {
		return nil, errors.New("unknown db error")
	}

	generals := common.TransformGenerals(records)

	return generals, nil
}
