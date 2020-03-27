package utils

import (
	"context"
	"golang-assignment/models"
)

// GetContextTokenClaims return user token claims
func GetContextTokenClaims(requestContext context.Context) *models.Token {
	return requestContext.Value(GetUserTable()).(*models.Token)
}

// const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
// // RandStringBytes generate random string for password
// func RandStringBytes(n int) string {
// 	b := make([]byte, n)
// 	for i := range b {
// 		b[i] = letterBytes[rand.Intn(len(letterBytes))]
// 	}
// 	return string(b)
// }
