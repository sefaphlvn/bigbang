package crud

func GetBootstrap(listenerName string) map[string]interface{} {
	authority := "bigbang.elchi.io"
	port_value := 80
	availableController := "1"

	data := map[string]interface{}{
		"node": map[string]interface{}{
			"id":      listenerName,
			"cluster": "aaasdasdsadsadasdsa",
		},
		"static_resources": map[string]interface{}{
			"clusters": []interface{}{
				map[string]interface{}{
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
													"port_value": port_value,
												},
											},
										},
									},
								},
							},
						},
					},
					"http2_protocol_options": map[string]interface{}{},
				},
			},
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
								"value": availableController,
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

	return data
}
