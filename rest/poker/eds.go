package poker

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"go.mongodb.org/mongo-driver/bson"
)

func PokerEds(context *db.AppContext, name string, processed *Processed) {
	filter := bson.D{{Key: "resource.resource.eds_cluster_config.service_name", Value: name}}

	rGeneral, err := resources.GetGenerals(context, "clusters", filter)
	if err != nil {
		context.Logger.Debug(err)
	}

	for _, general := range rGeneral {
		DetectChangedResource(general.GType, general.Name, context, processed)
	}
}
