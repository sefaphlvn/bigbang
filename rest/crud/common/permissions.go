package common

import (
	"github.com/sefaphlvn/bigbang/pkg/models"
)

func DetectSetPermissions(resource models.DBResourceClass, requestDetails models.RequestDetails) {
	resource.SetPermissions(getPermissions(requestDetails))
}

func getPermissions(requestDetails models.RequestDetails) *models.Permissions {
	if requestDetails.User.BaseGroup != "" {
		return &models.Permissions{Groups: []string{requestDetails.User.BaseGroup}, Users: []string{}}
	}
	if requestDetails.User.UserID != "" {
		return &models.Permissions{Groups: []string{}, Users: []string{requestDetails.User.UserID}}
	}
	return &models.Permissions{Groups: []string{}, Users: []string{}}
}
