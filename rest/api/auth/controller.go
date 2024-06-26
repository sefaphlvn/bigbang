package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type DBHandler crud.DbHandler

var validate = validator.New()

func NewUserHandler(db *db.WTF) *DBHandler {
	return &DBHandler{
		DB: db,
	}
}

func VerifyPassword(userHashedPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userHashedPassword), []byte(providedPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "username or password is incorrect"
		check = false
	}

	return check, msg
}

func ValidateToken(signedToken string) (claims *models.SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&models.SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(helper.SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*models.SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}

	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}

	return claims, msg
}

func UpdateAllTokens(userDB *DBHandler, signedToken string, signedRefreshToken string, userId string) {
	var userCollection *mongo.Collection = userDB.DB.Client.Collection("users")
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	now := time.Now()

	Updated_at := primitive.NewDateTimeFromTime(now)
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)

	if err != nil {
		log.Panic(err)
		return
	}
}

func ValidateRefreshToken(tokenString string) (models.SignedDetails, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&models.SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv(helper.SECRET_KEY)), nil
		},
	)
	if err != nil {
		return models.SignedDetails{}, fmt.Errorf("could not parse refresh token: %w", err)
	}

	claims, ok := token.Claims.(*models.SignedDetails)
	if !ok || !token.Valid {
		return models.SignedDetails{}, fmt.Errorf("invalid refresh token")
	}

	return *claims, nil
}

func respondWithJSON(c *gin.Context, status int, msg, userOrGroupID string) {
	c.JSON(status, gin.H{"message": msg, "id": userOrGroupID})
}
