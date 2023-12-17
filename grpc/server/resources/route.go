package resources

import (
	"encoding/json"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"google.golang.org/protobuf/encoding/protojson"
)

func (r *AllResources) GetRoutes(rdsName string, db *db.MongoDB) (*routev3.RouteConfiguration, error) {
	route, err := GetResource(db, "routes", rdsName)
	if err != nil {
		return nil, err
	}

	jsonRoute, err := json.Marshal(route.Resource.Resource)
	if err != nil {
		return nil, err
	}

	singleRouter := &routev3.RouteConfiguration{}
	err = protojson.Unmarshal(jsonRoute, singleRouter)
	if err != nil {
		return nil, err
	}

	return singleRouter, nil
}
