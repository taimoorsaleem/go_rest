package utils

import (
	"context"
	"go_rest/models/payloadmodels"
)

// GetContextTokenClaims return user token claims
func GetContextTokenClaims(requestContext context.Context) *payloadmodels.Token {
	return requestContext.Value(GetUserTable()).(*payloadmodels.Token)
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
