package models

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Username      *string            `json:"username" validate:"required,min=2,max=100"`
	Password      *string            `json:"password" validate:"required,min=6"`
	Email         *string            `json:"email" bson:"email" validate:"email,required"`
	Role          *string            `json:"role" bson:"role"`
	Token         *string            `json:"token" bson:"token"`
	BaseGroup     *string            `json:"base_group" bson:"base_group"`
	Active        *bool              `json:"active" bson:"active"`
	Refresh_token *string            `json:"refresh_token" bson:"refresh_token"`
	Created_at    primitive.DateTime `json:"created_at" bson:"created_at"`
	Updated_at    primitive.DateTime `json:"updated_at" bson:"updated_at"`
	User_id       string             `json:"user_id" bson:"user_id"`
}

type Group struct {
	ID         primitive.ObjectID `bson:"_id"`
	GroupName  *string            `json:"groupname" bson:"groupname" validate:"required,min=2,max=100"`
	Members    []string           `json:"members" bson:"members"`
	Created_at primitive.DateTime `json:"created_at" bson:"created_at"`
	Updated_at primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

type UserList struct {
	ID         primitive.ObjectID `bson:"_id"`
	Username   *string            `json:"username" validate:"required,min=2,max=100"`
	Email      *string            `json:"email" validate:"email,required"`
	Created_at primitive.DateTime `json:"created_at"`
	Updated_at primitive.DateTime `json:"updated_at"`
	User_id    string             `json:"user_id"`
	Groups     []string           `json:"groups"`
}

type SignedDetails struct {
	Email      string
	Username   string
	UserId     string
	Groups     []string
	Role       string
	BaseGroup  *string
	AdminGroup bool
	jwt.RegisteredClaims
}

type InnerPermission struct {
	Added   []string `json:"added,omitempty"`
	Removed []string `json:"removed,omitempty"`
}

type Permission struct {
	Listeners  *InnerPermission `json:"listeners,omitempty"`
	Routes     *InnerPermission `json:"routes,omitempty"`
	Clusters   *InnerPermission `json:"clusters,omitempty"`
	Endpoints  *InnerPermission `json:"endpoints,omitempty"`
	Secrets    *InnerPermission `json:"secrets,omitempty"`
	Extensions *InnerPermission `json:"extensions,omitempty"`
	Bootstrap  *InnerPermission `json:"bootstrap,omitempty"`
}
