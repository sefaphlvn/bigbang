package db

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WTF struct {
	Client *mongo.Database
	Ctx    context.Context
	Logger *logrus.Logger
	Config *config.AppConfig
}

var (
	admin_user      = "admin"
	admin_email     = "admin@navigazer.com"
	admin_role      = "admin"
	admin_active    = true
	admin_baseGroup = ""
)

func NewMongoDB(config *config.AppConfig, logger *logrus.Logger) *WTF {
	// connectionString := fmt.Sprintf("%s://%s%s", config.MongoDB.Scheme, hosts, config.MongoDB.Port)
	connectionString := fmt.Sprintf("%s://%s:%s@%s%s", config.MongoDB_Scheme, config.MongoDB_Username, config.MongoDB_Password, config.MongoDB_Hosts, config.MongoDB_Port)

	tM := reflect.TypeOf(bson.M{})
	reg := bson.NewRegistryBuilder().RegisterTypeMapEntry(bsontype.EmbeddedDocument, tM).Build()
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString).SetRegistry(reg))
	if err != nil {
		logger.Fatalf("%s", err)
	}

	database := client.Database(config.MongoDB_Database)
	_, err = collectCreateIndex(database, ctx, logger)
	if err != nil {
		logger.Fatalf("%s", err)
	}

	wtf := &WTF{
		Client: database,
		Ctx:    ctx,
		Logger: logger,
		Config: config,
	}

	userID, err := createAdminUser(wtf)
	if err != nil {
		logger.Infof("Admin user not created: %s", err)
	}

	if err := createAdminGroup(wtf, userID); err != nil {
		logger.Infof("Admin group not created: %s", err)
	}

	return wtf
}

func (db *WTF) GetGenerals(collectionName string) (*mongo.Cursor, error) {
	collection := db.Client.Collection(collectionName)
	findOptions := options.Find()
	findOptions.SetProjection(bson.D{{Key: "general", Value: 1}})

	cur, err := collection.Find(db.Ctx, bson.D{{}}, findOptions)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("listener not found")
		} else {
			return nil, errors.New("unknown db error")
		}
	}

	return cur, err
}

func indexExists(ctx context.Context, collection *mongo.Collection, indexName string) (bool, error) {
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return false, fmt.Errorf("indexExists: %w", err)
	}
	defer cursor.Close(ctx)
	var indexes []bson.M
	if err = cursor.All(ctx, &indexes); err != nil {
		return false, err
	}
	for _, index := range indexes {
		if index["name"] == indexName {
			return true, nil
		}
	}
	return false, nil
}

func collectCreateIndex(database *mongo.Database, ctx context.Context, logger *logrus.Logger) (interface{}, error) {
	indices := map[string]mongo.IndexModel{
		"users":      {Keys: bson.M{"username": 1}, Options: options.Index().SetUnique(true).SetName("username_1")},
		"groups":     {Keys: bson.M{"groupname": 1}, Options: options.Index().SetUnique(true).SetName("groupname_1")},
		"service":    {Keys: bson.M{"name": 1}, Options: options.Index().SetUnique(true).SetName("name_1")},
		"clusters":   {Keys: bson.M{"general.name": 1}, Options: options.Index().SetUnique(true).SetName("general_name_1")},
		"listeners":  {Keys: bson.M{"general.name": 1}, Options: options.Index().SetUnique(true).SetName("general_name_1")},
		"endpoints":  {Keys: bson.M{"general.name": 1}, Options: options.Index().SetUnique(true).SetName("general_name_1")},
		"routes":     {Keys: bson.M{"general.name": 1}, Options: options.Index().SetUnique(true).SetName("general_name_1")},
		"extensions": {Keys: bson.M{"general.name": 1}, Options: options.Index().SetUnique(true).SetName("general_name_1")},
		"vhds":       {Keys: bson.M{"general.name": 1}, Options: options.Index().SetUnique(true).SetName("general_name_1")},
		"bootstrap":  {Keys: bson.M{"general.name": 1}, Options: options.Index().SetUnique(true).SetName("general_name_1")},
		"secrets":    {Keys: bson.M{"general.name": 1}, Options: options.Index().SetUnique(true).SetName("general_name_1")},
	}

	for collectionName, index := range indices {
		collection := database.Collection(collectionName)
		indexName := getIndexName(index)
		if err := createIndex(ctx, collection, index, indexName); err != nil {
			logger.Fatalf("Failed to create index for %s: %v\n", collectionName, err)
		}
	}

	return nil, nil
}

func createIndex(ctx context.Context, collection *mongo.Collection, index mongo.IndexModel, indexName string) error {
	if indexName == "" {
		return errors.New("invalid index name")
	}
	exists, err := indexExists(ctx, collection, indexName)
	if err != nil {
		return fmt.Errorf("could not check for index existence: %v", err)
	}

	if !exists {
		_, err = collection.Indexes().CreateOne(ctx, index)
		if err != nil {
			return fmt.Errorf("could not create index for %v on collection %v: %w", index.Keys, collection.Name(), err)
		}
	}
	return nil
}

func getIndexName(index mongo.IndexModel) string {
	keys, ok := index.Keys.(bson.M)
	if !ok {
		return ""
	}

	nameParts := make([]string, 0, len(keys))

	for key, val := range keys {
		if nestedKeys, ok := val.(bson.M); ok {
			for nestedKey := range nestedKeys {
				nameParts = append(nameParts, key+"."+nestedKey+"_1")
			}
		} else {
			nameParts = append(nameParts, key+"_1")
		}
	}

	return strings.Join(nameParts, "_")
}

func createAdminUser(db *WTF) (string, error) {
	collection := db.Client.Collection("users")
	var user models.User
	err := collection.FindOne(db.Ctx, bson.M{"username": "admin"}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		hashedPassword := helper.HashPassword("admin")
		user.Password = &hashedPassword
		now := time.Now()

		user.Created_at = primitive.NewDateTimeFromTime(now)
		user.Updated_at = primitive.NewDateTimeFromTime(now)
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		user.Email = &admin_email
		user.Username = &admin_user
		user.Role = &admin_role
		user.BaseGroup = &admin_baseGroup
		user.Active = &admin_active

		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.Username, user.User_id, []string{}, nil, false, *user.Role)
		user.Token = &token
		user.Refresh_token = &refreshToken

		_, insertErr := collection.InsertOne(db.Ctx, user)
		if insertErr != nil {
			return "", insertErr
		}
	}
	return user.User_id, nil
}

func createAdminGroup(db *WTF, userID string) error {
	collection := db.Client.Collection("groups")
	var group models.Group
	err := collection.FindOne(db.Ctx, bson.M{"groupname": userID}).Decode(&group)
	if err == mongo.ErrNoDocuments && userID != "" {
		_, err = collection.InsertOne(db.Ctx, bson.M{
			"groupname":  "admin",
			"members":    []string{userID},
			"created_at": primitive.NewDateTimeFromTime(time.Now()),
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
