package server

import (
	"fmt"

	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"

	"github.com/sefaphlvn/bigbang/pkg/helper"
)

func Testit(req *discovery.DeltaDiscoveryRequest, resp *discovery.DeltaDiscoveryResponse) {
	if req != nil {
		fmt.Println("---------------Request Start---------------")
		helper.PrettyPrint(req)
		fmt.Println("---------------Request End---------------")
	}

	if resp != nil {
		fmt.Println("---------------Response Start---------------")
		helper.PrettyPrint(resp)
		fmt.Println("---------------Response End---------------")
	}
}
