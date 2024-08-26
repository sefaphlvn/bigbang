package crud

import (
	"time"

	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetBootstrap(listenerGeneral models.General, config *config.AppConfig) map[string]interface{} {
	now := time.Now()
	CreatedAt := primitive.NewDateTimeFromTime(now)
	UpdatedAt := primitive.NewDateTimeFromTime(now)
	authority := config.BIGBANG_ADDRESS
	portValue := 80
	if config.BIGBANG_TLS_ENABLED == "true" {
		portValue = 443
	}

	cluster := map[string]interface{}{
		"name":            "bigbang-controller",
		"type":            "STRICT_DNS",
		"connect_timeout": "0.50s",
		"lb_policy":       "ROUND_ROBIN",
		"load_assignment": map[string]interface{}{
			"cluster_name": "bigbang-controller",
			"endpoints": []interface{}{
				map[string]interface{}{
					"lb_endpoints": []interface{}{
						map[string]interface{}{
							"endpoint": map[string]interface{}{
								"address": map[string]interface{}{
									"socket_address": map[string]interface{}{
										"address":    authority,
										"port_value": portValue,
									},
								},
							},
						},
					},
				},
			},
		},
		"http2_protocol_options": map[string]interface{}{},
	}

	if config.BIGBANG_TLS_ENABLED == "true" {
		cluster["transport_socket"] = map[string]interface{}{
			"name": "envoy.transport_sockets.tls",
			"typed_config": map[string]interface{}{
				"@type": "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext",
				"common_tls_context": map[string]interface{}{
					"validation_context": map[string]interface{}{
						"trusted_ca": map[string]interface{}{
							"filename": "/etc/ssl/certs/ca-certificates.crt",
						},
					},
				},
			},
		}
	}

	data := map[string]interface{}{
		"node": map[string]interface{}{
			"id":      listenerGeneral.Name,
			"cluster": "aadsa",
		},
		"static_resources": map[string]interface{}{
			"clusters": []interface{}{cluster},
		},
		"dynamic_resources": map[string]interface{}{
			"ads_config": map[string]interface{}{
				"api_type":              "DELTA_GRPC",
				"transport_api_version": "V3",
				"grpc_services": []interface{}{
					map[string]interface{}{
						"envoy_grpc": map[string]interface{}{
							"cluster_name": "bigbang-controller",
							"authority":    authority,
						},
						"initial_metadata": []interface{}{
							map[string]interface{}{
								"key":   "bigbang-controller",
								"value": "1",
							},
						},
					},
				},
				"set_node_on_first_message_only": false,
			},
		},
		"admin": map[string]interface{}{
			"access_log": []interface{}{
				map[string]interface{}{
					"name": "envoy.access_loggers.stdout",
					"typed_config": map[string]interface{}{
						"@type": "type.googleapis.com/envoy.extensions.access_loggers.stream.v3.StdoutAccessLog",
					},
				},
			},
			"address": map[string]interface{}{
				"socket_address": map[string]interface{}{
					"address":    "127.0.0.1",
					"port_value": 9090,
				},
			},
		},
	}

	general := map[string]interface{}{
		"name":                 listenerGeneral.Name,
		"version":              listenerGeneral.Version,
		"type":                 "bootstrap",
		"gtype":                "envoy.config.bootstrap.v3.Bootstrap",
		"canonical_name":       "config.bootstrap.v3.Bootstrap",
		"category":             "",
		"permissions":          map[string]interface{}{"users": []interface{}{}, "groups": []interface{}{}},
		"additional_resources": []interface{}{},
		"created_at":           CreatedAt,
		"updated_at":           UpdatedAt,
		"config_discovery":     []interface{}{},
		"typed_config":         []interface{}{},
	}

	result := map[string]interface{}{
		"general":  general,
		"resource": map[string]interface{}{"version": "1", "resource": data},
	}

	return result
}
