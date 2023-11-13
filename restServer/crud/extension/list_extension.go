package extension

import (
	"errors"

	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/restServer/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (extension *DBHandler) ListExtensions(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	var records []bson.M
	collection := extension.DB.Client.Collection("extensions")
	filter := bson.M{"general.canonical_name": resourceDetails.CanonicalName}
	opts := options.Find().SetProjection(bson.M{"resource": 0})

	cursor, err := collection.Find(extension.DB.Ctx, filter, opts)
	if err != nil {
		return nil, errors.New("unknown db error")
	}

	if err = cursor.All(extension.DB.Ctx, &records); err != nil {
		return nil, errors.New("unknown db error")
	}

	var generals []models.General
	for _, record := range records {
		general := record["general"].(bson.M)
		g := models.General{
			Name:          helper.GetString(general, "name"),
			Version:       helper.GetString(general, "version"),
			Type:          helper.GetString(general, "type"),
			GType:         helper.GetString(general, "gtype"),
			CanonicalName: helper.GetString(general, "canonical_name"),
			CreatedAt:     helper.GetDateTime(general, "created_at"),
			UpdatedAt:     helper.GetDateTime(general, "updated_at"),
		}

		generals = append(generals, g)
	}

	return generals, nil
}
