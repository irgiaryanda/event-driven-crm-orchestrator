package models

import "time"

// Event represents an event in the CRM orchestrator
type Event struct {
	ID        string    `json:"id"`
	Payload   string    `json:"payload"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
