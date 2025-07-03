package clerk

import (
	"encoding/json"
	"log"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	clerkAPIKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkAPIKey == "" {
		log.Fatal("CLERK_SECRET_KEY environment variable is required")
	}

	clerk.SetKey(clerkAPIKey)
}

// handleUserCreated processes user.created webhook events
func HandleUserCreated(event ClerkWebhookEvent) error {
	var userData clerk.User
	if err := json.Unmarshal(event.Data, &userData); err != nil {
		return err
	}

	log.Printf("User created: %s (%s)", userData.ID, userData.EmailAddresses[0].EmailAddress)

	// TODO: Add your business logic here
	// Examples:
	// - Create user record in your database
	// - Send welcome email
	// - Initialize user preferences
	// - Create user profile

	return nil
}

// handleUserUpdated processes user.updated webhook events
func HandleUserUpdated(event ClerkWebhookEvent) error {
	var userData clerk.User
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
func HandleUserDeleted(event ClerkWebhookEvent) error {
	var userData clerk.User
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
func HandleSessionCreated(event ClerkWebhookEvent) error {
	log.Printf("Session created for event: %s", event.Type)

	// TODO: Add your business logic here
	// Examples:
	// - Log user activity
	// - Update last login time
	// - Track session metrics

	return nil
}

// handleSessionEnded processes session.ended webhook events
func HandleSessionEnded(event ClerkWebhookEvent) error {
	log.Printf("Session ended for event: %s", event.Type)

	// TODO: Add your business logic here
	// Examples:
	// - Log session duration
	// - Update user activity metrics
	// - Clean up session data

	return nil
}

// handleEmailCreated processes email.created webhook events
func HandleEmailCreated(event ClerkWebhookEvent) error {
	log.Printf("Email created for event: %s", event.Type)

	// TODO: Add your business logic here
	// Examples:
	// - Track email delivery
	// - Update email preferences
	// - Log email events

	return nil
}
