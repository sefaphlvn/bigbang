package custom

import (
	"errors"
	"fmt"
	"log"

	"github.com/sefaphlvn/bigbang/restServer/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Record struct {
	Type string `json:"type" bson:"type"`
	Name string `json:"name" bson:"name"`
}

func (custom *DBHandler) GetFilterChainFilters(resource models.DBResourceClass, resourceDetails models.ResourceDetails) (interface{}, error) {

	collection := custom.DB.Client.Collection("extensions")
	opts := options.Find()
	opts.SetProjection(bson.M{
		"general.name":    1,
		"general.subtype": 1,
	})

	cursor, err := collection.Find(custom.DB.Ctx, bson.M{"general.type": "filters"}, opts)
	if err != nil {
		return nil, errors.New("unknown db error")
	}

	var results []Record
	for cursor.Next(custom.DB.Ctx) {
		var doc struct {
			General struct {
				Name    string `bson:"name"`
				Subtype string `bson:"subtype"`
			} `bson:"general"`
		}
		cursor.Decode(&doc)
		results = append(results, Record{Type: doc.General.Subtype, Name: doc.General.Name})
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		fmt.Printf("Type: %s, Name: %s\n", result.Type, result.Name)
	}

	return results, nil
}
