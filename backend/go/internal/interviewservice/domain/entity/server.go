package entity

import (
	"time"
)

type Server struct {
	ID                       string    `json:"id"`
	UserID                   string    `json:"user_id"`
	InfrastructureResourceID string    `json:"infrastructure_resource_id"`
	SKU                      string    `json:"sku"`
	PowerStatus              string    `json:"power_status"` // "on" | "off"
	CreatedAt                time.Time `json:"created_at"`
}
