package db

import (
	"database/sql"
	"time"

	"github.com/irgiaryanda/event-driven-crm-orchestrator/internal/models"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

// Initialize sets up the database connection and creates tables
func Initialize(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// Create events table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		payload TEXT,
		status TEXT,
		created_at DATETIME
	);
	`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		return err
	}

	return nil
}

// InsertEvent inserts a new event, returns true if new, false if already exists (idempotent)
func InsertEvent(event *models.Event) (bool, error) {
	// Check if event already exists
	var exists int
	err := DB.QueryRow("SELECT COUNT(*) FROM events WHERE id = ?", event.ID).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists > 0 {
		return false, nil // Already exists, idempotent response
	}

	// Insert new event
	_, err = DB.Exec(
		"INSERT INTO events (id, payload, status, created_at) VALUES (?, ?, ?, ?)",
		event.ID,
		event.Payload,
		event.Status,
		event.CreatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

// UpdateEventStatus updates the status of an existing event
func UpdateEventStatus(id string, status string) error {
	_, err := DB.Exec("UPDATE events SET status = ? WHERE id = ?", status, id)
	return err
}

// GetAllEvents retrieves all events ordered by created_at DESC
func GetAllEvents() ([]models.Event, error) {
	rows, err := DB.Query("SELECT id, payload, status, created_at FROM events ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Payload, &event.Status, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
