package common

import (
	"github.com/sefaphlvn/bigbang/rest/models"
	"go.mongodb.org/mongo-driver/bson"
)

func AddUserFilter(details models.ResourceDetails, mainFilter bson.M) bson.M {
	// Kullanıcıya bağlı filtre oluştur
	userFilter := bson.M{}
	if !details.User.IsAdmin {
		userFilter = bson.M{"general.groups": bson.M{"$in": details.User.Groups}}
	}

	// Mevcut filtreyi kullanıcı filtreyi ile birleştir
	for key, value := range userFilter {
		mainFilter[key] = value
	}

	return mainFilter
}
