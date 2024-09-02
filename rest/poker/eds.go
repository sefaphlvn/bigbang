package poker

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/filters"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

func PokerEds(context *db.AppContext, name string, project string, processed *Processed) {
	filter := filters.EdsDownstreamFilters(name)

	rGeneral, err := resources.GetGenerals(context, filter.Collection, filter.Filter)
	if err != nil {
		context.Logger.Debug(err)
	}

	for _, general := range rGeneral {
		DetectChangedResource(general.GType, general.Name, project, context, processed)
	}
}
