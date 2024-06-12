package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func PrettyPrinter(data interface{}) {
	prettyJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatalf("JSON marshaling error: %v", err)
	}

	fmt.Println(string(prettyJSON))
}

func GetResourceType(data interface{}) {
	resourceType := reflect.TypeOf(data)
	fmt.Printf("Resource type: %v\n", resourceType)
}

func ToBool(strBool string) bool {
	boolean, err := strconv.ParseBool(strBool)
	if err != nil {
		fmt.Println(err)
	}

	return boolean
}

func ToInt(strInt string) int {
	number, err := strconv.Atoi(strInt)
	if err != nil {
		fmt.Printf("Cannot convert to int: %s", err)
		return 0
	}

	return number
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func GenerateAllTokens(email string, Username string, user_id string, groups []string, base_group *string, adminGroup bool, role string) (signedToken string, signedRefreshToken string, err error) {
	claims := &models.SignedDetails{
		Email:      email,
		Username:   Username,
		UserId:     user_id,
		Groups:     RemoveDuplicates(groups),
		BaseGroup:  base_group,
		AdminGroup: adminGroup,
		Role:       role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(60)).Unix(),
		},
	}

	refreshClaims := &models.SignedDetails{
		Email:      email,
		Username:   Username,
		UserId:     user_id,
		Groups:     RemoveDuplicates(groups),
		BaseGroup:  base_group,
		AdminGroup: adminGroup,
		Role:       role,
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

func RemoveDuplicates(strings []string) []string {
	uniqueStrings := make(map[string]bool)
	result := []string{}

	for _, str := range strings {
		if _, exists := uniqueStrings[str]; !exists {
			uniqueStrings[str] = true
			result = append(result, str)
		}
	}

	return result
}
