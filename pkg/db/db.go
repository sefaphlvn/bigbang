package db

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/errstr"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

type AppContext struct {
	Client *mongo.Database
	Ctx    context.Context
	Logger *logrus.Logger
	Config *config.AppConfig
}

var (
	adminUser                      = "admin"
	adminEmail                     = "admin@elchi.io"
	adminRole          models.Role = "owner"
	adminActive                    = true
	adminBaseGroup                 = ""
	generalProject                 = "general.project"
	generalName                    = "general.name"
	generalNameProject             = "general_name_project_1"
)

func NewMongoDB(config *config.AppConfig, logger *logrus.Logger) *AppContext {
	// connectionString := fmt.Sprintf("%s://%s%s", config.MongoDB.Scheme, hosts, config.MongoDB.Port)
	connectionString := fmt.Sprintf("%s://%s:%s@%s%s", config.MongodbScheme, config.MongodbUsername, config.MongodbPassword, config.MongodbHosts, config.MongodbPort)
	tM := reflect.TypeOf(bson.M{})
	reg := bson.NewRegistry()
	reg.RegisterTypeMapEntry(bson.TypeEmbeddedDocument, tM)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString).SetRegistry(reg))
	if err != nil {
		logger.Fatalf("%s", err)
	}

	database := client.Database(config.MongodbDatabase)
	err = collectCreateIndex(ctx, database, logger)
	if err != nil {
		logger.Fatalf("%s", err)
	}

	context := &AppContext{
		Client: database,
		Ctx:    ctx,
		Logger: logger,
		Config: config,
	}

	userID, err := createAdminUser(context)
	if err != nil {
		logger.Infof("Admin user not created: %s", err)
	}

	if err := createAdminGroup(context, userID); err != nil {
		logger.Infof("Admin group not created: %s", err)
	}

	if err := createDefaultProject(context, userID); err != nil {
		logger.Infof("Default project not created: %s", err)
	}

	return context
}

func (db *AppContext) GetGenerals(collectionName string) (*mongo.Cursor, error) {
	collection := db.Client.Collection(collectionName)
	findOptions := options.Find()
	findOptions.SetProjection(bson.D{{Key: "general", Value: 1}})

	cur, err := collection.Find(db.Ctx, bson.D{{}}, findOptions)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errstr.ErrListenerNotFound
		}
		return nil, errstr.ErrUnknownDBError
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

func collectCreateIndex(ctx context.Context, database *mongo.Database, logger *logrus.Logger) error {
	indices := map[string]mongo.IndexModel{
		"users":        {Keys: bson.M{"username": 1}, Options: options.Index().SetUnique(true).SetName("username_1")},
		"groups":       {Keys: bson.D{{Key: "groupname", Value: 1}, {Key: "project", Value: 1}}, Options: options.Index().SetUnique(true).SetName("groupname_project_1")},
		"service":      {Keys: bson.D{{Key: "name", Value: 1}, {Key: "project", Value: 1}}, Options: options.Index().SetUnique(true).SetName("name_project_1")},
		"clusters":     {Keys: bson.D{{Key: generalName, Value: 1}, {Key: generalProject, Value: 1}}, Options: options.Index().SetUnique(true).SetName(generalNameProject)},
		"listeners":    {Keys: bson.D{{Key: generalName, Value: 1}, {Key: generalProject, Value: 1}}, Options: options.Index().SetUnique(true).SetName(generalNameProject)},
		"endpoints":    {Keys: bson.D{{Key: generalName, Value: 1}, {Key: generalProject, Value: 1}}, Options: options.Index().SetUnique(true).SetName(generalNameProject)},
		"routes":       {Keys: bson.D{{Key: generalName, Value: 1}, {Key: generalProject, Value: 1}}, Options: options.Index().SetUnique(true).SetName(generalNameProject)},
		"virtual_host": {Keys: bson.D{{Key: generalName, Value: 1}, {Key: generalProject, Value: 1}}, Options: options.Index().SetUnique(true).SetName(generalNameProject)},
		"filters":      {Keys: bson.D{{Key: generalName, Value: 1}, {Key: generalProject, Value: 1}}, Options: options.Index().SetUnique(true).SetName(generalNameProject)},
		"secrets":      {Keys: bson.D{{Key: generalName, Value: 1}, {Key: generalProject, Value: 1}}, Options: options.Index().SetUnique(true).SetName(generalNameProject)},
		"extensions":   {Keys: bson.D{{Key: generalName, Value: 1}, {Key: generalProject, Value: 1}}, Options: options.Index().SetUnique(true).SetName(generalNameProject)},
		"bootstrap":    {Keys: bson.D{{Key: generalName, Value: 1}, {Key: generalProject, Value: 1}}, Options: options.Index().SetUnique(true).SetName(generalNameProject)},
		"projects":     {Keys: bson.M{"projectname": 1}, Options: options.Index().SetUnique(true).SetName("projectname_1")},
	}

	for collectionName, index := range indices {
		collection := database.Collection(collectionName)
		indexName := getIndexName(index)
		if err := createIndex(ctx, collection, index, indexName); err != nil {
			logger.Fatalf("Failed to create index for %s: %v\n", collectionName, err)
			return err
		}
	}

	return nil
}

func createIndex(ctx context.Context, collection *mongo.Collection, index mongo.IndexModel, indexName string) error {
	if indexName == "" {
		return errstr.ErrInvalidIndexName
	}
	exists, err := indexExists(ctx, collection, indexName)
	if err != nil {
		return fmt.Errorf("could not check for index existence: %w", err)
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
	var nameParts []string

	switch keys := index.Keys.(type) {
	case bson.M:
		for key, val := range keys {
			if nestedKeys, ok := val.(bson.M); ok {
				for nestedKey := range nestedKeys {
					nameParts = append(nameParts, key+"."+nestedKey+"_1")
				}
			} else {
				nameParts = append(nameParts, key+"_1")
			}
		}
	case bson.D:
		for _, keyVal := range keys {
			key := keyVal.Key
			if nestedKeys, ok := keyVal.Value.(bson.D); ok {
				for _, nestedKeyVal := range nestedKeys {
					nestedKey := nestedKeyVal.Key
					nameParts = append(nameParts, key+"."+nestedKey+"_1")
				}
			} else {
				nameParts = append(nameParts, key+"_1")
			}
		}
	default:
		return ""
	}

	return strings.Join(nameParts, "_")
}

func createAdminUser(db *AppContext) (string, error) {
	collection := db.Client.Collection("users")
	var user models.User
	err := collection.FindOne(db.Ctx, bson.M{"username": "admin"}).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		hashedPassword := helper.HashPassword("admin")
		user.Password = &hashedPassword
		now := time.Now()

		user.CreatedAt = primitive.NewDateTimeFromTime(now)
		user.UpdatedAt = primitive.NewDateTimeFromTime(now)
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()
		user.Email = &adminEmail
		user.Username = &adminUser
		user.Role = &adminRole
		user.BaseGroup = &adminBaseGroup
		user.Active = &adminActive

		token, refreshToken, _ := helper.GenerateAllTokens(user.Email, user.Username, user.UserID, nil, nil, nil, nil, false, user.Role)
		user.Token = &token
		user.RefreshToken = &refreshToken

		_, insertErr := collection.InsertOne(db.Ctx, user)
		if insertErr != nil {
			return "", insertErr
		}
	}
	return user.UserID, nil
}

func createAdminGroup(db *AppContext, userID string) error {
	if userID == "" {
		return errstr.ErrUserIDEmpty
	}

	collection := db.Client.Collection("groups")
	var group models.Group
	err := collection.FindOne(db.Ctx, bson.M{"groupname": "admin"}).Decode(&group)

	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		_, err = collection.InsertOne(db.Ctx, bson.M{
			"groupname":  "admin",
			"members":    []string{userID},
			"created_at": primitive.NewDateTimeFromTime(time.Now()),
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		})
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				db.Logger.Infof("admin group already exists: %v", err)
			} else {
				return fmt.Errorf("failed to create admin group: %w", err)
			}
		} else {
			db.Logger.Info("admin group created successfully")
		}
	case err != nil:
		return fmt.Errorf("failed to check for admin group: %w", err)
	default:
		db.Logger.Info("admin group already exists")
	}

	return nil
}

func createDefaultProject(db *AppContext, userID string) error {
	if userID == "" {
		return errstr.ErrUserIDEmpty
	}

	collection := db.Client.Collection("projects")
	var project models.Project
	err := collection.FindOne(db.Ctx, bson.M{"projectname": "default"}).Decode(&project)

	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		_, err = collection.InsertOne(db.Ctx, bson.M{
			"projectname": "default",
			"members":     []string{userID},
			"created_at":  primitive.NewDateTimeFromTime(time.Now()),
			"updated_at":  primitive.NewDateTimeFromTime(time.Now()),
		})
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				db.Logger.Infof("default project already exists: %v", err)
			} else {
				return fmt.Errorf("failed to create default project: %w", err)
			}
		} else {
			db.Logger.Info("default project created successfully")
		}
	case err != nil:
		return fmt.Errorf("failed to check for default project: %w", err)
	default:
		db.Logger.Info("default project already exists")
	}

	return nil
}
