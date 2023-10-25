package helper

import (
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