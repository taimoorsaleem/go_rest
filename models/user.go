package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User Model Defination
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	NAME     string             `json:"name,omitempty" bson:"name,omitempty"`
	EMAIL    string             `json:"email,omitempty" bson:"email,omitempty"`
	PASSWORD string             `json:"password,omitempty" bson:"password,omitempty"`
}
