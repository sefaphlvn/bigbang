package common

import (
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

func (ar *Resources) GetEndpoints(ep string, context *db.AppContext) {
	doc, _ := resources.GetResourceNGeneral(context, "endpoints", ep, ar.Project)
	singleEndpoint := &endpoint.ClusterLoadAssignment{}
	err := resources.MarshalUnmarshalWithType(doc.GetResource(), singleEndpoint)

	if err != nil {
		context.Logger.Debug(err)
	}

	ar.Endpoint = append(ar.Endpoint, singleEndpoint)
}
