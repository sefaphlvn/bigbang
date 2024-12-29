package config

type AppConfig struct {
	BigbangAddress    string `mapstructure:"BIGBANG_ADDRESS" yaml:"BIGBANG_ADDRESS"`
	BigbangPort       string `mapstructure:"BIGBANG_PORT" yaml:"BIGBANG_PORT"`
	BigbangTLSEnabled string `mapstructure:"BIGBANG_TLS_ENABLED" yaml:"BIGBANG_TLS_ENABLED"`

	MongodbHosts      string `mapstructure:"MONGODB_HOSTS" yaml:"MONGODB_HOSTS"`
	MongodbUsername   string `mapstructure:"MONGODB_USERNAME" yaml:"MONGODB_USERNAME"`
	MongodbPassword   string `mapstructure:"MONGODB_PASSWORD" yaml:"MONGODB_PASSWORD"`
	MongodbPort       string `mapstructure:"MONGODB_PORT" yaml:"MONGODB_PORT"`
	MongodbDatabase   string `mapstructure:"MONGODB_DATABASE" yaml:"MONGODB_DATABASE"`
	MongodbScheme     string `mapstructure:"MONGODB_SCHEME" yaml:"MONGODB_SCHEME"`
	MongodbReplicaSet string `mapstructure:"MONGODB_REPLICASET" yaml:"MONGODB_REPLICASET"`
	MongodbTimeoutMs  string `mapstructure:"MONGODB_TIMEOUTMS" yaml:"MONGODB_TIMEOUTMS"`
	MongodbTLSEnabled string `mapstructure:"MONGODB_TLS_ENABLED" yaml:"MONGODB_TLS_ENABLED"`

	LogLevel        string `mapstructure:"LOG_LEVEL" yaml:"LOG_LEVEL"`
	LogFormatter    string `mapstructure:"LOG_FORMATTER" yaml:"LOG_FORMATTER"`
	LogReportCaller string `mapstructure:"LOG_REPORTCALLER" yaml:"LOG_REPORTCALLER"`
}
