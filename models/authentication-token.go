package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthenticationToken Model Defination
type AuthenticationToken struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"omitempty,uuid"`
	USERID       primitive.ObjectID `json:"userid,omitempty" bson:"userid,omitempty" validate:"required"`
	ACCESSTOKEN  string             `json:"access_token,omitempty" bson:"access_token,omitempty" validate:"required"`
	REFRESHTOKEN string             `json:"refresh_token,omitempty" bson:"refresh_token,omitempty" validate:"required"`
}
