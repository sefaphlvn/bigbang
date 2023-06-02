package extension

import (
	"errors"

	"github.com/sefaphlvn/bigbang/restapi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (extension *DBHandler) ListExtensions(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	var records []bson.M
	collection := extension.DB.Client.Collection("extensions")
	opts := options.Find().SetProjection(bson.M{"resource": 0})

	cursor, err := collection.Find(extension.DB.Ctx, bson.M{}, opts)
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
			Name:      general["name"].(string),
			Version:   general["version"].(string),
			Type:      general["type"].(string),
			SubType:   general["subtype"].(string),
			CreatedAt: general["created_at"].(primitive.DateTime),
			UpdatedAt: general["updated_at"].(primitive.DateTime),
		}

		generals = append(generals, g)
	}

	return generals, nil
}
