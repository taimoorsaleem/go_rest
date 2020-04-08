package userservice

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"golang-assignment/models/entities"
	"golang-assignment/models/payloadmodels"
	"golang-assignment/models/responsemodels"

	"golang-assignment/utils"
	"golang-assignment/utils/auth"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// ResetPasswordLink validate user and send email on user account along with reset password link
func ResetPasswordLink(payload payloadmodels.ResetPasswordLink) (bool, error) {
	// fetch DB user
	user, fetchUserError := fetchUserByEmail(payload.Email)
	if fetchUserError != nil {
		fmt.Println("Error occurred while fetching user by email")
		fmt.Println(fetchUserError)
		return false, fetchUserError
	}
	// generate token for user link
	token, tokenError := auth.GenerateToken(user)
	if tokenError != nil {
		fmt.Println("Error occurred while generating token")
		fmt.Println(tokenError)
		return false, tokenError
	}
	// Send email to request user with reset password link
	emailBody := "Reset Password Link: \n http://localhost:8000/?token=" + token
	isEmailSend, emailError := sendEmail(emailBody, user)
	if !isEmailSend && emailError != nil {
		fmt.Println("Error occurred while sending email to user")
		fmt.Println(emailError)
		return false, emailError
	}
	// Upsert reset password link in DB
	resetTokenCollection := utils.GetCollection(utils.GetResetTokenTable())
	opt := options.FindOneAndUpdate().SetUpsert(true)
	resetTokenCollection.FindOneAndUpdate(
		context.TODO(),
		bson.D{{"_id", user.ID}},
		bson.D{{"$set", bson.D{{"token", token}}}}, opt)
	return true, nil
}

// ResetPassword validate token and save new password in db
func ResetPassword(payload payloadmodels.ResetPassword) (bool, error) {
	// Fetch reset token document by provided token
	var resetPasswordToken entities.ResetPasswordToken
	resetTokenCollection := utils.GetCollection(utils.GetResetTokenTable())
	resetPasswordTokenError := resetTokenCollection.FindOne(context.TODO(), bson.M{
		"token": payload.Token,
	}).Decode(&resetPasswordToken)
	if resetPasswordTokenError != nil {
		fmt.Println("Error occurred while checking provided token")
		fmt.Println(resetPasswordTokenError)
		return false, resetPasswordTokenError
	}
	// Generate new password
	password, _ := auth.GeneratePassword(payload.Password)
	// set new password for user in DB
	updatePayload := bson.D{{
		"$set", bson.D{
			{"password", password},
		},
	},
	}
	findAndUpdate(resetPasswordToken.ID, updatePayload)
	// Delete used token from db
	resetTokenCollection.DeleteOne(context.TODO(), bson.M{"_id": resetPasswordToken.ID})
	return true, nil
}

// ChangePassword change user password
func ChangePassword(payload payloadmodels.ChangePassword, user *payloadmodels.Token) (bool, error) {
	dbUser, fetchUserError := fetchUserByEmail(user.Email)
	if fetchUserError != nil {
		fmt.Println("Error occurred while fetching user")
		fmt.Println(fetchUserError)
		return false, fetchUserError
	}
	// compare database and payload password
	if isMatch, passError := auth.CompareHashAndPassword(dbUser.PASSWORD, payload.Password); !isMatch && passError != nil {
		fmt.Println("Invalid login credentials")
		fmt.Println(passError)
		return false, passError
	}
	// generate new password
	password, _ := auth.GeneratePassword(payload.NewPassword)
	// update user password
	updatePayload := bson.D{{
		"$set", bson.D{
			{"password", password},
		},
	},
	}
	findAndUpdate(dbUser.ID, updatePayload)
	return true, nil
}

// FetchUsers fetch all usersdata
func FetchUsers() (*[]entities.User, error){
	var users []entities.User
	// Load all user
	userCollection := utils.GetCollection(utils.GetUserTable())
	cursor, err := userCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
		return nil, err
	}

	if cursorError := cursor.All(context.TODO(), &users); cursorError != nil {
		fmt.Println(cursorError)
		fmt.Println(cursorError)
		return nil, cursorError
	}
	return &users, nil
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

// sendEmail send reset password link
func sendEmail(body string, user *entities.User) (bool, error) {
	to := user.EMAIL
	pass := os.Getenv("password")
	from := os.Getenv("from")
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Reset Password\n\n" +
		body
	err := smtp.SendMail("smtp.gmail.com:25",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return false, err
	}
	return true, nil
}

// findAndUpdate and update user
func findAndUpdate(id primitive.ObjectID, updatePayload bson.D) *mongo.SingleResult {
	userCollection := utils.GetCollection(utils.GetUserTable())
	updatedUser := userCollection.FindOneAndUpdate(context.TODO(), bson.M{"_id": id}, updatePayload)
	return updatedUser
}
