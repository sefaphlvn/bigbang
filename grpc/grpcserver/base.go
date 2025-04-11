package grpcserver

type GrpcServer struct {
	Name     string   `bson:"name" json:"name"`
	NodeIDs  []string `bson:"node_ids" json:"node_ids"`
	Address  string   `bson:"address" json:"address"`
	LastSync int64    `bson:"lastSync" json:"lastSync"`
}

type NodeData struct {
	Name     string `bson:"name" json:"name"`
	NodeID   string `bson:"node_id" json:"node_id"`
	Address  string `bson:"address" json:"address"`
	LastSync int64  `bson:"lastSync" json:"lastSync"`
}
