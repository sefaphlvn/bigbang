package models

const (
	APITypePrefix         = "type.googleapis.com/"
	HTTPConnectionManager = APITypePrefix + "envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager"
	Router                = APITypePrefix + "envoy.extensions.filters.http.router.v3.Router"
)
