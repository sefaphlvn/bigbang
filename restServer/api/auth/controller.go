package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/restServer/crud"
	"github.com/sefaphlvn/bigbang/restServer/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

type DBHandler crud.DbHandler

type SignedDetails struct {
	Email    string
	Username string
	User_id  string
	Groups   []string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("secret")

var validate = validator.New()

func NewUserHandler(db *db.MongoDB) *DBHandler {
	return &DBHandler{
		DB: db,
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
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

func (userDB *DBHandler) SignUp() gin.HandlerFunc {
	var userCollection *mongo.Collection = userDB.DB.Client.Collection("user")
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error occured while checking for the username"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "username already exists"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password
		now := time.Now()

		user.Created_at = primitive.NewDateTimeFromTime(now)
		user.Updated_at = primitive.NewDateTimeFromTime(now)
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		user.Groups = []string{}
		token, refreshToken, _ := GenerateAllTokens(*user.Email, *user.Username, user.User_id, user.Groups)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := "User item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
			return
		}

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func (userDB *DBHandler) Login() gin.HandlerFunc {
	var userCollection *mongo.Collection = userDB.DB.Client.Collection("user")
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

		token, refreshToken, _ := GenerateAllTokens(*foundUser.Email, *foundUser.Username, foundUser.User_id, foundUser.Groups)

		foundUser.Token = &token
		foundUser.Refresh_token = &refreshToken

		UpdateAllTokens(userDB, token, refreshToken, foundUser.User_id)

		c.JSON(http.StatusOK, foundUser)
	}
}

func (userDB *DBHandler) Logout() gin.HandlerFunc {
	var userCollection *mongo.Collection = userDB.DB.Client.Collection("user")
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Could not retrieve user id"})
			c.Abort()
			return
		}

		fmt.Println(userId)
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
	var userCollection *mongo.Collection = userDB.DB.Client.Collection("user")
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

		signedToken, signedRefreshToken, _ := GenerateAllTokens(*foundUser.Email, *foundUser.Username, foundUser.User_id, foundUser.Groups)
		UpdateAllTokens(userDB, signedToken, signedRefreshToken, foundUser.User_id)

		c.JSON(http.StatusOK, gin.H{
			"token":         signedToken,
			"refresh_token": signedRefreshToken,
		})
	}
}

func GenerateAllTokens(email string, Username string, user_id string, groups []string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:    email,
		Username: Username,
		User_id:  user_id,
		Groups:   groups,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(60)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		Email:    email,
		Username: Username,
		User_id:  user_id,
		Groups:   groups,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Fatal(err)
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}

	return claims, msg
}

func UpdateAllTokens(userDB *DBHandler, signedToken string, signedRefreshToken string, userId string) {
	var userCollection *mongo.Collection = userDB.DB.Client.Collection("user")
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

func ValidateRefreshToken(tokenString string) (SignedDetails, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv(SECRET_KEY)), nil
		},
	)
	if err != nil {
		return SignedDetails{}, fmt.Errorf("could not parse refresh token: %w", err)
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok || !token.Valid {
		return SignedDetails{}, fmt.Errorf("invalid refresh token")
	}

	return *claims, nil
}
