package resources

import (
	"encoding/json"
	"fmt"
	"log"

	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/sefaphlvn/bigbang/grpcServer/db"
	"github.com/sefaphlvn/bigbang/grpcServer/helper"
	"github.com/sefaphlvn/bigbang/grpcServer/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/encoding/protojson"
)

func (R *AllResources) DecodeListener(resource *models.Resource, db *db.MongoDB) {
	resArray, ok := resource.Resource.(primitive.A)
	R.Version = resource.Version

	if !ok {
		log.Fatal("Unexpected resource format")
	}

	for _, res := range resArray {
		data, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
		}

		data = R.mergeFilters(data, db)
		singleListener := &listener.Listener{}
		err = protojson.Unmarshal(data, singleListener)
		if err != nil {
			log.Fatal(err)
		}

		R.Listener = append(R.Listener, singleListener)
	}
}

func (R *AllResources) mergeFilters(data []byte, db *db.MongoDB) []byte {
	listener := helper.ToMapStringInterface(data)
	filterChains := helper.ItoGenericTypeConvert[[]interface{}](listener["filter_chains"])

	for _, filterChain := range filterChains { // filter_chains
		tFilterChain := helper.ItoGenericTypeConvert[map[string]interface{}](filterChain)
		filters := helper.ItoGenericTypeConvert[[]interface{}](tFilterChain["filters"])
		tFilterChain["filters"] = R.detectCollectFilter(filters, db)
	}
	newData, _ := json.Marshal(listener)

	return newData
}

func (R *AllResources) detectCollectFilter(filters []interface{}, db *db.MongoDB) interface{} {
	for i := range filters {
		filter := helper.ItoGenericTypeConvert[map[string]interface{}](filters[i])
		fmt.Println(filter)

		resource, err := GetResource(db, "extensions", filter["name"].(string))
		if err != nil {
			return nil
		}
		filters[i] = resource.Resource
	}

	return filters
}
