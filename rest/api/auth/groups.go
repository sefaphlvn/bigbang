package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sefaphlvn/bigbang/pkg/errstr"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

type GroupWithActiveStatus struct {
	models.Group
	IsCreate    bool               `json:"is_create"`
	Permissions *models.Permission `json:"permissions"`
}

const (
	ErrUserNotCreated  = "User item was not created"
	SuccessUserCreated = "Successfully created user"
)

func (handler *AppHandler) ListGroups(c *gin.Context) {
	ctx := c.Request.Context()
	var groupCollection *mongo.Collection = handler.Context.Client.Collection("groups")
	filter := bson.M{"project": c.Query("project")}

	if !handler.CheckUserProjectPermission(c) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user does not have permission to view list of groups"})
		return
	}

	opts := options.Find().SetProjection(bson.M{"groupname": 1, "members": 1, "created_at": 1, "updated_at": 1})
	cursor, err := groupCollection.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not find records"})
	}

	var records []bson.M
	if err = cursor.All(ctx, &records); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not decode records"})
	}

	c.JSON(http.StatusOK, records)
}

func (handler *AppHandler) GetGroup(c *gin.Context) {
	ctx := c.Request.Context()
	var userCollection *mongo.Collection = handler.Context.Client.Collection("groups")
	groupID := c.Param("group_id")
	objectID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group_id"})
		return
	}

	if !handler.CheckUserProjectPermission(c) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user does not have permission to view group"})
		return
	}

	filter := bson.M{"_id": objectID, "project": c.Query("project")}

	opts := options.FindOne().SetProjection(bson.M{"groupname": 1, "email": 1, "created_at": 1, "updated_at": 1, "members": 1})
	var record bson.M
	err = userCollection.FindOne(ctx, filter, opts).Decode(&record)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not find records"})
	}

	c.JSON(http.StatusOK, record)
}

func (handler *AppHandler) GetBaseGroup(ctx context.Context, userID string) *string {
	var usersCollection *mongo.Collection = handler.Context.Client.Collection("users")
	filters := bson.M{"user_id": userID}
	opts := options.Find()
	opts.SetProjection(bson.M{"base_group": 1})
	cursor, err := usersCollection.Find(ctx, filters, opts)
	if err != nil {
		handler.Context.Logger.Info(err)
		return nil
	}
	defer cursor.Close(ctx)

	var result struct {
		BaseGroup *string `bson:"base_group"`
	}

	if cursor.Next(ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			handler.Context.Logger.Info(err)
			return nil
		}
		return result.BaseGroup
	}

	return nil
}

func (handler *AppHandler) GetUserGroups(ctx context.Context, userID string) (*[]string, *string, bool) {
	var groupCollection *mongo.Collection = handler.Context.Client.Collection("groups")
	filters := bson.M{"members": userID}
	adminGroup := false

	opts := options.Find()
	opts.SetProjection(bson.M{"_id": 1, "groupname": 1})
	cursor, err := groupCollection.Find(ctx, filters, opts)
	if err != nil {
		handler.Context.Logger.Info(err)
		return nil, nil, false
	}

	defer cursor.Close(ctx)

	var results []string
	for cursor.Next(ctx) {
		var group models.Group
		if err := cursor.Decode(&group); err != nil {
			handler.Context.Logger.Info(err)
			continue
		}
		results = append(results, group.ID.Hex())
		if group.GroupName != nil && *group.GroupName == "admin" {
			adminGroup = true
		}
	}

	baseGroup := handler.GetBaseGroup(ctx, userID)
	if baseGroup != nil {
		results = append(results, *baseGroup)
	}

	if err := cursor.Err(); err != nil {
		handler.Context.Logger.Info(err)
	}

	return helper.RemoveDuplicates(&results), baseGroup, adminGroup
}

func (handler *AppHandler) SetUpdateGroup(c *gin.Context) {
	ctx := c.Request.Context()
	var userCollection *mongo.Collection = handler.Context.Client.Collection("groups")
	ctx, cancel := context.WithTimeout(ctx, 100*time.Second)
	var status int
	var msg, groupID string
	defer cancel()
	var groupWA GroupWithActiveStatus

	if !handler.CheckUserProjectPermission(c) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user does not have permission to update group"})
		return
	}

	if err := c.BindJSON(&groupWA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if groupWA.IsCreate {
		status, msg, groupID = handler.CreateGroup(ctx, userCollection, groupWA)
	} else {
		status, msg = handler.UpdateGroup(ctx, userCollection, groupWA, c.Param("group_id"))
		groupID = c.Param("group_id")
	}

	if groupWA.Permissions != nil {
		handler.SetPermission(*groupWA.Permissions, groupID, "groups")
	}

	respondWithJSON(c, status, msg, groupID)
}

func (handler *AppHandler) CreateGroup(ctx context.Context, groupCollection *mongo.Collection, groupWA GroupWithActiveStatus) (int, string, string) {
	count, err := groupCollection.CountDocuments(ctx, bson.M{"groupname": groupWA.GroupName})
	if err != nil {
		return http.StatusBadRequest, "error occurred while checking for the groupname", "0"
	}

	if count > 0 {
		return http.StatusBadRequest, "groupname already exists", "0"
	}

	validationErr := validate.Struct(groupWA.Group)
	if validationErr != nil {
		return http.StatusBadRequest, validationErr.Error(), "0"
	}

	now := time.Now()

	groupWA.CreatedAt = primitive.NewDateTimeFromTime(now)
	groupWA.UpdatedAt = primitive.NewDateTimeFromTime(now)
	groupWA.ID = primitive.NewObjectID()

	_, insertErr := groupCollection.InsertOne(ctx, groupWA.Group)

	if insertErr != nil {
		return http.StatusBadRequest, ErrUserNotCreated, "0"
	}

	return http.StatusOK, SuccessUserCreated, groupWA.ID.String()
}

func (handler *AppHandler) UpdateGroup(ctx context.Context, groupCollection *mongo.Collection, groupWA GroupWithActiveStatus, groupID string) (int, string) {
	objectID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return http.StatusBadRequest, "no group found with the given group_id"
	}
	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{},
	}

	updateMap, ok := update["$set"].(bson.M)
	if !ok {
		return http.StatusInternalServerError, errstr.ErrUnexpectedTypeBsonM.Error()
	}

	if groupWA.GroupName != nil {
		updateMap["groupname"] = groupWA.GroupName
	}

	if groupWA.Members != nil {
		updateMap["members"] = groupWA.Members
	}

	updateMap["updated_at"] = primitive.NewDateTimeFromTime(time.Now())
	result, err := groupCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("error updating group: %v", err)
	}

	if result.MatchedCount == 0 {
		return http.StatusBadRequest, "no group found with the given groupname"
	}

	return http.StatusOK, "group successfully updated"
}
