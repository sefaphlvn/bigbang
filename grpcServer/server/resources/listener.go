package resources

import (
	"encoding/json"
	"fmt"
	"log"

	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/sefaphlvn/bigbang/grpcServer/helper"
	"github.com/sefaphlvn/bigbang/grpcServer/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/encoding/protojson"
)

func (l *AllResources) DecodeListener(resource *models.Resource) {
	resArray, ok := resource.Resource.(primitive.A)
	if !ok {
		log.Fatal("Unexpected resource format")
	}

	for _, res := range resArray {
		fmt.Println(res)
		data, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
		}

		data = mergeFilters(data)
		singleListener := &listener.Listener{}
		err = protojson.Unmarshal(data, singleListener)
		if err != nil {
			fmt.Println("sadasdsadasdas")
			log.Fatal(err)

		}

		l.Listener = append(l.Listener, singleListener)
	}

	// prints out the string representation of all listeners
	/* for _, lis := range l.listener {
		fmt.Println(lis.String())
	} */

}

func mergeFilters(data []byte) []byte {
	listener := helper.ToInterface(data)
	filterChains := listener["filter_chains"].([]interface{})

	for _, filters := range filterChains { // filter_chains
		tFilters := filters.(map[string]interface{})
		ff := tFilters["filters"].([]interface{})

		for i := range ff {
			ff[i] = helper.ToInterface([]byte(`{
				"name": "envoy.filters.network.http_connection_manager",
				"typed_config": {
					"@type": "type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager",
					"rds": {
						"route_config_name": "abroute",
						"config_source": {
							"ads": { },
							"initial_fetch_timeout": "15s",
							"resource_api_version": "V3"
						}
					},
					"stat_prefix": "ddd"
				}
			}`))
		}
	}
	newData, _ := json.Marshal(listener)
	fmt.Println(string(newData))

	return newData
}
