package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User Model Defination
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"omitempty,uuid"`
	NAME     string             `json:"name,omitempty" bson:"name,omitempty" validate:"required"`
	EMAIL    string             `json:"email,omitempty" bson:"email,omitempty" validate:"required,email"`
	PASSWORD string             `json:"password,omitempty" bson:"password,omitempty" validate:"required,min=6,max=10"`
}

// SignInPayload
type SignInPayload struct {
	EMAIL    string `validate:"required,email"`
	PASSWORD string `validate:"required,min=6,max=10"`
}

// ResetPasswordLink
type ResetPasswordLink struct {
	EMAIL string `validate:"required,email"`
}
