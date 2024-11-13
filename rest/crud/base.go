package crud

import (
	"context"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/poker"
)

type Application struct {
	Context *db.AppContext
	Poke    *bridge.PokeServiceClient
}

func HandleResourceChange(ctx context.Context, resource models.DBResourceClass, requestDetails models.RequestDetails, context *db.AppContext, project string, poke *bridge.PokeServiceClient) *poker.Processed {
	if requestDetails.SaveOrPublish == "publish" {
		initialProcessed := poker.Processed{Listeners: []string{}, Depends: []string{}}
		changedResources := poker.DetectChangedResource(ctx, resource.GetGeneral().GType, requestDetails.Name, project, context, &initialProcessed, poke)
		return changedResources
	}
	return nil
}
