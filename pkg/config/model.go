package config

type AppConfig struct {
	BIGBANG_REST_SERVER_PORT string `mapstructure:"BIGBANG_REST_SERVER_PORT" yaml:"BIGBANG_REST_SERVER_PORT"`
	BIGBANG_ADDRESS          string `mapstructure:"BIGBANG_ADDRESS" yaml:"BIGBANG_ADDRESS"`
	BIGBANG_TLS_ENABLED      string `mapstructure:"BIGBANG_TLS_ENABLED" yaml:"BIGBANG_TLS_ENABLED"`
	BIGBANG_GRPC_POKE_PORT   string `mapstructure:"BIGBANG_GRPC_POKE_PORT" yaml:"BIGBANG_GRPC_POKE_PORT"`

	MONGODB_HOSTS          string `mapstructure:"MONGODB_HOSTS" yaml:"MONGODB_HOSTS"`
	MONGODB_USERNAME       string `mapstructure:"MONGODB_USERNAME" yaml:"MONGODB_USERNAME"`
	MONGODB_PASSWORD       string `mapstructure:"MONGODB_PASSWORD" yaml:"MONGODB_PASSWORD"`
	MONGODB_PORT           string `mapstructure:"MONGODB_PORT" yaml:"MONGODB_PORT"`
	MONGODB_DATABASE       string `mapstructure:"MONGODB_DATABASE" yaml:"MONGODB_DATABASE"`
	MONGODB_SCHEME         string `mapstructure:"MONGODB_SCHEME" yaml:"MONGODB_SCHEME"`
	MONGODB_REPLICASET     string `mapstructure:"MONGODB_REPLICASET" yaml:"MONGODB_REPLICASET"`
	MONGODB_TIMEOUTSECONDS string `mapstructure:"MONGODB_TIMEOUTSECONDS" yaml:"MONGODB_TIMEOUTSECONDS"`

	LOG_LEVEL        string `mapstructure:"LOG_LEVEL" yaml:"LOG_LEVEL"`
	LOG_FORMATTER    string `mapstructure:"LOG_FORMATTER" yaml:"LOG_FORMATTER"`
	LOG_REPORTCALLER string `mapstructure:"LOG_REPORTCALLER" yaml:"LOG_REPORTCALLER"`
}
