package resources

import (
	"encoding/json"
	"fmt"
	"log"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/sefaphlvn/bigbang/grpc/models"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
)

func (R *AllResources) CollectExtensions(resource []models.AdditionalResource, db *db.MongoDB) {
	var typedExtensionConfig = []*core.TypedExtensionConfig{}
	for _, additionalResource := range resource {
		for _, extension := range additionalResource.Extensions {
			anyResource, addadditionalResource, _ := R.CreateDynamicFilter(extension.GType, extension.Name, db)
			typedExtensionConfig = append(typedExtensionConfig, &core.TypedExtensionConfig{
				Name:        additionalResource.ParentName,
				TypedConfig: anyResource,
			})

			R.Extensions = append(R.Extensions, typedExtensionConfig...)

			if addadditionalResource != nil {
				R.CollectExtensions(addadditionalResource, db)
			}
		}
	}

}

func (R *AllResources) CreateDynamicFilter(typeUrl string, resourceName string, db *db.MongoDB) (*anypb.Any, []models.AdditionalResource, error) {
	var message *anypb.Any
	var additionalResource []models.AdditionalResource
	switch typeUrl {
	case HTTPConnectionManager:
		resource, err := GetResource(db, "extensions", resourceName)

		additionalResource = resource.GetGeneral().AdditionalResources
		if err != nil {
			log.Fatal(err)
		}

		data, err := json.Marshal(resource.Resource.Resource)
		if err != nil {
			log.Fatal(err)
		}

		hcmman := &hcm.HttpConnectionManager{}
		err = protojson.Unmarshal(data, hcmman)
		if err != nil {
			log.Fatal(err)
		}

		aa := hcmman.GetRds()

		fmt.Println(aa.RouteConfigName)
		message, _ = anypb.New(hcmman)

	case Router:
		resource, err := GetResource(db, "extensions", resourceName)
		if err != nil {
			log.Fatal(err)
		}

		data, err := json.Marshal(resource.Resource.Resource)
		if err != nil {
			log.Fatal(err)
		}

		router := &router.Router{}
		err = protojson.Unmarshal(data, router)
		if err != nil {
			log.Fatal(err)
		}
		message, _ = anypb.New(router)

	default:
		return nil, nil, fmt.Errorf("unknown type URL: %s", typeUrl)
	}

	return message, additionalResource, nil
}

const (
	APITypePrefix         = "type.googleapis.com/"
	HTTPConnectionManager = "envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager"
	Router                = "envoy.extensions.filters.http.router.v3.Router"
)
