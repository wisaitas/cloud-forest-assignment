package power

// Request for power action (server-id from path, action from body)
type Request struct {
	Action string `json:"action" validate:"required,oneof=on off"`
}

// Response after successful power action
type Response struct {
	Success bool   `json:"success"`
	State   string `json:"state"`
}
