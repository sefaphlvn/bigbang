package custom

import (
	"errors"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (custom *AppHandler) GetCustomHttpFilterList(resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	collection := custom.Context.Client.Collection(requestDetails.Collection)
	opts := options.Find()
	opts.SetProjection(bson.M{
		"general.name":           1,
		"general.canonical_name": 1,
		"general.gtype":          1,
		"general.type":           1,
		"general.category":       1,
	})

	var filters = bson.M{
		"general.version":              requestDetails.Version,
		"general.project":              requestDetails.Project,
		"general.category":             requestDetails.Category,
		"general.metadata.http_filter": bson.M{"$regex": requestDetails.Metadata["http_filter"], "$options": "i"},
	}

	cursor, err := collection.Find(custom.Context.Ctx, filters, opts)
	if err != nil {
		return nil, errors.New("unknown db error")
	}

	var results []Record
	for cursor.Next(custom.Context.Ctx) {
		var doc struct {
			General struct {
				Name          string `bson:"name"`
				CanonicalName string `bson:"canonical_name"`
				GType         string `bson:"gtype"`
				Type          string `bson:"type"`
				Category      string `bson:"category"`
				Collection    string `bson:"collection"`
			} `bson:"general"`
		}

		if err := cursor.Decode(&doc); err != nil {
			custom.Context.Logger.Debugf("Error decoding http filter: %v", err)
		}

		results = append(
			results,
			Record{
				Name:          doc.General.Name,
				CanonicalName: doc.General.CanonicalName,
				GType:         doc.General.GType,
				Type:          doc.General.Type,
				Category:      doc.General.Category,
				Collection:    requestDetails.Collection,
			},
		)
	}

	if err := cursor.Err(); err != nil {
		custom.Context.Logger.Debug(err)
	}

	return results, nil
}
