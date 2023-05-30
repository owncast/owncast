package models

import "time"

// IPAddress is a simple representation of an IP address.
type IPAddress struct {
	CreatedAt time.Time `json:"createdAt"`
	IPAddress string    `json:"ipAddress"`
	Notes     string    `json:"notes"`
}
