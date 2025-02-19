package xds

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
)

type Field struct {
	Name string
	Type string
}

type ResourceSchema map[string][]Field

func (xds *AppHandler) ListResource(ctx context.Context, _ models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	filter := bson.M{}
	collection := xds.Context.Client.Collection(requestDetails.Collection)
	opts := options.Find().SetProjection(bson.M{"resource": 0})

	if requestDetails.GType != "" {
		filter["general.gtype"] = requestDetails.GType.String()
	}

	filterWithRestriction := common.AddUserFilter(requestDetails, filter)
	cursor, err := collection.Find(ctx, filterWithRestriction, opts)
	if err != nil {
		return nil, fmt.Errorf("could not find records: %w", err)
	}

	var records []bson.M
	if err = cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("could not decode records: %w", err)
	}

	return common.TransformGenerals(records), nil
}
