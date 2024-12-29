package scenarios

const BasicHcm = `
{
	"general": {
		"name": "{{ .Data.name }}",
		"version": "{{ .Version }}",
		"type": "network_filter",
		"gtype": "envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager",
		"project": "{{ .Project }}",
		"collection": "filters",
		"canonical_name": "envoy.filters.network.http_connection_manager",
		"category": "envoy.filters.network",
		"metadata": { "from_template": true },
		"permissions": {
			"users": [],
			"groups": []
		}
	},
	"resource": {
		"version": "1",
		"resource": {
			"stat_prefix": "{{ .Data.stat_prefix }}",
			"codec_type": "{{ .Data.codec_type }}",
			"route_config": {
				"name": "route1",
				"virtual_hosts": [
					{
						"name": "virtualhost1",
						"domains": {{ .Data.domains | toJson }},
						"routes": [
							{
								"name": "route1_1",
								"match": {
									"{{ .Data.match_key }}": "{{ .Data.match_value }}"
								},
								"route": {
									"cluster": "{{ .Data.cluster }}"
								}
							}
						]
					}
				]
			}
		}
	}
}
`

const RDSHcm = `
{
	"general": {
		"name": "{{ .Data.name }}",
		"version": "{{ .Version }}",
		"type": "network_filter",
		"gtype": "envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager",
		"project": "{{ .Project }}",
		"collection": "filters",
		"canonical_name": "envoy.filters.network.http_connection_manager",
		"category": "envoy.filters.network",
		"metadata": { "from_template": true },
		"permissions": {
			"users": [],
			"groups": []
		}
	},
	"resource": {
		"version": "1",
		"resource": {
			"codec_type": "{{ .Data.codec_type }}",
			"stat_prefix": "{{ .Data.stat_prefix }}",
			"rds": {
				"config_source": {
					"ads": {},
					"initial_fetch_timeout": "2.0s",
					"resource_api_version": "V3"
				},
				"route_config_name": "{{ .Data.rds }}"
			}
		}
	}
}
`
