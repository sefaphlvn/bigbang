package poker

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

var handlers map[models.GTypes]ResourceHandler

type PokerHCMHandler struct{}
type PokerCdsHandler struct{}
type PokerEdsHandler struct{}
type PokerRouterHandler struct{}
type PokerRouteHandler struct{}
type PokerTcpProxyHandler struct{}
type PokerAccessLogHandler struct{}
type PokerTLSHandler struct {
	gType models.GTypes
}

type ResourceHandler interface {
	Handle(context *db.AppContext, resourceName string, processed *Processed)
}

// CDS
func (h *PokerCdsHandler) Handle(context *db.AppContext, resourceName string, processed *Processed) {
	PokerCds(context, resourceName, processed)
}

// EDS
func (h *PokerEdsHandler) Handle(context *db.AppContext, resourceName string, processed *Processed) {
	PokerEds(context, resourceName, processed)
}

// ROUTER
func (h *PokerRouterHandler) Handle(context *db.AppContext, resourceName string, processed *Processed) {
	PokerRouter(context, resourceName, processed)
}

// ROUTE
func (h *PokerRouteHandler) Handle(context *db.AppContext, resourceName string, processed *Processed) {
	PokerRoute(context, resourceName, processed)
}

// HCM
func (h *PokerHCMHandler) Handle(context *db.AppContext, resourceName string, processed *Processed) {
	PokerHCM(context, resourceName, processed)
}

// TCP PROXY
func (h *PokerTcpProxyHandler) Handle(context *db.AppContext, resourceName string, processed *Processed) {
	PokerTcpProxy(context, resourceName, processed)
}

// TCP PROXY
func (h *PokerAccessLogHandler) Handle(context *db.AppContext, resourceName string, processed *Processed) {
	PokerAccessLog(context, resourceName, processed)
}

// TLS
func (h *PokerTLSHandler) Handle(context *db.AppContext, resourceName string, processed *Processed) {
	PokerTLS(context, resourceName, h.gType, processed)
}

func init() {
	handlers = map[models.GTypes]ResourceHandler{
		models.Endpoint:                     &PokerEdsHandler{},
		models.Cluster:                      &PokerCdsHandler{},
		models.Router:                       &PokerRouterHandler{},
		models.Route:                        &PokerRouteHandler{},
		models.HTTPConnectionManager:        &PokerHCMHandler{},
		models.TcpProxy:                     &PokerTcpProxyHandler{},
		models.FileAccessLog:                &PokerAccessLogHandler{},
		models.FluentdAccessLog:             &PokerAccessLogHandler{},
		models.StdErrAccessLog:              &PokerAccessLogHandler{},
		models.StdoutAccessLog:              &PokerAccessLogHandler{},
		models.DownstreamTlsContext:         &PokerTLSHandler{gType: models.DownstreamTlsContext},
		models.UpstreamTlsContext:           &PokerTLSHandler{gType: models.UpstreamTlsContext},
		models.TlsCertificate:               &PokerTLSHandler{gType: models.TlsCertificate},
		models.CertificateValidationContext: &PokerTLSHandler{gType: models.CertificateValidationContext},
	}
}
