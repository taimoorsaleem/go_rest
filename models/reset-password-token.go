package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ResetPasswordToken Model Defination
type ResetPasswordToken struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Token string             `json:"token,omitempty" bson:"token,omitempty"`
}

// ResetPasswordTokenPayload s
type ResetPasswordTokenPayload struct {
	Password string
	Token    string
}
