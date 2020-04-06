package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ResetPasswordToken Model Defination
type ResetPasswordToken struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Token string             `json:"token,omitempty" bson:"token,omitempty"`
}

// ResetPassword
type ResetPassword struct {
	Password string `validate:"required,min=6,max=10"`
	Token    string `validate:"required"`
}
