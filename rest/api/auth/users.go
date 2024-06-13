package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserWithGroups struct {
	models.User
	Groups      []string           `json:"groups"`
	IsCreate    bool               `json:"is_create"`
	Permissions *models.Permission `json:"permissions"`
}

func (userDB *DBHandler) SetUpdateUser(c *gin.Context) {
	var userCollection *mongo.Collection = userDB.DB.Client.Collection("users")
	var ctx, cancel = context.WithTimeout(userDB.DB.Ctx, 100*time.Second)
	var status int
	var msg, userID string
	var userWG UserWithGroups
	defer cancel()

	if err := c.BindJSON(&userWG); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if userWG.IsCreate {
		status, msg, userID = userDB.CreateUser(ctx, userCollection, userWG)
	} else {
		status, msg = userDB.UpdateUser(ctx, userCollection, userWG, c.Param("user_id"))
		userID = c.Param("user_id")
	}

	if userWG.Permissions != nil {
		userDB.SetPermission(*userWG.Permissions, userID, "users")
	}

	respondWithJSON(c, status, msg, userID)
}

func (userDB *DBHandler) CreateUser(ctx context.Context, userCollection *mongo.Collection, userWG UserWithGroups) (int, string, string) {
	count, err := userCollection.CountDocuments(ctx, bson.M{"username": userWG.Username})
	if err != nil {
		return http.StatusBadRequest, "error occured while checking for the username", "0"
	}

	if count > 0 {
		return http.StatusBadRequest, "username already exists", "0"
	}

	validationErr := validate.Struct(userWG.User)
	if validationErr != nil {
		return http.StatusBadRequest, validationErr.Error(), "0"
	}

	password := helper.HashPassword(*userWG.Password)
	userWG.Password = &password
	now := time.Now()

	userWG.Created_at = primitive.NewDateTimeFromTime(now)
	userWG.Updated_at = primitive.NewDateTimeFromTime(now)
	userWG.ID = primitive.NewObjectID()
	userWG.User_id = userWG.ID.Hex()
	token, refreshToken, _ := helper.GenerateAllTokens(*userWG.Email, *userWG.Username, userWG.User_id, []string{}, nil, false, *userWG.Role)
	userWG.Token = &token
	userWG.Refresh_token = &refreshToken

	insertOneResult, insertErr := userCollection.InsertOne(ctx, userWG.User)

	if insertErr != nil {
		return http.StatusBadRequest, "User item was not created", userWG.User_id
	}

	if userWG.Groups != nil {
		fmt.Println(insertOneResult.InsertedID, userWG.Groups)
	}

	return http.StatusOK, "Successfully created user", userWG.User_id
}

func (userDB *DBHandler) UpdateUser(ctx context.Context, userCollection *mongo.Collection, userWG UserWithGroups, userID string) (int, string) {
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{},
	}

	if userWG.Username != nil {
		update["$set"].(bson.M)["username"] = userWG.Username
	}
	if userWG.Password != nil {
		hashedPassword := helper.HashPassword(*userWG.Password)
		update["$set"].(bson.M)["password"] = hashedPassword
	}
	if userWG.Email != nil {
		update["$set"].(bson.M)["email"] = userWG.Email
	}
	if userWG.Role != nil {
		update["$set"].(bson.M)["role"] = userWG.Role
	}
	if userWG.BaseGroup != nil {
		if *userWG.BaseGroup == "xremove" {
			update["$set"].(bson.M)["base_group"] = nil
		} else {
			update["$set"].(bson.M)["base_group"] = userWG.BaseGroup
		}
	}

	if userWG.Active != nil {
		update["$set"].(bson.M)["active"] = userWG.Active
	}

	update["$set"].(bson.M)["updated_at"] = primitive.NewDateTimeFromTime(time.Now())
	result, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("error updating user: %v", err)
	}

	if result.MatchedCount == 0 {
		return http.StatusBadRequest, "no user found with the given username"
	}

	return http.StatusOK, "user successfully updated"
}

func (userDB *DBHandler) Login() gin.HandlerFunc {
	var userCollection *mongo.Collection = userDB.DB.Client.Collection("users")
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "username or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*foundUser.Password, *user.Password)

		if !passwordIsValid {
			c.JSON(http.StatusBadRequest, gin.H{"message": msg})
			return
		}

		groups, base_group, adminGroup := userDB.GetUserGroups(foundUser.User_id)

		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.Username, foundUser.User_id, groups, base_group, adminGroup, *foundUser.Role)

		foundUser.Token = &token
		foundUser.Refresh_token = &refreshToken

		UpdateAllTokens(userDB, token, refreshToken, foundUser.User_id)

		c.JSON(http.StatusOK, foundUser)
	}
}

func (userDB *DBHandler) ListUsers(c *gin.Context) {
	var userCollection *mongo.Collection = userDB.DB.Client.Collection("users")
	filter := bson.M{}

	opts := options.Find().SetProjection(bson.M{"username": 1, "email": 1, "created_at": 1, "updated_at": 1, "user_id": 1, "groups": 1})
	cursor, err := userCollection.Find(userDB.DB.Ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not find records"})
	}

	var records []bson.M
	if err = cursor.All(userDB.DB.Ctx, &records); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not decode records"})
	}

	c.JSON(http.StatusOK, records)
}

func (userDB *DBHandler) GetUser(c *gin.Context) {
	var userCollection *mongo.Collection = userDB.DB.Client.Collection("users")
	filter := bson.M{"user_id": c.Param("user_id")}

	opts := options.FindOne().SetProjection(bson.M{"username": 1, "email": 1, "created_at": 1, "updated_at": 1, "user_id": 1, "groups": 1, "role": 1, "base_group": 1, "active": 1})
	var record bson.M
	err := userCollection.FindOne(userDB.DB.Ctx, filter, opts).Decode(&record)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not find records"})
	}

	groups, _, _ := userDB.GetUserGroups(record["user_id"].(string))
	record["groups"] = groups

	c.JSON(http.StatusOK, record)
}

func (userDB *DBHandler) Logout() gin.HandlerFunc {
	var userCollection *mongo.Collection = userDB.DB.Client.Collection("users")
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Could not retrieve user id"})
			c.Abort()
			return
		}

		filter := bson.M{"user_id": userId}
		update := bson.M{
			"$unset": bson.M{
				"token":         "",
				"refresh_token": "",
			},
		}

		_, err := userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to logout"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	}
}

func (userDB *DBHandler) Refresh() gin.HandlerFunc {
	var userCollection *mongo.Collection = userDB.DB.Client.Collection("users")
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var token models.User

		if err := c.BindJSON(&token); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.M{"refresh-token": token.Refresh_token}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid refresh token"})
			return
		}

		groups, base_group, admin_group := userDB.GetUserGroups(foundUser.User_id)

		signedToken, signedRefreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.Username, foundUser.User_id, groups, base_group, admin_group, *foundUser.Role)
		UpdateAllTokens(userDB, signedToken, signedRefreshToken, foundUser.User_id)

		c.JSON(http.StatusOK, gin.H{
			"token":         signedToken,
			"refresh_token": signedRefreshToken,
		})
	}
}
