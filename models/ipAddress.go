package models

import "time"

// IPAddress is a simple representation of an IP address.
type IPAddress struct {
	IPAddress string    `json:"ipAddress"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
}
