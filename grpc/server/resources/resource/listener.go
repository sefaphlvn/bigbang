package resource

import (
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/encoding/protojson"
)

var unmarshaler = protojson.UnmarshalOptions{
	AllowPartial: true,
	// DiscardUnknown: true,
}

func (ar *AllResources) DecodeListener(rawListenerResource *models.DBResource, context *db.AppContext, logger *logrus.Logger) {
	resArray, ok := rawListenerResource.Resource.Resource.(primitive.A)
	if !ok {
		logger.Fatal("Unexpected resource format")
	}

	ar.SetVersion(rawListenerResource.Resource.Version)
	ar.SetProject(rawListenerResource.General.Project)

	var lstnrs []types.Resource
	for _, lstnr := range resArray {
		listenerWithTransportSocket, _ := ar.GetTypedConfigs(models.ListenerTypedConfigPaths, lstnr, context)

		singleListener := &listener.Listener{}
		err := resources.MarshalUnmarshalWithType(listenerWithTransportSocket, singleListener)
		if err != nil {
			logger.Errorf("Listener Unmarshall err: %s", err)
		}

		lstnrs = append(lstnrs, singleListener)
	}
	// burasi for icerisindeydi disina aldim kontrol et
	ar.SetListener(lstnrs)

	ar.CollectExtensions(rawListenerResource.General.ConfigDiscovery, context, logger)
}
