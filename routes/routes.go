package routes

import (
	"go_rest/controller"
	"go_rest/utils/auth"

	"github.com/gorilla/mux"
)

// Handlers Set type
func Handlers() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/signup", controller.SignUp).Methods("POST")
	router.HandleFunc("/api/signin", controller.SignIn).Methods("POST")
	router.HandleFunc("/api/resetPasswordLink", controller.ResetPasswordLink).Methods("POST")
	router.HandleFunc("/api/resetPassword", controller.ResetPassword).Methods("PUT")
	subRouter := router.PathPrefix("/api").Subrouter()
	protectedRoute(subRouter)
	return router
}

func protectedRoute(subRoute *mux.Router) {
	subRoute.Use(auth.AuthenticateMiddleware)
	subRoute.HandleFunc("/users", controller.FetchUsers).Methods("GET")
	subRoute.HandleFunc("/changepassword", controller.ChangePassword).Methods("PUT")
}
