package poker

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"go.mongodb.org/mongo-driver/bson"
)

func PokerEds(wtf *db.WTF, name string) {
	filter := bson.D{{Key: "resource.resource.eds_cluster_config.service_name", Value: name}}

	rGeneral, err := resources.GetGenerals(wtf, "clusters", filter)
	if err != nil {
		wtf.Logger.Debug(err)
	}

	for _, general := range rGeneral {
		DetectChangedResource(general.GType, general.Name, wtf)
	}
}
