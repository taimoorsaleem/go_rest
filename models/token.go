package models

import jwt "github.com/dgrijalva/jwt-go"

// Token claim models for token
type Token struct {
	ID    string
	Name  string
	Email string
	*jwt.StandardClaims
}

// ChangePassword for pass
type ChangePassword struct {
	Password    string
	NewPassword string
}
