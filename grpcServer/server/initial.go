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

func (h *Handler) InitialSnapshots(db *db.MongoDB, ctx *Context, l Logger) {
	serviceNames := getListenerList(db)
	for _, serviceName := range serviceNames {
		allResource, err := h.GetConfigurationFromListener(serviceName)
		if err != nil {
			l.Errorf("BULK GetConfigurationFromListener(%v): %v", serviceName, err)
		}
		ctx.SetSnapshot(allResource, l)
	}
}

func (h *Handler) GetConfigurationFromListener(serviceName string) (*resources.AllResources, error) {
	rawListenerResource, err := resources.GetResource(h.DB, "listeners", serviceName)
	if err != nil {
		return nil, err
	}

	lis, err := resources.SetSnapshot(rawListenerResource)
	if err != nil {
		return nil, err
	}

	return lis, nil
}
