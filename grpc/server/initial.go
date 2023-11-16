package server

import (
	"github.com/sefaphlvn/bigbang/grpc/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func (h *Handler) InitialSnapshots() {
	serviceNames := h.getListenerList()
	for _, serviceName := range serviceNames {
		allResource, err := h.GetAllResourcesFromListener(serviceName)
		if err != nil {
			h.L.Errorf("BULK GetConfigurationFromListener(%v): %v", serviceName, err)
		}
		h.Ctx.SetSnapshot(allResource, h.L)
	}
}

func (h *Handler) getListenerList() []string {
	var serviceNames []string
	cur, err := h.DB.GetGenerals("listeners")
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(h.DB.Ctx) {
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
