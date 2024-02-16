package custom

import (
	"errors"
	"fmt"
	"log"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Record struct {
	Name          string `json:"name" bson:"name"`
	CanonicalName string `json:"canonical_name" bson:"canonical_name"`
	GType         string `json:"gtype" bson:"gtype"`
	Category      string `json:"category" bson:"category"`
}

func (custom *DBHandler) GetCustomResourceList(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {
	collection := custom.DB.Client.Collection(resourceDetails.Collection)
	opts := options.Find()
	opts.SetProjection(bson.M{
		"general.name":           1,
		"general.canonical_name": 1,
		"general.gtype":          1,
		"general.category":       1,
	})

	var filters = bson.M{"general.type": resourceDetails.Type, "general.version": resourceDetails.Version}

	if resourceDetails.Category != "" {
		filters["general.category"] = resourceDetails.Category
	}
	cursor, err := collection.Find(custom.DB.Ctx, filters, opts)
	if err != nil {
		return nil, errors.New("unknown db error")
	}

	var results []Record
	for cursor.Next(custom.DB.Ctx) {
		var doc struct {
			General struct {
				Name          string `bson:"name"`
				CanonicalName string `bson:"canonical_name"`
				GType         string `bson:"gtype"`
				Category      string `bson:"category"`
			} `bson:"general"`
		}
		cursor.Decode(&doc)
		results = append(
			results,
			Record{
				Name:          doc.General.Name,
				CanonicalName: doc.General.CanonicalName,
				GType:         doc.General.GType,
				Category:      doc.General.Category,
			},
		)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		fmt.Printf("Type: %s, Name: %s\n", result.CanonicalName, result.Name)
	}

	return results, nil
}
