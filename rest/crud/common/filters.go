package common

import (
	"github.com/sefaphlvn/bigbang/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

func AddUserFilter(details models.ResourceDetails, mainFilter bson.M) bson.M {
	userFilter := bson.M{}
	if !details.User.IsAdmin {
		userFilter = bson.M{"general.groups": bson.M{"$in": details.User.Groups}}
	}

	for key, value := range userFilter {
		mainFilter[key] = value
	}

	return mainFilter
}
