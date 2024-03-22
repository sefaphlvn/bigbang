package common

import (
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

func (ar *Resources) GetEndpoints(ep string, wtf *db.WTF) {
	doc, _ := resources.GetResource(wtf, "endpoints", ep)
	singleEndpoint := &endpoint.ClusterLoadAssignment{}
	err := resources.GetResourceWithType(doc.GetResource(), singleEndpoint)

	if err != nil {
		wtf.Logger.Debug(err)
	}

	ar.Endpoint = append(ar.Endpoint, singleEndpoint)
}
