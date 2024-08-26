package xds

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/sefaphlvn/bigbang/pkg/models"
	snapshotStats "github.com/sefaphlvn/bigbang/pkg/stats"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (xds *AppHandler) GetResource(resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error) {
	collection := xds.Context.Client.Collection(requestDetails.Collection)
	filter := bson.M{"general.name": requestDetails.Name}
	filterWithRestriction := common.AddUserFilter(requestDetails, filter)
	fmt.Println(filterWithRestriction)
	result := collection.FindOne(xds.Context.Ctx, filterWithRestriction)

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("not found: (" + requestDetails.Name + ")")
		} else {
			return nil, errors.New("unknown db error")
		}
	}

	GetSnapshotsFromServer("localhost:18000")

	err := result.Decode(resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func GetSnapshotsFromServer(serverAddress string) {
	// gRPC client bağlantısını oluştur
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := snapshotStats.NewSnapshotKeyServiceClient(conn)

	// Metadata oluştur
	md := metadata.Pairs(
		"key1", "value1",
		"key2", "value2",
	)

	// Metadata ile context oluştur
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Snapshot verilerini al
	resp, err := client.GetSnapshotKeys(ctx, &snapshotStats.Empty{})
	if err != nil {
		log.Fatalf("could not get snapshots: %v", err)
	}

	fmt.Printf("Snapshot keys: %s\n", resp.Keys)

}
