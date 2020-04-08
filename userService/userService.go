package userservice

import (
	"context"
	"fmt"

	"golang-assignment/models/entities"
	"golang-assignment/models/payloadmodels"
	"golang-assignment/models/responsemodels"

	"golang-assignment/utils"
	"golang-assignment/utils/auth"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SignUp create user and return response
func SignUp(user entities.User) (*responsemodels.SignupResponse, error) {
	user.PASSWORD, _ = auth.GeneratePassword(user.PASSWORD)
	userCollection := utils.GetCollection(utils.GetUserTable())
	insertedUser, insertError := userCollection.InsertOne(context.TODO(), user)
	if insertError != nil {
		fmt.Println("Error occurred while creating user")
		fmt.Println(insertError)
		return nil, insertError
	}
	user.ID, _ = insertedUser.InsertedID.(primitive.ObjectID)

	token, tokenError := auth.GenerateToken(&user)
	if tokenError != nil {
		fmt.Println("Error occurred while creating user")
		fmt.Println(tokenError)
		return nil, tokenError
	}

	var reqResponse = responsemodels.SignupResponse{
		Status:  true,
		Message: "User Signup successfully!",
		Id:      user.ID.Hex(),
		Name:    user.NAME,
		Email:   user.EMAIL,
		Token:   token,
	}
	return &reqResponse, nil
}

// SignIn sign in user and create token
func SignIn(payload payloadmodels.SignIn) (*responsemodels.SignupResponse, error) {
	user, fetchUserError := fetchUserByEmail(payload.Email)
	if fetchUserError != nil {
		fmt.Println("Error occurred while creating user")
		fmt.Println(fetchUserError)
		return nil, fetchUserError
	}
	if isMatch, passMatchError := auth.CompareHashAndPassword(user.PASSWORD, payload.Password); !isMatch && passMatchError != nil {
		fmt.Println("Invalid login credentials")
		fmt.Println(passMatchError)
		return nil, passMatchError
	}
	token, tokenError := auth.GenerateToken(user)
	if tokenError != nil {
		fmt.Println("Error occurred while creating user")
		fmt.Println(tokenError)
		return nil, tokenError
	}

	var reqResponse = responsemodels.SignupResponse{
		Status:  true,
		Message: "User Logged in successfully!",
		Id:      user.ID.Hex(),
		Name:    user.NAME,
		Email:   user.EMAIL,
		Token:   token,
	}
	return &reqResponse, nil
}

//***************************
// fetchUserByEmail fetch user by email and return user object
func fetchUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	userCollection := utils.GetCollection(utils.GetUserTable())
	err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		fmt.Println("Error occurred while fetching user by email ", err)
		return nil, err
	}
	return &user, nil
}
