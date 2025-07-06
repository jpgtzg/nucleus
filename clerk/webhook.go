package clerk

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"nucleus/types/clerk"
	"os"

	svix "github.com/svix/svix-webhooks/go"
)

type ClerkWebhookEvent = clerk.ClerkWebhookEvent

var wh *svix.Webhook

func init() {
	var err error
	wh, err = svix.NewWebhook(os.Getenv("CLERK_WEBHOOK_SECRET"))
	if err != nil {
		log.Fatalf("Error creating Clerk webhook: %v", err)
	}
}

// HandleWebhook is a handler that receives webhooks from Clerk and processes them
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	header := r.Header
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verify the webhook signature using svix
	err = wh.Verify(payload, header)
	if err != nil {
		log.Printf("Webhook verification failed: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Parse the webhook event
	var event ClerkWebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		log.Printf("Error parsing webhook event: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	go processWebhookEvent(&event)
}

// ProcessWebhook is a function that processes the webhook from Clerk
func processWebhookEvent(event *ClerkWebhookEvent) error {
	log.Printf("[CLERK] Processing webhook event: %s", event.Type)

	switch event.Type {
	case "organization.created":
		return HandleOrganizationCreated(event)
	case "organization.updated":
		return HandleOrganizationUpdated(event)
	case "organization.deleted":
		return HandleOrganizationDeleted(event)
	default:
		log.Printf("Unhandled webhook event type: %s", event.Type)
		return nil
	}
}
