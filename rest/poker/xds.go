package poker

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/filters"
)

func PokerCds(context *db.AppContext, clusterName string, project string, processed *Processed) {
	cdsFilters := filters.ClusterDownstreamFilters(clusterName)
	for _, filter := range cdsFilters {
		CheckResource(context, filter.Filter, filter.Collection, project, processed)
	}
}

func PokerEds(context *db.AppContext, name string, project string, processed *Processed) {
	filter := filters.EdsDownstreamFilters(name)
	CheckResource(context, filter.Filter, filter.Collection, project, processed)
}

func PokerRoute(context *db.AppContext, name string, project string, processed *Processed) {
	filter := filters.RouteDownstreamFilters(name)
	CheckResource(context, filter.Filter, filter.Collection, project, processed)
}

func PokerHCEFS(context *db.AppContext, name string, project string, processed *Processed) {
	filter := filters.HCEFSDownstreamFilters(name)
	for _, filter := range filter {
		CheckResource(context, filter.Filter, filter.Collection, project, processed)
	}
}
