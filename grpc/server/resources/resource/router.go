package resource

import (
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"google.golang.org/protobuf/types/known/anypb"
)

func (ar *AllResources) DecodeRouter(resourceName string, db *db.AppContext) (*anypb.Any, []*models.ConfigDiscovery, error) {
	var message *anypb.Any
	resource, err := resources.GetResourceNGeneral(db, "extensions", resourceName, ar.Project)
	if err != nil {
		return nil, nil, err
	}

	singleRouter := &router.Router{}
	err = resources.MarshalUnmarshalWithType(resource.GetResource(), singleRouter)
	if err != nil {
		return nil, nil, err
	}

	message, _ = anypb.New(singleRouter)

	return message, nil, nil
}
