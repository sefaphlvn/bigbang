package common

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sefaphlvn/bigbang/pkg/models"
)

func TransformGenerals(records []bson.M) interface{} {
	generals := make([]models.General, 0, len(records))

	for _, record := range records {
		bsonData, err := bson.Marshal(record["general"])
		if err != nil {
			return nil
		}

		var general models.General
		if err := bson.Unmarshal(bsonData, &general); err != nil {
			return nil
		}

		generals = append(generals, general)
	}
	return generals
}
