package grpcserver

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sefaphlvn/bigbang/grpc/server/snapshot"
)

func ResetGrpcServerNodeIDs(dbClient *mongo.Database) {
	currentTime := time.Now().Unix()
	collection := dbClient.Collection("grpc_servers")

	filter := bson.M{"name": GetHostname()}
	update := bson.M{
		"$set": bson.M{
			"lastSync": currentTime,
			"nodeIDs":  []string{},
		},
	}

	opts := options.Update().SetUpsert(false)
	_, _ = collection.UpdateOne(context.TODO(), filter, update, opts)
}

func AddOrUpdateGrpcServer(dbClient *mongo.Database, address, nodeID string) {
	currentTime := time.Now().Unix()
	collection := dbClient.Collection("grpc_servers")

	filter := bson.M{"name": GetHostname()}
	update := bson.M{
		"$set": bson.M{
			"address":  address,
			"lastSync": currentTime,
		},
		"$addToSet": bson.M{
			"nodeIDs": nodeID,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		fmt.Printf("Error adding or updating GRPC server: %v\n", err)
	}
}

func RemoveNodeID(dbClient *mongo.Database, nodeID string) {
	collection := dbClient.Collection("grpc_servers")

	filter := bson.M{"name": GetHostname()}
	update := bson.M{
		"$pull": bson.M{
			"nodeIDs": nodeID,
		},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Printf("Error removing node ID: %v\n", err)
	}
}

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "none"
	}

	return hostname
}

func ScheduleSetNodeIDs(ctxCache *snapshot.Context, client *mongo.Database) {
	time.AfterFunc(15*time.Second, func() {
		SetNodeIDs(ctxCache, client)
	})
}

func SetNodeIDs(ctxCache *snapshot.Context, client *mongo.Database) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		nodeIDs := ctxCache.Cache.Cache.GetStatusKeys()
		fmt.Println("Snapshot Keys:", nodeIDs)

		currentTime := time.Now().Unix()
		collection := client.Collection("grpc_servers")

		filter := bson.M{"name": GetHostname()}
		update := bson.M{
			"$set": bson.M{
				"lastSync": currentTime,
				"nodeIDs":  nodeIDs,
			},
		}

		opts := options.Update().SetUpsert(false)
		_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			fmt.Printf("Error updating GRPC server node IDs: %v\n", err)
		}
	}
}
