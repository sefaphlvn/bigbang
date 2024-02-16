package resource

import (
	"encoding/json"

	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/encoding/protojson"
)

func (ar *AllResources) DecodeListener(rawListenerResource *models.DBResource, db *db.WTF, logger *logrus.Logger) {
	resArray, ok := rawListenerResource.Resource.Resource.(primitive.A)
	if !ok {
		logger.Fatal("Unexpected resource format")
	}

	ar.SetVersion(rawListenerResource.Resource.Version)

	var lstnr []types.Resource
	for _, res := range resArray {
		data, err := json.Marshal(res)
		if err != nil {
			logger.Error(err)
		}

		singleListener := &listener.Listener{}
		err = protojson.Unmarshal(data, singleListener)
		if err != nil {
			logger.Error(err)
		}

		lstnr = append(lstnr, singleListener)
		ar.SetListener(lstnr)
	}

	ar.CollectExtensions(rawListenerResource.General.AdditionalResources, db, logger)
}
