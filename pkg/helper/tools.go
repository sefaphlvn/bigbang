package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sirupsen/logrus"
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

func PrettyPrint(data interface{}) {
	var jsonData interface{}

	// Eğer veri string ise JSON olarak unmarshall et
	switch v := data.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &jsonData); err != nil {
			log.Fatalf("JSON unmarshaling error: %v", err)
		}
	default:
		jsonData = v
	}

	// JSON verisini pretty print yap
	prettyJSON, err := json.MarshalIndent(jsonData, "", "    ")
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

func GenerateAllTokens(email *string, Username *string, user_id string, groups *[]string, projects *[]models.CombinedProjects, base_group *string, base_project *string, adminGroup bool, role *models.Role) (signedToken string, signedRefreshToken string, err error) {
	claims := &models.SignedDetails{
		Email:       email,
		Username:    Username,
		UserId:      user_id,
		Groups:      RemoveDuplicates(groups),
		Projects:    projects,
		BaseGroup:   base_group,
		BaseProject: base_project,
		AdminGroup:  adminGroup,
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
	}

	refreshClaims := &models.SignedDetails{
		Email:       email,
		Username:    Username,
		UserId:      user_id,
		Groups:      RemoveDuplicates(groups),
		Projects:    projects,
		BaseGroup:   base_group,
		BaseProject: base_project,
		AdminGroup:  adminGroup,
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
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

func RemoveDuplicates(strings *[]string) *[]string {
	if strings == nil {
		result := []string{}
		return &result
	}

	uniqueStrings := make(map[string]bool)
	result := []string{}

	for _, str := range *strings {
		if _, exists := uniqueStrings[str]; !exists {
			uniqueStrings[str] = true
			result = append(result, str)
		}
	}

	return &result
}

func MarshalJSON(data interface{}, logger *logrus.Logger) (string, error) {
	jsonString, err := json.Marshal(data)
	if err != nil {
		logger.Debugf("Error marshalling JSON: %v", err)
		return "", err
	}
	return string(jsonString), nil
}

func RemoveDuplicatesP(projects *[]models.CombinedProjects) *[]models.CombinedProjects {
	uniqueProjects := make(map[string]models.CombinedProjects)
	for _, project := range *projects {
		if _, exists := uniqueProjects[project.ProjectID]; !exists {
			uniqueProjects[project.ProjectID] = project
		}
	}

	result := make([]models.CombinedProjects, 0, len(uniqueProjects))
	for _, project := range uniqueProjects {
		result = append(result, project)
	}

	return &result
}
