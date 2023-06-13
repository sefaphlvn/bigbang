package helper

import (
	"encoding/json"
	"log"
)

func ToInterface(data []byte) map[string]interface{} {
	var dataInterface map[string]interface{}
	err := json.Unmarshal(data, &dataInterface)
	if err != nil {
		log.Fatal(err)
	}

	return dataInterface
}
