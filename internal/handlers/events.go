package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/irgiaryanda/event-driven-crm-orchestrator/internal/db"
)

// GetEvents returns all events as JSON
func GetEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	events, err := db.GetAllEvents()
	if err != nil {
		http.Error(w, "Failed to retrieve events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
