package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/irgiaryanda/event-driven-crm-orchestrator/internal/db"
	"github.com/irgiaryanda/event-driven-crm-orchestrator/internal/models"
)

// WebhookPayload represents the incoming webhook payload
type WebhookPayload struct {
	EventID string          `json:"event_id"`
	Data    json.RawMessage `json:"data"`
}

// HandleWebhook processes incoming webhook events
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Generate new UUID if event_id is not provided
	eventID := payload.EventID
	if eventID == "" {
		eventID = uuid.New().String()
	}

	// Marshal the data back to string for storage
	payloadBytes, _ := json.Marshal(payload.Data)

	event := &models.Event{
		ID:        eventID,
		Payload:   string(payloadBytes),
		Status:    "received",
		CreatedAt: time.Now(),
	}

	isNew, err := db.InsertEvent(event)
	if err != nil {
		http.Error(w, "Failed to store event", http.StatusInternalServerError)
		return
	}

	if isNew {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"status":   "created",
			"event_id": eventID,
		})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":   "already_exists",
			"event_id": eventID,
		})
	}
}
