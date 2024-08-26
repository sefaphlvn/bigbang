package poker

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateALFilters(name string) []MongoFilters {
	return []MongoFilters{
		{
			Collection: "extensions",
			Filter:     bson.D{{Key: "general.typed_config.name", Value: name}},
		},
		{
			Collection: "listeners",
			Filter:     bson.D{{Key: "general.typed_config.name", Value: name}},
		},
	}
}

func PokerAccessLog(context *db.AppContext, name string, processed *Processed) {
	filters := CreateALFilters(name)

	for _, filter := range filters {
		resourceGeneral, err := resources.GetGenerals(context, filter.Collection, filter.Filter)
		if err != nil {
			context.Logger.Debug(err)
		}

		for _, general := range resourceGeneral {
			DetectChangedResource(general.GType, general.Name, context, processed)
		}
	}

}
