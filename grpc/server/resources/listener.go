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

func (R *AllResources) DecodeListener(rawListenerResource *models.DBResource, db *db.MongoDB, logger *logrus.Logger) {
	resArray, ok := rawListenerResource.Resource.Resource.(primitive.A)
	if !ok {
		logger.Fatal("Unexpected resource format")
	}

	R.Version = rawListenerResource.Resource.Version

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

		R.CollectExtensions(rawListenerResource.General.AdditionalResources, db)
		R.Listener = append(R.Listener, singleListener)
	}
}
