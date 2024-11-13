package server

import (
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"

	"github.com/sefaphlvn/bigbang/pkg/helper"
)

func Testit(req *discovery.DeltaDiscoveryRequest, resp *discovery.DeltaDiscoveryResponse) {
	if req != nil {
		helper.PrettyPrint(req)
	}

	if resp != nil {
		helper.PrettyPrint(resp)
	}
}
