package poker

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/filters"
)

func PokerAccessLog(context *db.AppContext, name string, project string, processed *Processed) {
	filters := filters.ALSDownstreamFilters(name)
	for _, filter := range filters {
		CheckResource(context, filter.Filter, filter.Collection, project, processed)
	}
}

func PokerHCM(context *db.AppContext, name string, project string, processed *Processed) {
	filter := filters.HcmDownstreamFilters(name)
	CheckResource(context, filter.Filter, filter.Collection, project, processed)
}

func PokerTcpProxy(context *db.AppContext, name string, project string, processed *Processed) {
	filter := filters.TcpProxyDownstreamFilters(name)
	CheckResource(context, filter.Filter, filter.Collection, project, processed)
}

func PokerRouter(context *db.AppContext, name string, project string, processed *Processed) {
	filter := filters.RouterDownstreamFilters(name)
	CheckResource(context, filter.Filter, filter.Collection, project, processed)
}
