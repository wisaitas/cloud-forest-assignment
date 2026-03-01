package entity

import (
	"time"
)

type ActivityLog struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type,omitempty"`
	ResourceID   string    `json:"resource_id,omitempty"`
	Details      string    `json:"details,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
