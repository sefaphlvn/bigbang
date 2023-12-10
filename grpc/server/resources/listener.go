package resources

import (
	"encoding/json"

	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/sefaphlvn/bigbang/grpc/models"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/encoding/protojson"
)

func (r *AllResources) DecodeListener(rawListenerResource *models.DBResource, db *db.MongoDB, logger *logrus.Logger) {
	resArray, ok := rawListenerResource.Resource.Resource.(primitive.A)
	if !ok {
		logger.Fatal("Unexpected resource format")
	}

	r.Version = rawListenerResource.Resource.Version

	for _, res := range resArray {
		data, err := json.Marshal(res)
		if err != nil {
			logger.Fatal(err)
		}

		singleListener := &listener.Listener{}
		err = protojson.Unmarshal(data, singleListener)
		if err != nil {
			logger.Fatal(err)
		}

		route, err := r.GetRoutes(db)
		if err != nil {
			logger.Fatal(err)
		}

		r.Route = route

		r.CollectExtensions(rawListenerResource.General.AdditionalResources, db)

		r.Listener = append(r.Listener, singleListener)
	}
}
