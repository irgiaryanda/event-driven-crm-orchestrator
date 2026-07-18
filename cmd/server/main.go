package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/irgiaryanda/event-driven-crm-orchestrator/internal/db"
	"github.com/irgiaryanda/event-driven-crm-orchestrator/internal/handlers"
)

func main() {
	// Initialize database
	dbPath := "crm_orchestrator.db"
	if err := db.Initialize(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	// Setup HTTP routes
	http.HandleFunc("/api/webhook", handlers.HandleWebhook)

	// Start server
	server := &http.Server{
		Addr: ":8080",
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		server.Close()
	}()

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}

	log.Println("Server stopped")
}
