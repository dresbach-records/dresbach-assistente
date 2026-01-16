package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dresbach/project/handlers"
	"github.com/dresbach/project/session"
	"github.com/dresbach/project/statemachine"
	"github.com/dresbach/project/whatsapp"
)

func main() {
	// Initialize the state machine
	sm := statemachine.NewStateMachine()

	// Register all the handlers
	handlers.RegisterHandlers(sm)

	// Initialize the session store
	store := session.NewStore(30 * time.Minute)

	// Initialize the webhook handler
	webhookHandler := whatsapp.NewWebhookHandler(sm, store)

	// Register the webhook handler
	http.Handle("/webhook", webhookHandler)

	// Start the HTTP server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
