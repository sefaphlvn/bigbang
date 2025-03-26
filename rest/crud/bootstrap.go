package crud

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

func GetBootstrap(listenerGeneral models.General, config *config.AppConfig) map[string]any {
	now := time.Now()
	CreatedAt := primitive.NewDateTimeFromTime(now)
	UpdatedAt := primitive.NewDateTimeFromTime(now)
	nodeID := fmt.Sprintf("%s:%s", listenerGeneral.Name, listenerGeneral.Project)

	cluster := createClusterConfig()
	admin := createAdminConfig()
	data := createDataConfig(nodeID, config.BigbangAddress, listenerGeneral.Version, cluster, admin)
	general := createGeneralConfig(listenerGeneral, CreatedAt, UpdatedAt)

	return map[string]any{
		"general":  general,
		"resource": map[string]any{"version": "1", "resource": data},
	}
}

func createClusterConfig() map[string]any {
	cluster := map[string]any{
		"name": "bigbang-controller",
	}
	return cluster
}

func createDataConfig(nodeID, authority, version string, cluster, admin map[string]any) map[string]any {
	return map[string]any{
		"node": map[string]any{
			"id": nodeID,
			"cluster": nodeID,
		},
		"static_resources": map[string]any{
			"clusters": []any{cluster},
		},
		"dynamic_resources": map[string]any{
			"lds_config": map[string]any{
				"ads": map[string]any{},
				"resource_api_version": "V3",
			},
			"cds_config": map[string]any{
				"ads": map[string]any{},
				"resource_api_version": "V3",
			},
			"ads_config": map[string]any{
				"api_type":              "DELTA_GRPC",
				"transport_api_version": "V3",
				"grpc_services": []any{
					map[string]any{
						"envoy_grpc": map[string]any{
							"cluster_name": "bigbang-controller",
							"authority":    authority,
						},
						"initial_metadata": []any{
							map[string]any{
								"key":   "nodeid",
								"value": nodeID,
							},
							map[string]any{
								"key":   "envoy-version",
								"value": version,
							},
						},
					},
				},
				"set_node_on_first_message_only": false,
			},
		},
		"admin": admin,
	}
}

func createAdminConfig() map[string]any {
	return map[string]any{
		"address": map[string]any{
			"socket_address": map[string]any{
				"Protocol":   "TCP",
				"address":    "0.0.0.0",
				"port_value": 30090,
			},
		},
	}
}

func createGeneralConfig(listenerGeneral models.General, createdAt, updatedAt primitive.DateTime) map[string]any {
	return map[string]any{
		"name":                 listenerGeneral.Name,
		"version":              listenerGeneral.Version,
		"type":                 "bootstrap",
		"gtype":                "envoy.config.bootstrap.v3.Bootstrap",
		"canonical_name":       "config.bootstrap.v3.Bootstrap",
		"category":             "bootstrap",
		"collection":           "bootstrap",
		"project":              listenerGeneral.Project,
		"permissions":          map[string]any{"users": []any{}, "groups": []any{}},
		"additional_resources": []any{},
		"created_at":           createdAt,
		"updated_at":           updatedAt,
		"config_discovery":     []any{},
		"typed_config":         []any{},
	}
}

/* func createTLSTransportSocket() map[string]any {
	return map[string]any{
		"name": "envoy.transport_sockets.tls",
		"typed_config": map[string]any{
			"@type": "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext",
		},
	}
} */
