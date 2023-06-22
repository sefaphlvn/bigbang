package xds

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/restServer/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (xds *DBHandler) SetResource(resource models.DBResourceClass, collectionName models.ResourceDetails) (interface{}, error) {
	general := resource.GetGeneral()
	now := time.Now()
	general.CreatedAt = primitive.NewDateTimeFromTime(now)
	general.UpdatedAt = primitive.NewDateTimeFromTime(now)
	resource.SetGeneral(&general)

	fmt.Println(general)

	collection := xds.DB.Client.Collection(collectionName.Type)
	_, err := collection.InsertOne(xds.DB.Ctx, resource)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("name already exists")
		}
		return nil, err
	}
	return gin.H{"message": "Success"}, nil
}
