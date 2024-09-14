package resource

import (
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"google.golang.org/protobuf/types/known/anypb"
)

func (ar *AllResources) DecodeHTTPConnectionManager(resourceName string, context *db.AppContext) (*anypb.Any, []*models.ConfigDiscovery, error) {
	var message *anypb.Any
	var configDiscovery []*models.ConfigDiscovery

	resource, err := resources.GetResourceNGeneral(context, "extensions", resourceName, ar.Project)
	configDiscovery = resource.GetGeneral().ConfigDiscovery
	if err != nil {
		return nil, nil, err
	}

	hcmResource := resource.GetResource()
	hcmWithAccessLog, _ := ar.GetTypedConfigs(models.GeneralAccessLogTypedConfigPaths, hcmResource, context)

	httpConnectionManager := &hcm.HttpConnectionManager{}
	err = resources.MarshalUnmarshalWithType(hcmWithAccessLog, httpConnectionManager)
	if err != nil {
		return nil, nil, err
	}

	rds := httpConnectionManager.GetRds()
	if rds != nil {
		err = ar.GetRoutes(rds.RouteConfigName, context)
		if err != nil {
			return nil, nil, err
		}
	}

	message, err = anypb.New(httpConnectionManager)

	if err != nil {
		return nil, nil, err
	}

	return message, configDiscovery, nil
}
