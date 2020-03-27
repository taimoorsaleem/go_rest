package main

import (
	"context"
	"encoding/json"
	"golang-assignment/routes"
	"helper"
	"log"
	"models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	// router := mux.NewRouter()
	// router.HandleFunc("/api/users", getUsers).Methods("GET")
	// log.Fatal(http.ListenAndServe(":8000", router))
	http.Handle("/", routes.Handlers())
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// we created Book array
	var users []models.User

	//Connection mongoDB with helper class
	collection := helper.ConnectDB()

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var user models.User
		// & character returns the memory address of the following variable.
		err := cur.Decode(&user) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(users) // encode similar to serialize process.
}
