package crud

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

func GetBootstrap(listenerGeneral models.General, config *config.AppConfig) map[string]interface{} {
	now := time.Now()
	CreatedAt := primitive.NewDateTimeFromTime(now)
	UpdatedAt := primitive.NewDateTimeFromTime(now)
	nodeID := fmt.Sprintf("%s:%s", listenerGeneral.Name, listenerGeneral.Project)

	cluster := createClusterConfig()
	admin := createAdminConfig()
	data := createDataConfig(nodeID, config.BigbangAddress, listenerGeneral.Version, cluster, admin)
	general := createGeneralConfig(listenerGeneral, CreatedAt, UpdatedAt)

	return map[string]interface{}{
		"general":  general,
		"resource": map[string]interface{}{"version": "1", "resource": data},
	}
}

func createClusterConfig() map[string]interface{} {
	cluster := map[string]interface{}{
		"name": "bigbang-controller",
	}
	return cluster
}

func createDataConfig(nodeID, authority, version string, cluster, admin map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"node": map[string]interface{}{
			"id": nodeID,
			"cluster": nodeID,
		},
		"static_resources": map[string]interface{}{
			"clusters": []interface{}{cluster},
		},
		"dynamic_resources": map[string]interface{}{
			"lds_config": map[string]interface{}{
				"ads": map[string]interface{}{},
				"resource_api_version": "V3",
			},
			"cds_config": map[string]interface{}{
				"ads": map[string]interface{}{},
				"resource_api_version": "V3",
			},
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
								"key":   "nodeid",
								"value": nodeID,
							},
							map[string]interface{}{
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

func createAdminConfig() map[string]interface{} {
	return map[string]interface{}{
		"address": map[string]interface{}{
			"socket_address": map[string]interface{}{
				"Protocol":   "TCP",
				"address":    "0.0.0.0",
				"port_value": 30090,
			},
		},
	}
}

func createGeneralConfig(listenerGeneral models.General, createdAt, updatedAt primitive.DateTime) map[string]interface{} {
	return map[string]interface{}{
		"name":                 listenerGeneral.Name,
		"version":              listenerGeneral.Version,
		"type":                 "bootstrap",
		"gtype":                "envoy.config.bootstrap.v3.Bootstrap",
		"canonical_name":       "config.bootstrap.v3.Bootstrap",
		"category":             "bootstrap",
		"collection":           "bootstrap",
		"project":              listenerGeneral.Project,
		"permissions":          map[string]interface{}{"users": []interface{}{}, "groups": []interface{}{}},
		"additional_resources": []interface{}{},
		"created_at":           createdAt,
		"updated_at":           updatedAt,
		"config_discovery":     []interface{}{},
		"typed_config":         []interface{}{},
	}
}

func createTLSTransportSocket() map[string]interface{} {
	return map[string]interface{}{
		"name": "envoy.transport_sockets.tls",
		"typed_config": map[string]interface{}{
			"@type": "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext",
		},
	}
}
