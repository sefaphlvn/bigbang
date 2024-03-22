package resource

import (
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/sefaphlvn/bigbang/grpc/server/resources/common"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"google.golang.org/protobuf/types/known/anypb"
)

func (ar *AllResources) DecodeHTTPConnectionManager(arp *common.Resources, resourceName string, wtf *db.WTF) (*anypb.Any, []*models.ConfigDiscovery, error) {
	var message *anypb.Any
	var configDiscovery []*models.ConfigDiscovery

	resource, err := resources.GetResource(wtf, "extensions", resourceName)
	configDiscovery = resource.GetGeneral().ConfigDiscovery
	if err != nil {
		return nil, nil, err
	}

	httpConnectionManager := &hcm.HttpConnectionManager{}
	err = resources.GetResourceWithType(resource.GetResource(), httpConnectionManager)
	if err != nil {
		return nil, nil, err
	}

	rds := httpConnectionManager.GetRds()
	if rds != nil {
		err = ar.GetRoutes(rds.RouteConfigName, wtf)
		if err != nil {
			return nil, nil, err
		}
	}

	message, _ = anypb.New(httpConnectionManager)

	return message, configDiscovery, nil
}
