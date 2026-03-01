package getservers

// Request for list servers (user_id from JWT Locals)
type Request struct {
	UserID string
}

// ServerItem is one server in the list response
type ServerItem struct {
	ID                       string `json:"id"`
	InfrastructureResourceID string `json:"infrastructure_resource_id"`
	SKU                      string `json:"sku"`
	PowerStatus              string `json:"power_status"`
}

// Response after successful list
type Response struct {
	Servers []ServerItem `json:"servers"`
}
