package xds

import (
	"fmt"

	"github.com/sefaphlvn/bigbang/rest/crud/common"
	"github.com/sefaphlvn/bigbang/rest/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Field struct {
	Name string
	Type string
}

type ResourceSchema map[string][]Field

func (xds *DBHandler) ListResource(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	collection := xds.DB.Client.Collection(resourceDetails.Type)
	opts := options.Find().SetProjection(bson.M{"resource": 0})

	filterWithRestriction := common.AddUserFilter(resourceDetails, bson.M{})
	cursor, err := collection.Find(xds.DB.Ctx, filterWithRestriction, opts)
	if err != nil {
		return nil, fmt.Errorf("could not find records: %v", err)
	}

	var records []bson.M
	if err = cursor.All(xds.DB.Ctx, &records); err != nil {
		return nil, fmt.Errorf("could not decode records: %v", err)
	}

	return common.TransformGenerals(records), nil
}
