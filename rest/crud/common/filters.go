package common

import (
	"fmt"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

func AddUserFilter(details models.ResourceDetails, mainFilter bson.M) bson.M {
	userFilter := bson.M{}
	if !details.User.IsAdmin {
		userFilter = bson.M{
			"$or": []bson.M{
				{"general.permissions.groups": bson.M{"$in": details.User.Groups}},
				{"general.permissions.users": details.User.UserID},
			},
		}
	}

	for key, value := range userFilter {
		mainFilter[key] = value
	}

	fmt.Println(mainFilter)
	return mainFilter
}
