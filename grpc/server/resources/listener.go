package resources

import (
	"encoding/json"
	"errors"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/sefaphlvn/bigbang/grpc/models"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
)

func (R *AllResources) DecodeListener(resource *models.Resource, db *db.MongoDB, logger *logrus.Logger) {
	resArray, ok := resource.Resource.(primitive.A)
	R.Version = resource.Version

	if !ok {
		logger.Fatal("Unexpected resource format")
	}

	for _, res := range resArray {
		data, err := json.Marshal(res)
		if err != nil {
			logger.Fatal(err)
		}

		datas, _ := removeFilterChains(res)
		datass, err := json.Marshal(datas)

		filterChains := R.mergeFilters(data, db, logger)
		singleListener := &listener.Listener{}

		err = protojson.Unmarshal(datass, singleListener)

		singleListener.FilterChains = filterChains
		if err != nil {
			logger.Fatal(err, "sss")
		}

		R.Listener = append(R.Listener, singleListener)
	}
}

func removeFilterChains(data interface{}) (interface{}, error) {
	// Veriyi map[string]interface{} türüne dönüştür

	mapData, ok := data.(primitive.M)
	if !ok {
		return nil, errors.New("data is not a map[string]interface{}")
	}

	delete(mapData, "filter_chains")

	return mapData, nil
}

func (R *AllResources) mergeFilters(data []byte, db *db.MongoDB, logger *logrus.Logger) []*listener.FilterChain {
	listenera := helper.ToMapStringInterface(data)
	filterChains := helper.ItoGenericTypeConvert[[]interface{}](listenera["filter_chains"])
	fc := []*listener.FilterChain{}
	for _, filterChain := range filterChains { // filter_chains
		tFilterChain := helper.ItoGenericTypeConvert[map[string]interface{}](filterChain)
		filters := helper.ItoGenericTypeConvert[[]interface{}](tFilterChain["filters"])
		fc = []*listener.FilterChain{{Filters: R.detectCollectFilter(filters, db, logger)}}
	}

	//newData, _ := json.Marshal(listenera)
	//helper.PrettyPrinter(fc)
	return fc
}

func (R *AllResources) detectCollectFilter(filters []interface{}, db *db.MongoDB, logger *logrus.Logger) []*listener.Filter {
	var newFilters []*listener.Filter
	for i := range filters {
		filter := helper.ItoGenericTypeConvert[map[string]interface{}](filters[i])

		resource, err := GetResource(db, "extensions", filter["name"].(string))
		if err != nil {
			logger.Fatal(err)
		}

		data, err := json.Marshal(resource.Resource)
		if err != nil {
			logger.Fatal(err)
		}

		hcmman := &hcm.HttpConnectionManager{}
		_ = protojson.Unmarshal(data, hcmman)

		pbst, _ := anypb.New(hcmman)
		aaas := &listener.Filter{
			Name: filter["name"].(string),
			ConfigType: &listener.Filter_TypedConfig{
				TypedConfig: pbst,
			},
		}

		newFilters = append(newFilters, aaas)
	}
	//helper.PrettyPrinter(newFilters)
	return newFilters
}
