package responseModels

// SignupResponse signup response struct
type SignupResponse struct {
	Status  bool
	Message string
	Id      string
	Name    string
	Email   string
	Token   string
}
