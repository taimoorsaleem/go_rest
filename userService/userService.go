package userService

import (
	"context"
	"fmt"
	"golang-assignment/models/entities"
	"golang-assignment/models/responseModels"
	"golang-assignment/utils"
	"golang-assignment/utils/auth"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SignUp
func SignUp(user entities.User) (*responseModels.SignupResponse, error) {
	user.PASSWORD, _ = auth.GeneratePassword(user.PASSWORD)
	userCollection := utils.GetCollection(utils.GetUserTable())
	insertedUser, insertError := userCollection.InsertOne(context.TODO(), user)
	if insertError != nil {
		fmt.Println("Error occurred while creating user")
		fmt.Println(insertError)
		return nil, insertError
	}
	user.ID, _ = insertedUser.InsertedID.(primitive.ObjectID)

	token, tokenError := auth.GenerateToken(user)
	if tokenError != nil {
		fmt.Println("Error occurred while creating user")
		fmt.Println(tokenError)
		return nil, tokenError
	}

	var reqResponse = responseModels.SignupResponse{
		Status:  true,
		Message: "User Signup successfully!",
		Id:      user.ID.Hex(),
		Name:    user.NAME,
		Email:   user.EMAIL,
		Token:   token,
	}
	return &reqResponse, nil
}
