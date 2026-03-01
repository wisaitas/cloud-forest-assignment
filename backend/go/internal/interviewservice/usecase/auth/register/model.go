package register

// Request for user registration
type Request struct {
	Email           string `json:"email" validate:"required,email" example:"user@example.com"`
	Password        string `json:"password" validate:"required,min=8" example:"password123"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password" example:"password123"`
}
