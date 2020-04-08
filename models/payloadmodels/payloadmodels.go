package payloadmodels

// SignIn models
type SignIn struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6,max=10"`
}
