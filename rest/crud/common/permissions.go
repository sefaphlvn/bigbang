package common

import (
	"github.com/sefaphlvn/bigbang/pkg/models"
)

func DetectSetPermissions(resource models.DBResourceClass, requestDetails models.RequestDetails) {
	var permission = models.Permissions{Users: []string{}, Groups: []string{}}
	if requestDetails.User.Role == models.RoleAdmin || requestDetails.User.IsOwner {
		resource.SetPermissions(&permission)
	} else {
		resource.SetPermissions(getPermissions(requestDetails))
	}
}

func getPermissions(requestDetails models.RequestDetails) *models.Permissions {
	if requestDetails.User.BaseGroup != "" {
		return &models.Permissions{Groups: []string{requestDetails.User.BaseGroup}, Users: []string{}}
	} else if requestDetails.User.UserID != "" {
		return &models.Permissions{Groups: []string{}, Users: []string{requestDetails.User.UserID}}
	} else {
		return &models.Permissions{Groups: []string{}, Users: []string{}}
	}
}
