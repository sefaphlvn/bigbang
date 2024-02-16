package poker

import (
	"fmt"

	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"go.mongodb.org/mongo-driver/bson"
)

func PokerLds(db *db.WTF, name string) {
	filter := bson.D{
		{Key: "$match", Value: bson.M{"resource.resource.eds_cluster_config.service_name": name}},
	}

	clusters, err := resources.GetGenerals(db, "clusters", filter)

	if err != nil {
		fmt.Println(err)
	}

	for _, cluster := range clusters {
		fmt.Println(cluster)
	}
}
