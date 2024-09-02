package filters

import (
	"go.mongodb.org/mongo-driver/bson"
)

type MongoFilters struct {
	Collection string
	Filter     bson.D
}
