package clerk

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"nucleus/types/clerk"
	"os"

	clerkSDK "github.com/clerk/clerk-sdk-go/v2"
	"github.com/joho/godotenv"
	svix "github.com/svix/svix-webhooks/go"
)

type ClerkWebhookEvent = clerk.ClerkWebhookEvent
type UserData = clerk.UserData

var wh *svix.Webhook

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	clerkAPIKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkAPIKey == "" {
		log.Fatal("CLERK_SECRET_KEY environment variable is required")
	}

	clerkSDK.SetKey(clerkAPIKey)

	wh, err = svix.NewWebhook(os.Getenv("CLERK_WEBHOOK_SECRET"))
	if err != nil {
		log.Fatalf("Error creating Clerk webhook: %v", err)
	}
}

// WebhookHandler is a handler that receives webhooks from Clerk and processes them
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
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

	// Process the webhook event
	if err := processWebhookEvent(event); err != nil {
		log.Printf("Error processing webhook event: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook processed successfully"))
}

// ProcessWebhook is a function that processes the webhook from Clerk
func processWebhookEvent(event ClerkWebhookEvent) error {
	log.Printf("Processing webhook event: %s", event.Type)

	switch event.Type {
	case "user.created":
		return handleUserCreated(event)
	case "user.updated":
		return handleUserUpdated(event)
	case "user.deleted":
		return handleUserDeleted(event)
	case "session.created":
		return handleSessionCreated(event)
	case "session.ended":
		return handleSessionEnded(event)
	case "email.created":
		return handleEmailCreated(event)
	default:
		log.Printf("Unhandled webhook event type: %s", event.Type)
		return nil
	}
}

// handleUserCreated processes user.created webhook events
func handleUserCreated(event ClerkWebhookEvent) error {
	var userData UserData
	if err := json.Unmarshal(event.Data, &userData); err != nil {
		return err
	}

	log.Printf("User created: %s (%s)", userData.ID, userData.GetPrimaryEmail())

	// TODO: Add your business logic here
	// Examples:
	// - Create user record in your database
	// - Send welcome email
	// - Initialize user preferences
	// - Create user profile

	return nil
}

// handleUserUpdated processes user.updated webhook events
func handleUserUpdated(event ClerkWebhookEvent) error {
	var userData UserData
	if err := json.Unmarshal(event.Data, &userData); err != nil {
		return err
	}

	log.Printf("User updated: %s", userData.ID)

	// TODO: Add your business logic here
	// Examples:
	// - Update user record in your database
	// - Sync user data with other services
	// - Update user preferences

	return nil
}

// handleUserDeleted processes user.deleted webhook events
func handleUserDeleted(event ClerkWebhookEvent) error {
	var userData UserData
	if err := json.Unmarshal(event.Data, &userData); err != nil {
		return err
	}

	log.Printf("User deleted: %s", userData.ID)

	// TODO: Add your business logic here
	// Examples:
	// - Mark user as deleted in your database
	// - Clean up user data
	// - Cancel user subscriptions
	// - Archive user data

	return nil
}

// handleSessionCreated processes session.created webhook events
func handleSessionCreated(event ClerkWebhookEvent) error {
	log.Printf("Session created for event: %s", event.Type)

	// TODO: Add your business logic here
	// Examples:
	// - Log user activity
	// - Update last login time
	// - Track session metrics

	return nil
}

// handleSessionEnded processes session.ended webhook events
func handleSessionEnded(event ClerkWebhookEvent) error {
	log.Printf("Session ended for event: %s", event.Type)

	// TODO: Add your business logic here
	// Examples:
	// - Log session duration
	// - Update user activity metrics
	// - Clean up session data

	return nil
}

// handleEmailCreated processes email.created webhook events
func handleEmailCreated(event ClerkWebhookEvent) error {
	log.Printf("Email created for event: %s", event.Type)

	// TODO: Add your business logic here
	// Examples:
	// - Track email delivery
	// - Update email preferences
	// - Log email events

	return nil
}
