package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProjectWithActiveStatus struct {
	models.Project
	IsCreate bool `json:"is_create"`
}

func (handler *AppHandler) ListProjects(c *gin.Context) {
	UserID, _ := c.Get("user_id")
	userId, ok := UserID.(string)
	if !ok {
		userId = ""
	}
	var projectCollection *mongo.Collection = handler.Context.Client.Collection("projects")
	projects, _ := handler.GetUserProject(userId)

	var projectIDs []primitive.ObjectID
	for _, projectID := range *projects {
		objID, err := primitive.ObjectIDFromHex(projectID.ProjectID)
		if err != nil {
			handler.Context.Logger.Infof("Invalid ObjectID: %s", projectID.ProjectID)
			continue
		}
		projectIDs = append(projectIDs, objID)
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": projectIDs,
		},
	}

	opts := options.Find().SetProjection(bson.M{"projectname": 1, "members": 1, "created_at": 1, "updated_at": 1})
	cursor, err := projectCollection.Find(handler.Context.Ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not find records"})
	}

	var records []bson.M
	if err = cursor.All(handler.Context.Ctx, &records); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not decode records"})
	}

	c.JSON(http.StatusOK, records)
}

func (handler *AppHandler) GetProject(c *gin.Context) {
	var userCollection *mongo.Collection = handler.Context.Client.Collection("projects")
	var record bson.M

	projectID := c.Param("project_id")
	objectID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project_id"})
		return
	}

	filter := bson.M{"_id": objectID}
	opts := options.FindOne().SetProjection(bson.M{"projectname": 1, "email": 1, "created_at": 1, "updated_at": 1, "members": 1})
	err = userCollection.FindOne(handler.Context.Ctx, filter, opts).Decode(&record)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not find records"})
	}

	c.JSON(http.StatusOK, record)
}

func (handler *AppHandler) GetBaseProjectAndRole(userID string) (*string, bool) {
	var usersCollection *mongo.Collection = handler.Context.Client.Collection("users")
	var filters = bson.M{"user_id": userID}
	opts := options.Find()
	opts.SetProjection(bson.M{"base_project": 1, "username": 1, "role": 1})
	cursor, err := usersCollection.Find(handler.Context.Ctx, filters, opts)
	if err != nil {
		handler.Context.Logger.Info(err)
		return nil, false
	}
	defer cursor.Close(handler.Context.Ctx)

	var result struct {
		BaseProject *string `bson:"base_project"`
		UserName    string  `bson:"username"`
		Role        string  `bson:"role"`
	}

	if cursor.Next(handler.Context.Ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			handler.Context.Logger.Info(err)
			return nil, false
		}
		isOwner := result.Role == "owner"
		return result.BaseProject, isOwner
	}

	return nil, false
}

func (handler *AppHandler) GetUserProject(userID string) (*[]models.CombinedProjects, *string) {
	projectCollection := handler.Context.Client.Collection("projects")
	var projects []models.CombinedProjects

	baseProject, isOwner := handler.GetBaseProjectAndRole(userID)
	if baseProject != nil {
		projects = append(projects, models.CombinedProjects{
			ProjectID:   *baseProject,
			ProjectName: handler.getProjectName(projectCollection, *baseProject),
		})
	}

	if isOwner {
		allProjects, err := handler.getAllProjectNamesAndIDs(projectCollection)
		if err != nil {
			handler.Context.Logger.Info("Error getting all project names and IDs:", err)
			return nil, baseProject
		}
		projects = append(projects, allProjects...)
	} else {
		filters := bson.M{"members": userID}
		opts := options.Find().SetProjection(bson.M{"_id": 1, "projectname": 1})
		cursor, err := projectCollection.Find(handler.Context.Ctx, filters, opts)
		if err != nil {
			handler.Context.Logger.Info("Error finding projects:", err)
			return nil, baseProject
		}
		defer cursor.Close(handler.Context.Ctx)

		for cursor.Next(handler.Context.Ctx) {
			var project models.Project
			if err := cursor.Decode(&project); err != nil {
				handler.Context.Logger.Info("Error decoding project:", err)
				continue
			}
			projects = append(projects, models.CombinedProjects{
				ProjectID:   project.ID.Hex(),
				ProjectName: *project.ProjectName,
			})
		}

		if err := cursor.Err(); err != nil {
			handler.Context.Logger.Info("Cursor error:", err)
		}
	}

	return helper.RemoveDuplicatesP(&projects), baseProject
}

func (handler *AppHandler) SetUpdateProject(c *gin.Context) {
	var userCollection *mongo.Collection = handler.Context.Client.Collection("projects")
	var ctx, cancel = context.WithTimeout(handler.Context.Ctx, 100*time.Second)
	var status int
	var msg, projectID string
	defer cancel()
	var projectWA ProjectWithActiveStatus

	if err := c.BindJSON(&projectWA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if projectWA.IsCreate {
		if !c.GetBool("isOwner") {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to create a project"})
			return
		}
		status, msg, projectID = handler.CreateProject(ctx, userCollection, projectWA)

	} else {
		status, msg = handler.UpdateProject(ctx, userCollection, projectWA, c.Param("project_id"))
		projectID = c.Param("project_id")
	}

	respondWithJSON(c, status, msg, projectID)
}

func (handler *AppHandler) CreateProject(ctx context.Context, projectCollection *mongo.Collection, projectWA ProjectWithActiveStatus) (int, string, string) {
	count, err := projectCollection.CountDocuments(ctx, bson.M{"projectname": projectWA.ProjectName})
	if err != nil {
		return http.StatusBadRequest, "error occured while checking for the projectname", "0"
	}

	if count > 0 {
		return http.StatusBadRequest, "projectname already exists", "0"
	}

	validationErr := validate.Struct(projectWA.Project)
	if validationErr != nil {
		return http.StatusBadRequest, validationErr.Error(), "0"
	}

	now := time.Now()

	projectWA.Created_at = primitive.NewDateTimeFromTime(now)
	projectWA.Updated_at = primitive.NewDateTimeFromTime(now)
	projectWA.ID = primitive.NewObjectID()

	_, insertErr := projectCollection.InsertOne(ctx, projectWA.Project)

	if insertErr != nil {
		return http.StatusBadRequest, "User item was not created", "0"
	}

	return http.StatusOK, "Successfully created user", projectWA.ID.String()
}

func (handler *AppHandler) UpdateProject(ctx context.Context, projectCollection *mongo.Collection, projectWA ProjectWithActiveStatus, projectID string) (int, string) {
	objectID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return http.StatusBadRequest, "no project found with the given project_id"
	}
	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{},
	}

	if projectWA.ProjectName != nil {
		update["$set"].(bson.M)["projectname"] = projectWA.ProjectName
	}
	if projectWA.Members != nil {
		update["$set"].(bson.M)["members"] = projectWA.Members
	}

	update["$set"].(bson.M)["updated_at"] = primitive.NewDateTimeFromTime(time.Now())
	result, err := projectCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("error updating project: %v", err)
	}

	if result.MatchedCount == 0 {
		return http.StatusBadRequest, "no project found with the given projectname"
	}

	return http.StatusOK, "project successfully updated"
}

func (handler *AppHandler) getProjectName(projectCollection *mongo.Collection, projectID string) string {
	opts := options.FindOne().SetProjection(bson.M{"projectname": 1, "_id": 0})
	objectID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		handler.Context.Logger.Info(err)
		return ""
	}

	filter := bson.M{"_id": objectID}
	var result struct {
		ProjectName string `bson:"projectname"`
	}

	err = projectCollection.FindOne(handler.Context.Ctx, filter, opts).Decode(&result)
	if err != nil {
		return ""
	}

	return result.ProjectName
}

func (handler *AppHandler) getAllProjectNamesAndIDs(projectCollection *mongo.Collection) ([]models.CombinedProjects, error) {
	opts := options.Find().SetProjection(bson.M{"_id": 1, "projectname": 1})
	var projects []models.CombinedProjects
	cursor, err := projectCollection.Find(handler.Context.Ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	for cursor.Next(handler.Context.Ctx) {
		var result struct {
			ID          primitive.ObjectID `bson:"_id"`
			ProjectName string             `bson:"projectname"`
		}

		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		projects = append(projects, models.CombinedProjects{ProjectID: result.ID.Hex(), ProjectName: result.ProjectName})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}
