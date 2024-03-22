package resource

import (
	"encoding/json"
	"fmt"

	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/encoding/protojson"
)

var unmarshaler = protojson.UnmarshalOptions{
	AllowPartial: true, // Eğer tüm alanları doldurmadıysanız
	// DiscardUnknown: true, // Bilinmeyen alanları yok say
}

func (ar *AllResources) DecodeListener(rawListenerResource *models.DBResource, wtf *db.WTF, logger *logrus.Logger) {
	resArray, ok := rawListenerResource.Resource.Resource.(primitive.A)
	if !ok {
		logger.Fatal("Unexpected resource format")
	}

	ar.SetVersion(rawListenerResource.Resource.Version)

	var lstnr []types.Resource
	for _, singleListener := range resArray {
		listenerWithTransportSocket := ar.GetTransportSockets(models.TransportSocketPath, singleListener, wtf, logger)
		data, err := json.Marshal(listenerWithTransportSocket)
		if err != nil {
			logger.Error(err)
		}

		singleListener := &listener.Listener{}
		err = unmarshaler.Unmarshal(data, singleListener)
		if err != nil {
			logger.Errorf("Listener Unmarshall err: %s", err)
		}

		lstnr = append(lstnr, singleListener)
		ar.SetListener(lstnr)
	}

	ar.CollectExtensions(rawListenerResource.General.ConfigDiscovery, wtf, logger)
}

func (ar *AllResources) GetTransportSockets(pathd models.TypedPaths, jsonData interface{}, wtf *db.WTF, logger *logrus.Logger) interface{} {
	jsonString, err := json.Marshal(jsonData)
	if err != nil {
		logger.Debugf("Error marshalling JSON: %v", err)
		return jsonData
	}

	jsonStringStr := string(jsonString)
	for i := range gjson.Get(jsonStringStr, "filter_chains").Array() {
		path := fmt.Sprintf(pathd.String(), i)
		tempTypedConfig := resources.GetTypedConfigValue(jsonStringStr, path+".value", logger)
		if tempTypedConfig == nil {
			continue
		}

		conf, err := resources.GetResource(wtf, tempTypedConfig.Type, tempTypedConfig.Name)
		if err != nil {
			logger.Debugf("Error getting resource from DB: %v", err)
			continue
		}

		resource := conf.GetResource()
		ar.DecodeDownstreamTLS(conf, wtf)
		typed_config, ok := resource.(primitive.M)
		if !ok {
			logger.Debugf("Resource is not a map[string]interface{}")
			continue
		}

		typed_config["@type"] = "type.googleapis.com/" + tempTypedConfig.Gtype
		if jsonStringStr, err = sjson.Set(jsonStringStr, path, typed_config); err != nil {
			logger.Debugf("Error setting new config value with sjson.Set: %v", err)
		}
	}

	var updatedJSONData interface{}
	if err := json.Unmarshal([]byte(jsonStringStr), &updatedJSONData); err != nil {
		logger.Debugf("Error unmarshalling updated JSON: %v", err)
	}

	return updatedJSONData
}
