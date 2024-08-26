package common

import (
	"github.com/sefaphlvn/bigbang/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

func AddUserFilter(details models.RequestDetails, mainFilter bson.M) bson.M {
	userFilter := bson.M{"general.project": details.Project}
	if !details.User.IsOwner && details.User.Role != models.RoleAdmin {
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

	return mainFilter
}
