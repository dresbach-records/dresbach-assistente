package whatsapp

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dresbach/project/session"
	"github.com/dresbach/project/statemachine"
)

// WebhookHandler handles incoming WhatsApp messages.
type WebhookHandler struct {
	sm    *statemachine.StateMachine
	store *session.Store
}

// NewWebhookHandler creates a new WebhookHandler.
func NewWebhookHandler(sm *statemachine.StateMachine, store *session.Store) *WebhookHandler {
	return &WebhookHandler{
		sm:    sm,
		store: store,
	}
}

// ServeHTTP handles the HTTP request.
func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Placeholder for WhatsApp webhook verification

	var payload struct {
		// WhatsApp message payload structure
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Error decoding payload: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Placeholder for extracting user ID and message from payload
	userID := "user123"
	message := "1"

	session := h.store.Get(userID)

	response, nextState := h.sm.Handle(session.State, userID, message)

	session.State = nextState
	h.store.Set(userID, session)

	// Placeholder for sending response to WhatsApp
	log.Printf("User: %s, Message: %s, Response: %s, NextState: %s", userID, message, response, nextState)

	w.WriteHeader(http.StatusOK)
}
