package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-assignment/models"
	"golang-assignment/utils"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// GenerateToken generate token with provided user claims
func GenerateToken(user models.User) (string, error) {
	expiresAt := time.Now().Add(time.Minute * 2880).Unix()
	tokenClaims := &models.Token{
		ID:    user.ID.String(),
		Name:  user.NAME,
		Email: user.EMAIL,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}
	encodedToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenClaims)
	token, err := encodedToken.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println("Error occurred while hashing password", err)
		return "", err
	}
	return token, nil
}

// GenerateRefreshToken generate token with provided user claims
func GenerateRefreshToken(user models.User) (string, error) {
	expiresAt := time.Now().Add(time.Minute * 10000).Unix()
	tokenClaims := &models.Token{
		ID: user.ID.String(),
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}
	encodedToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenClaims)
	token, err := encodedToken.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println("Error occurred while hashing password", err)
		return "", err
	}
	return token, nil
}

// GeneratePassword generate password of provided string
func GeneratePassword(pass string) (string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error occurred while hashing password", err)
		return "", err
	}
	return string(hashPass), nil
}

// CompareHashAndPassword compare hash and password
func CompareHashAndPassword(dbPassword string, userPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(userPassword))
	if err != nil || err == bcrypt.ErrMismatchedHashAndPassword {
		return false, err
	}
	return true, nil
}

// AuthenticateMiddleware use to authenicate user
func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var token string = request.Header.Get("x-access-token")
		token = strings.TrimSpace(token)
		if token == "" {
			response.WriteHeader(http.StatusForbidden)
			json.NewEncoder(response).Encode(map[string]interface{}{"Message": "Missing auth token"})
			return
		}
		tokenClaims := &models.Token{}
		_, err := jwt.ParseWithClaims(token, tokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			response.WriteHeader(http.StatusForbidden)
			json.NewEncoder(response).Encode(map[string]interface{}{"Message": err.Error()})
			return
		}
		userContext := context.WithValue(request.Context(), utils.GetUserTable(), tokenClaims)
		next.ServeHTTP(response, request.WithContext(userContext))
	})
}
