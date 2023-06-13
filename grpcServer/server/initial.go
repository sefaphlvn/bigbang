package server

import (
	"log"

	"github.com/sefaphlvn/bigbang/grpcServer/db"
	"github.com/sefaphlvn/bigbang/grpcServer/server/resources"
	"github.com/sefaphlvn/bigbang/restServer/models"
	"go.mongodb.org/mongo-driver/bson"
)

func getListenerList(db *db.MongoDB) []string {
	var serviceNames []string
	cur, err := db.GetGenerals("listeners")
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(db.Ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		var general models.General
		bsonBytes, _ := bson.Marshal(result["general"])
		bson.Unmarshal(bsonBytes, &general)

		serviceNames = append(serviceNames, general.Name)
	}
	return serviceNames
}

func InitialSnapshots(db *db.MongoDB, ctx *Context, l Logger) {
	serviceNames := getListenerList(db)
	var ss *resources.AllResources
	for _, serviceName := range serviceNames {
		rawListenerResource, err := resources.GetResource(db, "listeners", serviceName)
		if err != nil {
			log.Fatal(err)
		}

		lis, err := resources.SetSnapshot(rawListenerResource)
		if err != nil {
			log.Fatal(err)
		}
		ss = lis
	}

	ctx.SetSnashot(ss, l)
}
