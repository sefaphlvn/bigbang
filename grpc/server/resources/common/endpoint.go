package common

import (
	"fmt"

	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

func (ar *Resources) GetEndpoints(ep string, db *db.WTF) {
	doc, _ := resources.GetResource(db, "endpoints", ep)
	singleEndpoint := &endpoint.ClusterLoadAssignment{}
	err := resources.GetResourceWithType(doc, singleEndpoint)

	if err != nil {
		fmt.Println(err)
	}

	ar.Endpoint = append(ar.Endpoint, singleEndpoint)
}
