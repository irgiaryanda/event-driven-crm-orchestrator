package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/irgiaryanda/event-driven-crm-orchestrator/internal/db"
	"github.com/irgiaryanda/event-driven-crm-orchestrator/internal/models"
	"github.com/irgiaryanda/event-driven-crm-orchestrator/internal/services"
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
		Status:    "PENDING",
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

	// Process asynchronously
	go processEvent(eventID, string(payloadBytes))
}

func processEvent(eventID string, payload string) {
	log.Printf("Processing event: %s", eventID)

	// Call LLM to categorize
	category, err := services.CategorizePayload(payload)
	if err != nil {
		log.Printf("LLM categorization failed for %s: %v", eventID, err)
		db.UpdateEventStatus(eventID, "FAILED")
		return
	}

	// Send Telegram notification
	message := fmt.Sprintf("🔔 New Webhook!\nCategory: %s\nID: %s", category, eventID)
	if err := services.SendNotification(message); err != nil {
		log.Printf("Telegram notification failed for %s: %v", eventID, err)
		db.UpdateEventStatus(eventID, "FAILED")
		return
	}

	// Update status to PROCESSED
	if err := db.UpdateEventStatus(eventID, "PROCESSED"); err != nil {
		log.Printf("Failed to update status for %s: %v", eventID, err)
		return
	}

	log.Printf("Event %s processed successfully with category: %s", eventID, category)
}
