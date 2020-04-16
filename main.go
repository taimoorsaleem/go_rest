package main

import (
	"go_rest/routes"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", routes.Handlers())
	log.Fatal(http.ListenAndServe(":8000", nil))
}
