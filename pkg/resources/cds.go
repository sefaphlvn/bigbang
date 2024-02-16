package resources

import (
	"encoding/json"
	"fmt"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"google.golang.org/protobuf/encoding/protojson"
)

func GetResourceDependencyForCds(data *models.DBResource) *cluster.Cluster {
	jsonData, err := json.Marshal(data.Resource.Resource)
	if err != nil {
		fmt.Println(err)
	}

	cds := &cluster.Cluster{}
	err = protojson.Unmarshal(jsonData, cds)
	if err != nil {
		return nil
	}
	return cds
}

