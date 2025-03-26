package common

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/sefaphlvn/bigbang/pkg/models"
)

func AddUserFilter(details models.RequestDetails, mainFilter bson.M) bson.M {
	if mainFilter == nil {
		mainFilter = bson.M{}
	}

	userFilter := bson.M{}
	if !details.User.IsOwner && details.User.Role != models.RoleAdmin {
		userFilter = bson.M{
			"$or": []bson.M{
				{"general.permissions.groups": bson.M{"$in": details.User.Groups}},
				{"general.permissions.users": details.User.UserID},
			},
		}
	}

	mainFilter["general.project"] = details.Project

	for key, value := range userFilter {
		mainFilter[key] = value
	}

	return mainFilter
}

func AddResourceIDFilter(requestDetails models.RequestDetails, mainFilter bson.M) (bson.M, error) {
	if mainFilter == nil {
		mainFilter = bson.M{}
	}

	objectID, err := primitive.ObjectIDFromHex(requestDetails.ResourceID)
	if err != nil {
		return mainFilter, err
	}

	mainFilter["_id"] = objectID
	fmt.Println(mainFilter)
	return mainFilter, nil
}
