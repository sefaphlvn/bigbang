package poker

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/filters"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

func PokerCds(context *db.AppContext, clusterName string, project string, processed *Processed) {
	cdsFilters := filters.ClusterDownstreamFilters(clusterName)
	for _, filter := range cdsFilters {
		resourceGeneral, err := resources.GetGenerals(context, filter.Collection, filter.Filter)
		if err != nil {
			context.Logger.Debug(err)
		}

		for _, general := range resourceGeneral {
			DetectChangedResource(general.GType, general.Name, project, context, processed)
		}
	}
}
