package poker

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"go.mongodb.org/mongo-driver/bson"
)

func PokerRouter(wtf *db.WTF, name string) {
	filter := bson.D{{Key: "general.additional_resources.extensions.name", Value: name}}

	rGeneral, err := resources.GetGenerals(wtf, "extensions", filter)
	if err != nil {
		wtf.Logger.Debug(err)
	}

	for _, general := range rGeneral {
		DetectChangedResource(general.GType, general.Name, wtf)
	}
}
