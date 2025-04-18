package common

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"

	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/models/downstreamfilters"
)

func IsDeletable(ctx context.Context, appCtx *db.AppContext, gtype models.GTypes, dfm downstreamfilters.DownstreamFilter) []string {
	downstreamFilters := gtype.DownstreamFilters(dfm)
	var deletableNames []string

	for _, filter := range downstreamFilters {
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

func IsDefaultResource(ctx context.Context, appCtx *db.AppContext, name string, collection string, project string) (bool, error) {
	defaultResourceNames := map[string][]string{
		"users":      {"admin"},
		"groups":     {"default"},
		"projects":   {"default"},
		"extensions": {"bigbang-controller-hpo"},
		"tls":        {"bigbang-controller-tls"},
		"clusters":   {"bigbang-controller"},
	}

	if names, exists := defaultResourceNames[collection]; exists {
		if slices.Contains(names, name) {
			return true, nil
		}
	}

	return false, nil
}
