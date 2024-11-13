package config

type AppConfig struct {
	BigbangRestServerPort string `mapstructure:"BIGBANG_REST_SERVER_PORT" yaml:"BIGBANG_REST_SERVER_PORT"`
	BigbangAddress        string `mapstructure:"BIGBANG_ADDRESS" yaml:"BIGBANG_ADDRESS"`
	BigbangTLSEnabled     string `mapstructure:"BIGBANG_TLS_ENABLED" yaml:"BIGBANG_TLS_ENABLED"`
	BigbangGrpcPokePort   string `mapstructure:"BIGBANG_GRPC_POKE_PORT" yaml:"BIGBANG_GRPC_POKE_PORT"`

	MongodbHosts          string `mapstructure:"MONGODB_HOSTS" yaml:"MONGODB_HOSTS"`
	MongodbUsername       string `mapstructure:"MONGODB_USERNAME" yaml:"MONGODB_USERNAME"`
	MongodbPassword       string `mapstructure:"MONGODB_PASSWORD" yaml:"MONGODB_PASSWORD"`
	MongodbPort           string `mapstructure:"MONGODB_PORT" yaml:"MONGODB_PORT"`
	MongodbDatabase       string `mapstructure:"MONGODB_DATABASE" yaml:"MONGODB_DATABASE"`
	MongodbScheme         string `mapstructure:"MONGODB_SCHEME" yaml:"MONGODB_SCHEME"`
	MongodbReplicaSet     string `mapstructure:"MONGODB_REPLICASET" yaml:"MONGODB_REPLICASET"`
	MongodbTimeoutSeconds string `mapstructure:"MONGODB_TIMEOUTSECONDS" yaml:"MONGODB_TIMEOUTSECONDS"`

	LogLevel        string `mapstructure:"LOG_LEVEL" yaml:"LOG_LEVEL"`
	LogFormatter    string `mapstructure:"LOG_FORMATTER" yaml:"LOG_FORMATTER"`
	LogReportCaller string `mapstructure:"LOG_REPORTCALLER" yaml:"LOG_REPORTCALLER"`
}
