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

func (R *AllResources) DecodeListener(resource *models.DBResource, db *db.MongoDB, logger *logrus.Logger) {
	resArray, ok := resource.Resource.Resource.(primitive.A)
	R.Version = resource.Resource.Version

	if !ok {
		logger.Fatal("Unexpected resource format")
	}

	for _, res := range resArray {
		data, err := json.Marshal(res)
		if err != nil {
			logger.Fatal(err)
		}
		singleListener := &listener.Listener{}
		err = protojson.Unmarshal(data, singleListener)
		if err != nil {
			logger.Fatal(err, "sss")
		}

		R.CollectExtensions(resource.General.AdditionalResources, db)
		R.Listener = append(R.Listener, singleListener)
	}
}
