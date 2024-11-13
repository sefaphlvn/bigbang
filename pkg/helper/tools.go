package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/sefaphlvn/bigbang/pkg/models"
)

var Unmarshaler = protojson.UnmarshalOptions{
	AllowPartial:   true,
	DiscardUnknown: true,
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func PrettyPrint(data interface{}) {
	// Eğer veri nil ise çıkış yap
	if data == nil {
		return
	}

	// JSON verisini saklamak için
	var jsonData interface{}

	switch v := data.(type) {
	case string:
		// String'i JSON olarak unmarshall etmeyi dene
		if err := json.Unmarshal([]byte(v), &jsonData); err != nil {
			// JSON değilse, doğrudan string'i döndür
			fmt.Println(v)
			return
		}
	default:
		// JSON olmayan yapıları doğrudan al
		jsonData = v
	}

	// Eğer jsonData bir JSON yapıdaysa, pretty print yap
	prettyJSON, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		log.Fatalf("JSON marshaling error: %v", err)
	}

	fmt.Println(string(prettyJSON))
}

/* func GetResourceType(data interface{}) {
	resourceType := reflect.TypeOf(data)
	fmt.Printf("Resource type: %v\n", resourceType)
} */

func ToBool(strBool string) bool {
	boolean, err := strconv.ParseBool(strBool)
	if err != nil {
		fmt.Println(err)
	}

	return boolean
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func GenerateAllTokens(email, username *string, userID string, groups *[]string, projects *[]models.CombinedProjects, baseGroup, baseProject *string, adminGroup bool, role *models.Role) (signedToken, signedRefreshToken string, err error) {
	claims := &models.SignedDetails{
		Email:       email,
		Username:    username,
		UserID:      userID,
		Groups:      RemoveDuplicates(groups),
		Projects:    projects,
		BaseGroup:   baseGroup,
		BaseProject: baseProject,
		AdminGroup:  adminGroup,
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
	}

	refreshClaims := &models.SignedDetails{
		Email:       email,
		Username:    username,
		UserID:      userID,
		Groups:      RemoveDuplicates(groups),
		Projects:    projects,
		BaseGroup:   baseGroup,
		BaseProject: baseProject,
		AdminGroup:  adminGroup,
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SecretKey))
	if err != nil {
		log.Fatal(err)
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SecretKey))
	if err != nil {
		log.Fatal(err)
		return "", "", err
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
		logger.Debugf("Error marshaling JSON: %v", err)
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

func MarshalUnmarshalWithType(data interface{}, msg proto.Message) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = Unmarshaler.Unmarshal(jsonData, msg)
	if err != nil {
		return err
	}

	return nil
}

func ConvertToJSON(v interface{}, log *logrus.Logger) string {
	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Infof("JSON convert err: %v", err)
	}
	return string(jsonData)
}

func EscapePointKey(key string) string {
	return strings.ReplaceAll(key, ".", `\.`)
}
