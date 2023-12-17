package resources

import (
	"encoding/json"
	"fmt"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/sefaphlvn/bigbang/grpc/models"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
)

func (r *AllResources) CollectExtensions(resource []models.AdditionalResource, db *db.MongoDB, logger *logrus.Logger) {
	var typedExtensionConfig []*core.TypedExtensionConfig
	for _, additionalResource := range resource {
		for _, extension := range additionalResource.Extensions {
			anyResource, additionalResources, err := r.CreateDynamicFilter(extension.GType, extension.Name, db)
			if err != nil {
				logger.Fatal(err)
			}

			typedExtensionConfig = append(typedExtensionConfig, &core.TypedExtensionConfig{
				Name:        additionalResource.ParentName,
				TypedConfig: anyResource,
			})

			r.Extensions = append(r.Extensions, typedExtensionConfig...)

			if additionalResources != nil {
				r.CollectExtensions(additionalResources, db, logger)
			}
		}
	}

}

func (r *AllResources) CreateDynamicFilter(typeUrl string, resourceName string, db *db.MongoDB) (*anypb.Any, []models.AdditionalResource, error) {
	var message *anypb.Any
	var additionalResource []models.AdditionalResource
	switch typeUrl {
	case HTTPConnectionManager:
		resource, err := GetResource(db, "extensions", resourceName)

		additionalResource = resource.GetGeneral().AdditionalResources
		if err != nil {
			return nil, nil, err
		}

		data, err := json.Marshal(resource.Resource.Resource)
		if err != nil {
			return nil, nil, err
		}

		httpConnectionManager := &hcm.HttpConnectionManager{}
		err = protojson.Unmarshal(data, httpConnectionManager)
		if err != nil {
			return nil, nil, err
		}

		rds := httpConnectionManager.GetRds()
		if rds != nil {
			r.Route, err = r.GetRoutes(rds.RouteConfigName, db)
			if err != nil {
				return nil, nil, err
			}
		}

		message, _ = anypb.New(httpConnectionManager)

	case Router:
		resource, err := GetResource(db, "extensions", resourceName)
		if err != nil {
			return nil, nil, err
		}

		data, err := json.Marshal(resource.Resource.Resource)
		if err != nil {
			return nil, nil, err
		}

		routerJson := &router.Router{}
		err = protojson.Unmarshal(data, routerJson)
		if err != nil {
			return nil, nil, err
		}

		message, _ = anypb.New(routerJson)

	default:
		return nil, nil, fmt.Errorf("unknown type URL: %s", typeUrl)
	}

	return message, additionalResource, nil
}

const (
	HTTPConnectionManager = "envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager"
	Router                = "envoy.extensions.filters.http.router.v3.Router"
)
