package clerk

import (
	"encoding/json"
	"log"
	"nucleus/supabase"
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

func HandleOrganizationCreated(event *ClerkWebhookEvent) error {
	var organization clerk.Organization
	err := json.Unmarshal(event.Data, &organization)
	if err != nil {
		return err
	}

	customer, err := customer.New(&stripe.CustomerParams{
		Name: stripe.String(organization.Name),
	})
	if err != nil {
		return err
	}

	err = supabase.CreateOrganization(organization.ID, customer.ID)
	if err != nil {
		return err
	}

	return nil
}

func HandleOrganizationUpdated(event *ClerkWebhookEvent) error {
	return nil
}

func HandleOrganizationDeleted(event *ClerkWebhookEvent) error {
	var eventData map[string]interface{}
	err := json.Unmarshal(event.Data, &eventData)
	if err != nil {
		return err
	}

	organizationId, ok := eventData["id"].(string)
	if !ok {
		return err
	}

	organization, err := supabase.GetOrganizationByClerkID(organizationId)
	if err != nil {
		supabase.DebugOrganizations()
		return err
	}

	_, err = customer.Del(organization.StripeCustomerID, nil)
	if err != nil {
		return err
	}

	return nil
}
