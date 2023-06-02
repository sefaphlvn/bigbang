package extension

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/restApi/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (extension *DBHandler) SetExtension(resource models.DBResourceClass, collectionName models.ResourceDetails) (interface{}, error) {
	general := resource.GetGeneral()
	now := time.Now()
	general.CreatedAt = primitive.NewDateTimeFromTime(now)
	general.UpdatedAt = primitive.NewDateTimeFromTime(now)
	resource.SetGeneral(&general)

	collection := extension.DB.Client.Collection("extensions")
	_, err := collection.InsertOne(extension.DB.Ctx, resource)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("name already exists")
		}
		return nil, err
	}
	return gin.H{"message": "Success"}, nil
}
