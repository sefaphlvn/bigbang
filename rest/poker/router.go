package poker

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"go.mongodb.org/mongo-driver/bson"
)

func PokerRouter(context *db.AppContext, name string) {
	filter := bson.D{{Key: "general.config_discovery.extensions.name", Value: name}}

	rGeneral, err := resources.GetGenerals(context, "extensions", filter)
	if err != nil {
		context.Logger.Debug(err)
	}

	for _, general := range rGeneral {
		DetectChangedResource(general.GType, general.Name, context)
	}
}
