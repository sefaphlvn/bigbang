package config

type AppConfig struct {
	ServerPort         string `mapstructure:"SERVERPORT" yaml:"ServerPort"`
	PokePort string `mapstructure:"POKEPORT" yaml:"PokePort"`

	MongoDB_Hosts          string `mapstructure:"MONGODB_HOSTS" yaml:"MongoDB_Hosts"`
	MongoDB_Username       string `mapstructure:"MONGODB_USERNAME" yaml:"MongoDB_Username"`
	MongoDB_Password       string `mapstructure:"MONGODB_PASSWORD" yaml:"MongoDB_Password"`
	MongoDB_Port           string `mapstructure:"MONGODB_PORT" yaml:"MongoDB_Port"`
	MongoDB_Database       string `mapstructure:"MONGODB_DATABASE" yaml:"MongoDB_Database"`
	MongoDB_Scheme         string `mapstructure:"MONGODB_SCHEME" yaml:"MongoDB_Scheme"`
	MongoDB_ReplicaSet     string `mapstructure:"MONGODB_REPLICASET" yaml:"MongoDB_ReplicaSet"`
	MongoDB_TimeoutSeconds string `mapstructure:"MONGODB_TIMEOUTSECONDS" yaml:"MongoDB_TimeoutSeconds"`

	Log_Level        string `mapstructure:"LOG_LEVEL" yaml:"Log_Level"`
	Log_Formatter    string `mapstructure:"LOG_FORMATTER" yaml:"Log_Formatter"`
	Log_ReportCaller string `mapstructure:"LOG_REPORTCALLER" yaml:"Log_ReportCaller"`
}
