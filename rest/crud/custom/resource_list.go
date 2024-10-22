package custom

import (
	"context"
	"fmt"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Record struct {
	Name          string `json:"name" bson:"name"`
	CanonicalName string `json:"canonical_name" bson:"canonical_name"`
	GType         string `json:"gtype" bson:"gtype"`
	Type          string `json:"type" bson:"type"`
	Category      string `json:"category" bson:"category"`
	Collection    string `json:"collection" bson:"collection"`
}

func (custom *AppHandler) GetCustomResourceList(resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	collection := custom.Context.Client.Collection(requestDetails.Collection)

	opts := options.Find().SetProjection(bson.M{
		"general.name":           1,
		"general.canonical_name": 1,
		"general.gtype":          1,
		"general.type":           1,
		"general.category":       1,
	})

	filters := buildFilters(requestDetails)

	cursor, err := collection.Find(custom.Context.Ctx, filters, opts)
	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}
	defer cursor.Close(custom.Context.Ctx)

	results, decodeErr := decodeResults(cursor, requestDetails.Collection, custom.Context.Logger)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return results, nil
}

func buildFilters(details models.RequestDetails) bson.M {
	filters := bson.M{
		"general.version": details.Version,
		"general.project": details.Project,
	}

	if details.GType != "" {
		filters["general.gtype"] = details.GType
	}

	if details.Category != "" {
		filters["general.category"] = details.Category
	}

	if details.CanonicalName != "" {
		filters["general.canonical_name"] = details.CanonicalName
	}

	return filters
}

func decodeResults(cursor *mongo.Cursor, collectionName string, logger *logrus.Logger) ([]Record, error) {
	var results []Record

	for cursor.Next(context.TODO()) {
		var doc struct {
			General struct {
				Name          string `bson:"name"`
				CanonicalName string `bson:"canonical_name"`
				GType         string `bson:"gtype"`
				Type          string `bson:"type"`
				Category      string `bson:"category"`
			} `bson:"general"`
		}

		if err := cursor.Decode(&doc); err != nil {
			logger.Debugf("Decode fail: %v", err)
			continue
		}

		results = append(results, Record{
			Name:          doc.General.Name,
			CanonicalName: doc.General.CanonicalName,
			GType:         doc.General.GType,
			Type:          doc.General.Type,
			Category:      doc.General.Category,
			Collection:    collectionName,
		})
	}

	if err := cursor.Err(); err != nil {
		logger.Debugf("Cursor error: %v", err)
		return nil, err
	}

	return results, nil
}
