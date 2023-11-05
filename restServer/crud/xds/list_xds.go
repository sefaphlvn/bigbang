package xds

import (
	"fmt"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/restServer/models"
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

	filter := bson.M{}
	if !resourceDetails.User.IsAdmin {
		filter = bson.M{
			"general.groups": bson.M{
				"$in": resourceDetails.User.Groups,
			},
		}
	}

	cursor, err := collection.Find(xds.DB.Ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("could not find records: %v", err)
	}

	var records []bson.M
	if err = cursor.All(xds.DB.Ctx, &records); err != nil {
		return nil, fmt.Errorf("could not decode records: %v", err)
	}

	var schema = ResourceSchema{
		"listeners": []Field{
			{Name: "agent", Type: "bool"},
			{Name: "team", Type: "string"},
			{Name: "service", Type: "string"},
		},
	}

	return transformRecords(records, resourceDetails.Type, schema), nil
}

func transformRecords(records []bson.M, resourceType string, schema ResourceSchema) interface{} {
	var generals []models.General
	fields, ok := schema[resourceType]
	if !ok {
		fields = []Field{}
	}

	for _, record := range records {
		general, ok := record["general"].(bson.M)
		if !ok {
			continue
		}
		extra, exOK := general["extra"].(bson.M)

		g := models.General{
			Name:      helper.GetString(general, "name"),
			Version:   helper.GetString(general, "version"),
			Type:      helper.GetString(general, "type"),
			SubType:   helper.GetString(general, "subtype"),
			CreatedAt: helper.GetDateTime(general, "created_at"),
			UpdatedAt: helper.GetDateTime(general, "updated_at"),
			Extra:     map[string]interface{}{},
		}

		if exOK {
			for _, field := range fields {
				switch field.Type {
				case "string":
					g.Extra[field.Name] = helper.GetString(extra, field.Name)
				case "bool":
					g.Extra[field.Name] = helper.GetBool(extra, field.Name)
				}
			}
		}

		generals = append(generals, g)
	}
	return generals
}
