package payloadmodels

import "github.com/dgrijalva/jwt-go"

// SignIn models
type SignIn struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6,max=10"`
}

// ResetPasswordLink reset password request payload struct
type ResetPasswordLink struct {
	Email string `validate:"required,email"`
}

// ResetPassword request paload
type ResetPassword struct {
	Password string `validate:"required,min=6,max=10"`
	Token    string `validate:"required"`
}

// ChangePassword request payload
type ChangePassword struct {
	Password    string `validate:"required,min=6,max=10"`
	NewPassword string `validate:"required,min=6,max=10"`
}

// Token claim models for token
type Token struct {
	ID    string
	Name  string
	Email string
	*jwt.StandardClaims
}
