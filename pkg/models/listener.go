package models

type ListenerResource struct {
	Version  string      `json:"version" bson:"version"`
	Resource interface{} `json:"resource" bson:"resource"`
}
