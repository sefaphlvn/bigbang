package poker

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/filters"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

func PokerAccessLog(context *db.AppContext, name string, project string, processed *Processed) {
	filters := filters.ALSDownstreamFilters(name)

	for _, filter := range filters {
		resourceGeneral, err := resources.GetGenerals(context, filter.Collection, filter.Filter)
		if err != nil {
			context.Logger.Debug(err)
		}

		for _, general := range resourceGeneral {
			DetectChangedResource(general.GType, general.Name, project, context, processed)
		}
	}
}
