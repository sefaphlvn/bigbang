package scenario

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/rest/crud/extension"
	"github.com/sefaphlvn/bigbang/rest/crud/xds"
)

type AppHandler struct {
	Context   *db.AppContext
	XDS       *xds.AppHandler
	Extension *extension.AppHandler
}

func NewScenarioHandler(context *db.AppContext) *AppHandler {
	xdsHandler := xds.NewXDSHandler(context)
	extensionHandler := extension.NewExtensionHandler(context)
	return &AppHandler{
		Context:   context,
		XDS:       xdsHandler,
		Extension: extensionHandler,
	}
}
