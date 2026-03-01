package login

// Request for login
type Request struct {
	Email    string `json:"email" validate:"required,email" example:"john.smith@gmail.com"`
	Password string `json:"password" validate:"required,min=8" example:"not-so-secure-password"`
}

// Response after successful login
type Response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
