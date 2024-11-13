package common

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

func IsDeletable(ctx context.Context, appCtx *db.AppContext, gtype models.GTypes, name string) []string {
	downstreamFilters := gtype.DownstreamFilters(name)
	var deletableNames []string

	for _, filter := range downstreamFilters {
		fmt.Println("Filter: ", filter)
		collection := appCtx.Client.Collection(filter.Collection)
		cursor, err := collection.Find(ctx, filter.Filter, options.Find())
		if err != nil {
			log.Printf("Error finding documents: %v", err)
			continue
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var result struct {
				General struct {
					Name  string `bson:"name"`
					GType string `bson:"gtype"`
				} `bson:"general"`
			}

			if err := cursor.Decode(&result); err != nil {
				log.Printf("Error decoding document: %v", err)
				continue
			}
			targetGtype := models.GTypes(result.General.GType)
			combined := fmt.Sprintf("%s - %s", result.General.Name, targetGtype.PrettyName())
			deletableNames = append(deletableNames, combined)
		}
		if err := cursor.Err(); err != nil {
			log.Printf("Cursor error: %v", err)
		}
	}

	return deletableNames
}
