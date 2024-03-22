package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func GetString(m bson.M, key string) string {
	if m != nil {
		if str, ok := m[key].(string); ok {
			return str
		}
	}
	return "None"
}

func GetBool(m bson.M, key string) bool {
	b, ok := m[key].(bool)
	if !ok {
		b = false
	}

	return b
}

func GetDateTime(m bson.M, key string) primitive.DateTime {
	dt, _ := m[key].(primitive.DateTime)
	return dt
}

func ToMapStringInterface(data []byte) map[string]interface{} {
	var typedData map[string]interface{}
	err := json.Unmarshal(data, &typedData)
	if err != nil {
		log.Fatal(err)
	}

	return typedData
}

func ItoGenericTypeConvert[T any](data interface{}) T {
	typedData, ok := data.(T)
	if !ok {
		log.Fatal("invalid type")
	}

	return typedData
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
		fmt.Println("Hata: MongoDB_Timeout değeri integer'a çevrilemiyor.")
		return 0
	}

	return number
}
