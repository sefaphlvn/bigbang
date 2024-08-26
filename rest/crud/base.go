package crud

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/poker"
)

type Application struct {
	Context *db.AppContext
}

func HandleResourceChange(resource models.DBResourceClass, requestDetails models.RequestDetails, context *db.AppContext) *poker.Processed {
	if requestDetails.SaveOrPublish == "publish" {
		initialProcessed := poker.Processed{Listeners: []string{}, Depends: []string{}}
		changedResources := poker.DetectChangedResource(resource.GetGeneral().GType, requestDetails.Name, context, &initialProcessed)
		return changedResources
	}
	return nil
}
