package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang-assignment/models/entities"
	"golang-assignment/models/payloadmodels"
	"golang-assignment/userservice"
	"golang-assignment/utils"
	"net/http"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var trans ut.Translator
var validate *validator.Validate

// SignUp controller request handler
func SignUp(response http.ResponseWriter, request *http.Request) {
	// Decode payload from request
	var user entities.User
	decoderError := json.NewDecoder(request.Body).Decode(&user)
	if decoderError != nil {
		utils.GetError(decoderError, response)
		return
	}
	// validate payload
	validationError := validate.Struct(user)
	if validationError != nil {
		fmt.Println(validationError.(validator.ValidationErrors)[0].Translate(trans))
		utils.GetError(errors.New(validationError.(validator.ValidationErrors)[0].Translate(trans)), response)
		return
	}
	// create user and generate token for created user
	reqResponse, signError := userservice.SignUp(user)
	if signError != nil {
		fmt.Println(signError)
		utils.GetError(signError, response)
		return
	}
	json.NewEncoder(response).Encode(reqResponse)
}

// SignIn user request handler
func SignIn(response http.ResponseWriter, request *http.Request) {
	// Decode payload from request
	var payload payloadmodels.SignIn
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
	// sign in user and handle error if occured
	reqResponse, signInError := userservice.SignIn(payload)
	if signInError != nil {
		fmt.Println("Error occurred while signin user")
		fmt.Println(signInError)
		utils.GetError(signInError, response)
		return
	}
	json.NewEncoder(response).Encode(reqResponse)
}

// ResetPasswordLink forget password for user request handler
func ResetPasswordLink(response http.ResponseWriter, request *http.Request) {
	// decode payload from request
	var payload payloadmodels.ResetPasswordLink
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
	// send reset password link on email address and save token in db for validation
	_, sendEmailError := userservice.ResetPasswordLink(payload)
	if sendEmailError != nil {
		fmt.Println("Error occurred while send reset password email to user")
		fmt.Println(sendEmailError)
		utils.GetError(sendEmailError, response)
		return
	}
	json.NewEncoder(response).Encode(map[string]string{
		"Message": "link created succssfully!",
	})
}

// ResetPassword reset user password of provided token
func ResetPassword(response http.ResponseWriter, request *http.Request) {
	// Decode request payload
	var payload payloadmodels.ResetPassword
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
	// Save new password in db
	_, resetPassError := userservice.ResetPassword(payload)
	if resetPassError != nil {
		fmt.Println("Error occurred while checking provided token")
		fmt.Println(resetPassError)
		utils.GetError(resetPassError, response)
		return
	}
	json.NewEncoder(response).Encode(map[string]string{
		"Message": "Password updated successfully!",
	})
}

// // ChangePassword request handler
// func ChangePassword(response http.ResponseWriter, request *http.Request) {
// 	var payload models.ChangePassword
// 	decoderError := json.NewDecoder(request.Body).Decode(&payload)
// 	if decoderError != nil {
// 		utils.GetError(decoderError, response)
// 		return
// 	}
// 	// validate payload
// 	validationError := validate.Struct(payload)
// 	if validationError != nil {
// 		fmt.Println(validationError.(validator.ValidationErrors)[0].Translate(trans))
// 		utils.GetError(errors.New(validationError.(validator.ValidationErrors)[0].Translate(trans)), response)
// 		return
// 	}
// 	// fetch users
// 	user := utils.GetContextTokenClaims(request.Context())
// 	dbUser, err := fetchUserByEmail(user.Email)
// 	if err != nil {
// 		fmt.Println("Error occurred while fetching user")
// 		fmt.Println(err)
// 		utils.GetError(err, response)
// 		return
// 	}
// 	// compare password
// 	if isMatch, passError := auth.CompareHashAndPassword(dbUser.PASSWORD, payload.Password); !isMatch && passError != nil {
// 		fmt.Println("Invalid login credentials")
// 		fmt.Println(passError)
// 		utils.GetError(passError, response)
// 		return
// 	}
// 	// generate new password
// 	password, _ := auth.GeneratePassword(payload.NewPassword)
// 	updatePayload := bson.D{{
// 		"$set", bson.D{
// 			{"password", password},
// 		},
// 	},
// 	}
// 	// Update new password
// 	findAndUpdate(dbUser.ID, updatePayload)
// 	json.NewEncoder(response).Encode(map[string]string{
// 		"Message": "Password updated successfully!",
// 	})
// }

// // FetchUsers fetch all user request handler
// func FetchUsers(response http.ResponseWriter, request *http.Request) {
// 	users := fetchUsers()
// 	json.NewEncoder(response).Encode(users)
// }

// func sendEmail(body string, user entities.User) (bool, error) {
// 	to := user.EMAIL
// 	pass := os.Getenv("password")
// 	from := os.Getenv("from")
// 	msg := "From: " + from + "\n" +
// 		"To: " + to + "\n" +
// 		"Subject: Reset Password\n\n" +
// 		body

// 	err := smtp.SendMail("smtp.gmail.com:25",
// 		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
// 		from, []string{to}, []byte(msg))

// 	if err != nil {
// 		log.Printf("smtp error: %s", err)
// 		return false, err
// 	}
// 	return true, nil
// }

// // fetchUserByEmail
// func fetchUserByEmail(email string) (entities.User, error) {
// 	var user entities.User = entities.User{}
// 	userCollection := utils.GetCollection(utils.GetUserTable())
// 	err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
// 	if err != nil {
// 		fmt.Println("Error occurred while fetching user by email ", err)
// 		return user, err
// 	}
// 	return user, nil
// }

// // findAndUpdate
// func findAndUpdate(id primitive.ObjectID, updatePayload bson.D) *mongo.SingleResult {
// 	userCollection := utils.GetCollection(utils.GetUserTable())
// 	updatedUser := userCollection.FindOneAndUpdate(context.TODO(), bson.M{"_id": id}, updatePayload)
// 	return updatedUser
// }

// // fetchUsers
// func fetchUsers() []models.User {
// 	var users []models.User
// 	userCollection := utils.GetCollection(utils.GetUserTable())

// 	// cursor, err := userCollection.Find(context.TODO(), bson.M{})
// 	// if err != nil {
// 	// 	return users
// 	// }
// 	// defer cursor.Close(context.TODO())
// 	// for cursor.Next(context.TODO()) {
// 	// 	var user models.User
// 	// 	fmt.Println("Y")
// 	// 	err := cursor.Decode(&user)
// 	// 	if err != nil {
// 	// 		return users
// 	// 	}
// 	// 	users = append(users, user)
// 	// }
// 	// if err := cursor.Err(); err != nil {
// 	// 	fmt.Println(err)
// 	// }

// 	// Load all user
// 	cursor, err := userCollection.Find(context.TODO(), bson.M{})

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	if cursorError := cursor.All(context.TODO(), &users); cursorError != nil {
// 		fmt.Println(cursorError)
// 	}
// 	return users
// }

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
