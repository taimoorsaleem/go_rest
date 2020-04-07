package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// enviorment variables
var connectionString = ""
var databaseName = ""
var userTabel string = "user"
var resetPasswordTokenTable string = "reset-token"
var authenticationTokenTable string = "authentication-token"

// connectDB Get database connection client
func connectDB() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		fmt.Print("error occurred while connecting database", err)
	}
	fmt.Println("Connected to MongoDB!")
	return client
}

// GetUserTable Get table name for users
func GetUserTable() string {
	return userTabel
}

// GetResetTokenTable Get table name for users reset token
func GetResetTokenTable() string {
	return resetPasswordTokenTable
}

// GetAuthenticationTokenTable Get authentication token
func GetAuthenticationTokenTable() string {
	return authenticationTokenTable
}

func getDatabaseName() string {
	return databaseName
}

// GetCollection get table/collection by provided name
func GetCollection(collectionName string) *mongo.Collection {
	client := connectDB()
	return client.Database(getDatabaseName()).Collection(collectionName)
}

// init Load enviorment variable
func init() {
	err := godotenv.Load()

	if err != nil {
		fmt.Print("error occurred while loading enviorment varaibales", err)
	}

	databaseName = os.Getenv("databaseName")
	connectionString = os.Getenv("CONNECTIONSTRING")
}
