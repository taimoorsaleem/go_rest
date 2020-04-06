package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang-assignment/models"
	"golang-assignment/utils"
	"golang-assignment/utils/auth"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var trans ut.Translator
var validate *validator.Validate

// SignUp controller request handler
func SignUp(response http.ResponseWriter, request *http.Request) {
	var user models.User

	decoderError := json.NewDecoder(request.Body).Decode(&user)
	if decoderError != nil {
		utils.GetError(decoderError, response)
		return
	}

	validationError := validate.Struct(user)
	if validationError != nil {
		fmt.Println(validationError.(validator.ValidationErrors)[0].Translate(trans))
		utils.GetError(errors.New(validationError.(validator.ValidationErrors)[0].Translate(trans)), response)
		return
	}

	user.PASSWORD, _ = auth.GeneratePassword(user.PASSWORD)
	userCollection := utils.GetCollection(utils.GetUserTable())
	insertedUser, insertError := userCollection.InsertOne(context.TODO(), user)
	if insertError != nil {
		fmt.Println("Error occurred while creating user")
		fmt.Println(insertError)
		utils.GetError(insertError, response)
		return
	}
	user.ID, _ = insertedUser.InsertedID.(primitive.ObjectID)

	token, tokenError := auth.GenerateToken(user)
	if tokenError != nil {
		fmt.Println("Error occurred while creating user")
		fmt.Println(tokenError)
		utils.GetError(tokenError, response)
		return
	}

	var reqResponse = map[string]interface{}{
		"status":  true,
		"message": "User Signup successfully!",
		"id":      user.ID,
		"name":    user.NAME,
		"email":   user.EMAIL,
		"token":   token,
	}
	json.NewEncoder(response).Encode(reqResponse)
}

// SignIn user request handler
func SignIn(response http.ResponseWriter, request *http.Request) {
	var payload models.SignIn
	decoderError := json.NewDecoder(request.Body).Decode(&payload)
	if decoderError != nil {
		utils.GetError(decoderError, response)
		return
	}

	validationError := validate.Struct(payload)
	if validationError != nil {
		fmt.Println(validationError.(validator.ValidationErrors)[0].Translate(trans))
		utils.GetError(errors.New(validationError.(validator.ValidationErrors)[0].Translate(trans)), response)
		return
	}

	dbUser, err := fetchUserByEmail(payload.EMAIL)
	if err != nil {
		fmt.Println("Error occurred while creating user")
		fmt.Println(err)
		utils.GetError(err, response)
		return
	}
	if isMatch, passError := auth.CompareHashAndPassword(dbUser.PASSWORD, payload.PASSWORD); !isMatch && passError != nil {
		fmt.Println("Invalid login credentials")
		fmt.Println(passError)
		utils.GetError(passError, response)
		return
	}

	token, tokenError := auth.GenerateToken(dbUser)
	if tokenError != nil {
		fmt.Println("Error occurred while creating user")
		fmt.Println(tokenError)
		utils.GetError(tokenError, response)
		return
	}

	var reqResponse = map[string]interface{}{
		"status":  true,
		"message": "User Logged in successfully!",
		"id":      dbUser.ID,
		"name":    dbUser.NAME,
		"email":   dbUser.EMAIL,
		"token":   token,
	}
	json.NewEncoder(response).Encode(reqResponse)
}

// ResetPasswordLink forget password for user request handler
func ResetPasswordLink(response http.ResponseWriter, request *http.Request) {
	var payload models.ResetPasswordLink
	decoderError := json.NewDecoder(request.Body).Decode(&payload)
	if decoderError != nil {
		utils.GetError(decoderError, response)
		return
	}
	// validate payload
	validationError := validate.Struct(payload)
	if validationError != nil {
		fmt.Println(validationError.(validator.ValidationErrors)[0].Translate(trans))
		utils.GetError(errors.New(validationError.(validator.ValidationErrors)[0].Translate(trans)), response)
		return
	}
	// fetch DB user
	dbUser, err := fetchUserByEmail(payload.EMAIL)
	if err != nil {
		fmt.Println("Error occurred while fetching user by email")
		fmt.Println(err)
		utils.GetError(err, response)
		return
	}
	// generate token
	token, tokenError := auth.GenerateToken(dbUser)
	if tokenError != nil {
		fmt.Println("Error occurred while generating token")
		fmt.Println(tokenError)
		utils.GetError(tokenError, response)
		return
	}
	// Send email to request user with reset password link
	emailBody := "Reset Password Link: \n http://localhost:8000/?token=" + token
	isEmailSend, emailError := sendEmail(emailBody, dbUser)
	if !isEmailSend && emailError != nil {
		fmt.Println("Error occurred while sending email to user")
		fmt.Println(emailError)
		utils.GetError(emailError, response)
		return
	}
	// set
	resetTokenCollection := utils.GetCollection(utils.GetResetTokenTable())
	// var updatedDocument bson.M
	opt := options.FindOneAndUpdate().SetUpsert(true)
	resetTokenCollection.FindOneAndUpdate(
		context.TODO(),
		bson.D{{"_id", dbUser.ID}},
		bson.D{{"$set", bson.D{{"token", token}}}}, opt) //.Decode(updatedDocument)
	json.NewEncoder(response).Encode(map[string]string{
		"Message": "link created succssfully!",
		"link":    emailBody,
	})
}

// ResetPassword reset user password of provided token
func ResetPassword(response http.ResponseWriter, request *http.Request) {
	var payload models.ResetPassword
	decoderError := json.NewDecoder(request.Body).Decode(&payload)
	if decoderError != nil {
		utils.GetError(decoderError, response)
		return
	}
	// validate payload
	validationError := validate.Struct(payload)
	if validationError != nil {
		fmt.Println(validationError.(validator.ValidationErrors)[0].Translate(trans))
		utils.GetError(errors.New(validationError.(validator.ValidationErrors)[0].Translate(trans)), response)
		return
	}
	// Fetch reset token document by provided token
	var resetPasswordToken models.ResetPasswordToken
	resetTokenCollection := utils.GetCollection(utils.GetResetTokenTable())
	resetPasswordTokenError := resetTokenCollection.FindOne(context.TODO(), bson.M{
		"token": payload.Token,
	}).Decode(&resetPasswordToken)
	if resetPasswordTokenError != nil {
		fmt.Println("Error occurred while checking provided token")
		fmt.Println(resetPasswordTokenError)
		utils.GetError(resetPasswordTokenError, response)
		return
	}
	// Generate password hash
	password, _ := auth.GeneratePassword(payload.Password)
	updatePayload := bson.D{{
		"$set", bson.D{
			{"password", password},
		},
	},
	}
	// Update user password
	findAndUpdate(resetPasswordToken.ID, updatePayload)
	// Delete used token
	resetTokenCollection.DeleteOne(context.TODO(), bson.M{"_id": resetPasswordToken.ID})
	json.NewEncoder(response).Encode(map[string]string{
		"Message": "Password updated successfully!",
	})
}

// ChangePassword request handler
func ChangePassword(response http.ResponseWriter, request *http.Request) {
	var payload models.ChangePassword
	decoderError := json.NewDecoder(request.Body).Decode(&payload)
	if decoderError != nil {
		utils.GetError(decoderError, response)
		return
	}
	// validate payload
	validationError := validate.Struct(payload)
	if validationError != nil {
		fmt.Println(validationError.(validator.ValidationErrors)[0].Translate(trans))
		utils.GetError(errors.New(validationError.(validator.ValidationErrors)[0].Translate(trans)), response)
		return
	}
	// fetch users
	user := utils.GetContextTokenClaims(request.Context())
	dbUser, err := fetchUserByEmail(user.Email)
	if err != nil {
		fmt.Println("Error occurred while fetching user")
		fmt.Println(err)
		utils.GetError(err, response)
		return
	}
	// compare password
	if isMatch, passError := auth.CompareHashAndPassword(dbUser.PASSWORD, payload.Password); !isMatch && passError != nil {
		fmt.Println("Invalid login credentials")
		fmt.Println(passError)
		utils.GetError(passError, response)
		return
	}
	// generate new password
	password, _ := auth.GeneratePassword(payload.NewPassword)
	updatePayload := bson.D{{
		"$set", bson.D{
			{"password", password},
		},
	},
	}
	// Update new password
	findAndUpdate(dbUser.ID, updatePayload)
	json.NewEncoder(response).Encode(map[string]string{
		"Message": "Password updated successfully!",
	})
}

// FetchUsers fetch all user request handler
func FetchUsers(response http.ResponseWriter, request *http.Request) {
	users := fetchUsers()
	json.NewEncoder(response).Encode(users)
}

func sendEmail(body string, user models.User) (bool, error) {
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

// fetchUserByEmail
func fetchUserByEmail(email string) (models.User, error) {
	var user models.User = models.User{}
	userCollection := utils.GetCollection(utils.GetUserTable())
	err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		fmt.Println("Error occurred while fetching user by email ", err)
		return user, err
	}
	return user, nil
}

// findAndUpdate
func findAndUpdate(id primitive.ObjectID, updatePayload bson.D) *mongo.SingleResult {
	userCollection := utils.GetCollection(utils.GetUserTable())
	updatedUser := userCollection.FindOneAndUpdate(context.TODO(), bson.M{"_id": id}, updatePayload)
	return updatedUser
}

// fetchUsers
func fetchUsers() []models.User {
	var users []models.User
	userCollection := utils.GetCollection(utils.GetUserTable())

	// cursor, err := userCollection.Find(context.TODO(), bson.M{})
	// if err != nil {
	// 	return users
	// }
	// defer cursor.Close(context.TODO())
	// for cursor.Next(context.TODO()) {
	// 	var user models.User
	// 	fmt.Println("Y")
	// 	err := cursor.Decode(&user)
	// 	if err != nil {
	// 		return users
	// 	}
	// 	users = append(users, user)
	// }
	// if err := cursor.Err(); err != nil {
	// 	fmt.Println(err)
	// }

	// Load all user
	cursor, err := userCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		fmt.Println(err)
	}
	if cursorError := cursor.All(context.TODO(), &users); cursorError != nil {
		fmt.Println(cursorError)
	}
	return users
}

func init() {
	en := en.New()
	uni := ut.New(en, en)
	var found bool
	trans, found = uni.GetTranslator("en")
	if !found {
		fmt.Print("translator not found")
	}
	validate = validator.New()
	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		fmt.Print(err)
	}

	_ = validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = validate.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

}
