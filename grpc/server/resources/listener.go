package resources

import (
	"encoding/json"

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

		singleListener := R.mergeFilters(data, db, resource.General.AdditionalResources, logger)

		R.Listener = append(R.Listener, singleListener)
	}
}

func removeFilterChains(data interface{}) (interface{}, error) {
	mapData := helper.ItoGenericTypeConvert[map[string]interface{}](data)
	delete(mapData, "filter_chains")
	return mapData, nil
}

func (R *AllResources) mergeFilters(data []byte, db *db.MongoDB, additionalResource []models.AdditionalResource, logger *logrus.Logger) *listener.Listener {
	listenera := helper.ToMapStringInterface(data)
	filterChains := helper.ItoGenericTypeConvert[[]interface{}](listenera["filter_chains"])
	datas, _ := removeFilterChains(listenera)
	datass, _ := json.Marshal(datas)

	singleListener := &listener.Listener{}
	_ = protojson.Unmarshal(datass, singleListener)

	typedFilterChains := make([]*listener.FilterChain, 0, len(filterChains))
	for _, filterChain := range filterChains {
		tFilterChain := helper.ItoGenericTypeConvert[map[string]interface{}](filterChain)
		filters := []*listener.Filter{}
		for _, addResource := range additionalResource {
			if tFilterChain["name"] == addResource.ParentName {
				for _, extensions := range addResource.Extensions {
					anyExtension, additionalResource2, _ := R.CreateDynamicFilter(extensions.GType, extensions.Name, addResource.ParentName, db)
					filters = append(filters, &listener.Filter{
						Name: extensions.Name,
						ConfigType: &listener.Filter_TypedConfig{
							TypedConfig: anyExtension,
						}})
					R.CollectExtensions(additionalResource2, db)
				}
			}
		}

		delete(tFilterChain, "filters")
		jsonBytes, _ := json.Marshal(tFilterChain)
		var filterChain2 listener.FilterChain
		_ = protojson.Unmarshal(jsonBytes, &filterChain2)

		filterChain2.Filters = filters
		typedFilterChains = append(typedFilterChains, &filterChain2)
	}

	singleListener.FilterChains = typedFilterChains
	return singleListener
}

func (R *AllResources) detectCollectFilter(filters []interface{}, db *db.MongoDB, logger *logrus.Logger) []*listener.Filter {
	var newFilters []*listener.Filter
	for i := range filters {
		filter := helper.ItoGenericTypeConvert[map[string]interface{}](filters[i])

		resource, err := GetResource(db, "extensions", filter["name"].(string))
		if err != nil {
			logger.Fatal(err)
		}

		data, err := json.Marshal(resource)
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
