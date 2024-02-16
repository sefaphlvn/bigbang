package config

import "time"

type AppConfig struct {
	ServerPort  string `json:"server_port" yaml:"ServerPort"`
	GrpcService string `json:"grpc_service" yaml:"GrpcService"`
	MongoDB     struct {
		Hosts          []string      `json:"hosts" yaml:"hosts"`
		Username       string        `json:"username" yaml:"username"`
		Password       string        `json:"password" yaml:"password"`
		Port           string        `json:"port" yaml:"port"`
		Database       string        `json:"database" yaml:"database"`
		Scheme         string        `json:"scheme" yaml:"scheme"`
		ReplicaSet     string        `json:"replica_set" yaml:"replicaSet"`
		TimeoutSeconds time.Duration `json:"timeout_seconds" yaml:"timeoutSeconds"`
	} `json:"mongo_db" yaml:"mongoDB"`
	Log struct {
		Level        string `json:"level" yaml:"level"`
		Formatter    string `json:"formatter" yaml:"formatter"`
		ReportCaller bool   `json:"report_caller" yaml:"reportCaller"`
	} `json:"log" yaml:"log"`
}
