package xds

import (
	"fmt"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Field struct {
	Name string
	Type string
}

type ResourceSchema map[string][]Field

func (xds *DBHandler) ListResource(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	filter := bson.M{}
	collection := xds.DB.Client.Collection(resourceDetails.Type.String())
	opts := options.Find().SetProjection(bson.M{"resource": 0})

	if resourceDetails.GType != "" {
		filter = bson.M{"general.gtype": resourceDetails.GType.String()}
	}

	filterWithRestriction := common.AddUserFilter(resourceDetails, filter)
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
