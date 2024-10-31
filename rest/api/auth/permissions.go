package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	general_name = "general.name"
)

func (handler *AppHandler) GetPermissions(c *gin.Context) {
	project := c.Query("project")
	var userOrGroup = c.Param("kind")
	var filter = bson.M{"general.project": project}

	if userOrGroup == "users" {
		filter["general.permissions.users"] = c.Param("id")
	} else {
		filter["general.permissions.groups"] = c.Param("id")
	}

	all, err := handler.GetData(bson.M{"general.project": project}, c.Param("type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	selected, err := handler.GetData(filter, c.Param("type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	records := map[string]interface{}{"all": all, "selected": selected}
	c.JSON(http.StatusOK, records)
}

func (handler *AppHandler) SetPermission(permissions models.Permission, userOrGroupID string, kind string) {
	checkAndAct := func(name string, p *models.InnerPermission) {
		if p != nil {
			if len(p.Added) > 0 || len(p.Removed) > 0 {
				collection := handler.Context.Client.Collection(name)
				for _, addedName := range p.Added {
					filter := bson.M{general_name: addedName}
					update := bson.M{
						"$addToSet": bson.M{"general.permissions." + kind: userOrGroupID},
					}
					_, err := collection.UpdateMany(context.TODO(), filter, update)
					if err != nil {
						handler.Context.Logger.Warnf("failed to add user/group to %s: %v", name, err)
					}
					fmt.Printf("User/Group %s added to %s: %v\n", userOrGroupID, name, addedName)
				}
				for _, removedName := range p.Removed {
					filter := bson.M{general_name: removedName}
					update := bson.M{
						"$pull": bson.M{"general.permissions." + kind: userOrGroupID},
					}
					_, err := collection.UpdateMany(context.TODO(), filter, update)
					if err != nil {
						handler.Context.Logger.Warnf("failed to remove user/group from %s: %v", name, err)
					}
					fmt.Printf("User/Group %s removed from %s: %v\n", userOrGroupID, name, removedName)
				}
			}
		}
	}

	fields := map[string]*models.InnerPermission{
		"listeners":  permissions.Listeners,
		"routes":     permissions.Routes,
		"clusters":   permissions.Clusters,
		"endpoints":  permissions.Endpoints,
		"secrets":    permissions.Secrets,
		"extensions": permissions.Extensions,
		"filters":    permissions.Filters,
		"bootstrap":  permissions.Bootstrap,
	}

	for name, field := range fields {
		checkAndAct(name, field)
	}
}

func (handler *AppHandler) GetData(filter bson.M, typ string) ([]bson.M, error) {
	var resourceCollection *mongo.Collection = handler.Context.Client.Collection(typ)
	opts := options.Find().SetProjection(bson.M{general_name: 1})

	cursor, err := resourceCollection.Find(handler.Context.Ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var records []bson.M
	if err = cursor.All(handler.Context.Ctx, &records); err != nil {
		return nil, err
	}

	return records, err
}
