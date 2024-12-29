package scenarios

const TcpProxy = `
{
	"general": {
		"name": "{{ .Data.name }}",
		"version": "{{ .Version }}",
		"type": "network_filter",
		"gtype": "envoy.extensions.filters.network.tcp_proxy.v3.TcpProxy",
		"project": "{{ .Project }}",
		"collection": "filters",
		"canonical_name": "envoy.filters.network.tcp_proxy",
		"category": "envoy.filters.network",
		"metadata": { "from_template": true },
		"permissions": {
			"users": [],
			"groups": []
		},
		"config_discovery": [],
		"typed_config": null
	},
	"resource": {
		"version": "1",
		"resource": {
			"stat_prefix": "{{ .Data.stat_prefix }}",
			"cluster": "{{ .Data.cluster }}"
		}
	}
}
`
