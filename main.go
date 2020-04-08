package main

import (
	"golang-assignment/routes"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", routes.Handlers())
	log.Fatal(http.ListenAndServe(":8000", nil))
}
