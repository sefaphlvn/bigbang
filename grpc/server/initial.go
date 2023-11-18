package server

import (
	"github.com/sefaphlvn/bigbang/grpc/models"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type Handler struct {
	Ctx    *Context
	DB     *db.MongoDB
	Logger *logrus.Logger
}

func (h *Handler) InitialSnapshots() {
	serviceNames := h.getListenerList()
	for _, serviceName := range serviceNames {
		allResource, err := h.GetAllResourcesFromListener(serviceName)
		if err != nil {
			h.Logger.Errorf("BULK GetConfigurationFromListener(%v): %v", serviceName, err)
		}
		err = h.Ctx.SetSnapshot(allResource, h.Logger)
		if err != nil {
			h.Logger.Errorf("%s", err)
		}
	}
}

func (h *Handler) getListenerList() []string {
	var serviceNames []string
	cur, err := h.DB.GetGenerals("listeners")
	if err != nil {
		h.Logger.Fatal(err)
	}

	for cur.Next(h.DB.Ctx) {
		var result bson.M
		err = cur.Decode(&result)
		if err != nil {
			h.Logger.Fatal(err)
		}

		var general models.General
		bsonBytes, _ := bson.Marshal(result["general"])

		err = bson.Unmarshal(bsonBytes, &general)
		if err != nil {
			h.Logger.Errorf("%s", err)
			return nil
		}

		serviceNames = append(serviceNames, general.Name)
	}
	return serviceNames
}
