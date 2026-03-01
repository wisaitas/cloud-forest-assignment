package provision

// Request for provision server
type Request struct {
	SKU string `json:"sku" validate:"required" example:"C1-R1GB-D20GB"`
}

// Response after successful provision
type Response struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
