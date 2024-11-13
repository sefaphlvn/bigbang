package poker

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	bridgeClient "github.com/sefaphlvn/bigbang/rest/bridge"
)

type Processed struct {
	ProcessedResources []string
	Listeners          []string
	Depends            []string
}

func DetectChangedResource(gType models.GTypes, resourceName, project string, context *db.AppContext, processed *Processed, poke *bridge.PokeServiceClient) *Processed {
	pathWithGtype := gType.String() + "===" + resourceName
	if gType != models.Listener {
		processed.Depends = append(processed.Depends, pathWithGtype)
	}

	if helper.Contains(processed.ProcessedResources, pathWithGtype) {
		return processed
	}
	processed.ProcessedResources = append(processed.ProcessedResources, pathWithGtype)

	if gType == models.Listener {
		if !helper.Contains(processed.Listeners, resourceName) {
			_, err := bridgeClient.PokeNode(*poke, resourceName, project)
			if err != nil {
				context.Logger.Debugf("Poke failed: %s\n", err)
			}

			processed.Listeners = append(processed.Listeners, resourceName)
			result := strings.Join(processed.Depends, " \n ")
			context.Logger.Infof("new version added to snapshot for (%s) processed resource paths: \n %s", resourceName, result)
		}
	} else {
		ProcessResource(context, gType, resourceName, project, processed, poke)
	}

	return processed
}

func ProcessResource(context *db.AppContext, gType models.GTypes, resourceName, project string, processed *Processed, poke *bridge.PokeServiceClient) {
	filterResults := gType.DownstreamFilters(resourceName)

	for _, filterResult := range filterResults {
		CheckResource(context, filterResult.Filter, filterResult.Collection, project, processed, poke)
	}
}

func CheckResource(context *db.AppContext, filter primitive.D, collection, project string, processed *Processed, poke *bridge.PokeServiceClient) {
	rGeneral, err := resources.GetGenerals(context, collection, filter)
	if err != nil {
		context.Logger.Debug(err)
	}

	for _, general := range rGeneral {
		DetectChangedResource(general.GType, general.Name, project, context, processed, poke)
	}
}
