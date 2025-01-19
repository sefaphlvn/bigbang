package poker

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/models/downstreamfilters"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	bridgeClient "github.com/sefaphlvn/bigbang/rest/bridge"
)

type Processed struct {
	ProcessedResources []string
	Listeners          []string
	Depends            []string
}

func DetectChangedResource(ctx context.Context, gType models.GTypes, version, resourceName, project string, context *db.AppContext, processed *Processed, poke *bridge.PokeServiceClient) *Processed {
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
			_, err := bridgeClient.PokeNode(ctx, *poke, resourceName, project, version)
			if err != nil {
				context.Logger.Debugf("Poke failed: %s\n", err)
			}

			processed.Listeners = append(processed.Listeners, resourceName)
			result := strings.Join(processed.Depends, " \n ")
			context.Logger.Infof("new version added to snapshot for (%s) processed resource paths: \n %s", resourceName, result)
		}
	} else {
		ProcessResource(ctx, context, gType, version, resourceName, project, processed, poke)
	}

	return processed
}

func ProcessResource(ctx context.Context, context *db.AppContext, gType models.GTypes, version, resourceName, project string, processed *Processed, poke *bridge.PokeServiceClient) {
	dfm := downstreamfilters.DownstreamFilter{
		Name:    resourceName,
		Project: project,
		Version: version,
	}
	filterResults := gType.DownstreamFilters(dfm)

	for _, filterResult := range filterResults {
		CheckResource(ctx, context, filterResult.Filter, filterResult.Collection, project, version, processed, poke)
	}
}

func CheckResource(ctx context.Context, context *db.AppContext, filter primitive.D, collection, project, version string, processed *Processed, poke *bridge.PokeServiceClient) {
	rGeneral, err := resources.GetGenerals(ctx, context, collection, filter)
	if err != nil {
		context.Logger.Debug(err)
	}

	for _, general := range rGeneral {
		DetectChangedResource(ctx, general.GType, version, general.Name, project, context, processed, poke)
	}
}
