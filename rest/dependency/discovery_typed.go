package dependency

import (
	"github.com/tidwall/gjson"

	"github.com/sefaphlvn/bigbang/pkg/models"
)

func parseConfigDiscovery(ctx *AppHandler, rootResult gjson.Result, activeResource Depend) []Depend {
	var dependencies []Depend

	rootResult.Get("general.config_discovery").ForEach(func(_, discoveryItem gjson.Result) bool {
		gtypeStr := discoveryItem.Get("gtype").String()
		if gtypeStr == "" {
			return true
		}

		gtype := models.GTypes(gtypeStr)
		cdName := discoveryItem.Get("name").String()
		cdID, _ := ctx.getResourceData(gtype.CollectionString(), cdName, activeResource.Project)
		dependencies = append(dependencies, Depend{Name: cdName, Gtype: gtype, Collection: gtype.CollectionString(), Project: activeResource.Project, ID: cdID})
		return true
	})

	return dependencies
}

func parseTypedConfig(ctx *AppHandler, rootResult gjson.Result, activeResource Depend) []Depend {
	var dependencies []Depend

	rootResult.Get("general.typed_config").ForEach(func(_, typedItem gjson.Result) bool {
		gtypeStr := typedItem.Get("gtype").String()
		if gtypeStr == "" {
			return true
		}

		gtype := models.GTypes(gtypeStr)
		tcName := typedItem.Get("name").String()
		tcID, _ := ctx.getResourceData(gtype.CollectionString(), tcName, activeResource.Project)
		dependencies = append(dependencies, Depend{Name: tcName, Gtype: gtype, Collection: gtype.CollectionString(), Project: activeResource.Project, ID: tcID})
		return true
	})

	return dependencies
}
