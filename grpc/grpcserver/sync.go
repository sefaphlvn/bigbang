package grpcserver

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sefaphlvn/bigbang/grpc/server/snapshot"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func AddOrUpdateGrpcServer(dbClient *mongo.Database, address, nodeID string) error {
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
	return err
}

func isNodeIDUnique(dbClient *mongo.Database, nodeID string) (bool, error) {
	collection := dbClient.Collection("grpc_servers")
	filter := bson.M{
		"$or": bson.A{
			bson.M{"nodeIDs": nodeID},
			bson.M{"name": GetHostname()},
		},
	}

	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func RemoveNodeID(dbClient *mongo.Database, nodeID string) error {
	collection := dbClient.Collection("grpc_servers")

	filter := bson.M{"name": GetHostname()}
	update := bson.M{
		"$pull": bson.M{
			"nodeIDs": nodeID,
		},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
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

	for {
		select {
		case <-ticker.C:
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
			collection.UpdateOne(context.TODO(), filter, update, opts)
		}
	}
}
