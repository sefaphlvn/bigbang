package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Username      *string            `json:"username" validate:"required,min=2,max=100"`
	Password      *string            `json:"password" validate:"required,min=6"`
	Email         *string            `json:"email" validate:"email,required"`
	Token         *string            `json:"token"`
	Refresh_token *string            `json:"refresh_token"`
	Created_at    primitive.DateTime `json:"created_at"`
	Updated_at    primitive.DateTime `json:"updated_at"`
	User_id       string             `json:"user_id"`
	Groups        []string           `json:"groups"`
}
