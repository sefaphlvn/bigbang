package helper

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
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
	s, _ := m[key].(string)
	return s
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
