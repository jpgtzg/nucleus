package clerk

import (
	"encoding/json"
	"log"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
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

func HandleOrganizationCreated(event ClerkWebhookEvent) error {
	var organization clerk.Organization
	err := json.Unmarshal(event.Data, &organization)
	if err != nil {
		return err
	}

	// Create the single stripe customer for the organization,
	// with the organization id in the metadata (required for bidirectional sync)
	customer, err := customer.New(&stripe.CustomerParams{
		Name: stripe.String(organization.Name),
		Metadata: map[string]string{
			"clerk_organization_id": organization.ID,
		},
	})
	if err != nil {
		return err
	}

	// Update the organization metadata with the stripe customer id
	// This is required for bidirectional sync
	metadata := map[string]interface{}{
		"stripe_customer_id": customer.ID,
	}
	if err := UpdateOrganizationPublicMetadata(organization.ID, metadata); err != nil {
		return err
	}

	return nil
}

func HandleOrganizationUpdated(event ClerkWebhookEvent) error {
	log.Printf("Organization updated: %+v", event)
	return nil
}

func HandleOrganizationDeleted(event ClerkWebhookEvent) error {
	log.Printf("Organization deleted: %+v", event)
	return nil
}
