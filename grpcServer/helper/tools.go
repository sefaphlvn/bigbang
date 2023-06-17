package helper

import (
	"encoding/json"
	"log"
)

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
