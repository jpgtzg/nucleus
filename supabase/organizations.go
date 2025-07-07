package supabase

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	supabaseTypes "nucleus/types/supabase"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

type Organization = supabaseTypes.Organization

var client *supabase.Client

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	var err error
	client, err = supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllOrganizations() ([]Organization, error) {
	data, _, err := client.From("organizations").Select("*", "exact", false).Execute()
	if err != nil {
		return nil, err
	}

	var organizations []Organization
	if err := json.Unmarshal(data, &organizations); err != nil {
		return nil, err
	}

	return organizations, nil
}

// DebugOrganizations logs all organizations in the database
func DebugOrganizations() {
	organizations, err := GetAllOrganizations()
	if err != nil {
		log.Printf("Failed to get organizations: %v", err)
		return
	}

	log.Printf("Found %d organizations in database:", len(organizations))
	for i, org := range organizations {
		log.Printf("  %d: ID=%d, ClerkID=%s, StripeCustomerID=%s", i+1, org.ID, org.ClerkID, org.StripeCustomerID)
	}
}

func GetOrganizationByClerkID(clerkID string) (Organization, error) {
	data, _, err := client.From("organizations").Select("*", "exact", false).Eq("clerk_organization_id", clerkID).Execute()
	if err != nil {
		return Organization{}, err
	}

	var organizations []Organization
	if err := json.Unmarshal(data, &organizations); err != nil {
		return Organization{}, err
	}

	if len(organizations) == 0 {
		return Organization{}, fmt.Errorf("organization not found with clerk ID: %s", clerkID)
	}

	return organizations[0], nil
}

func GetOrganizationByStripeCustomerID(stripeCustomerID string) (Organization, error) {
	data, _, err := client.From("organizations").Select("*", "exact", false).Eq("stripe_customer_id", stripeCustomerID).Execute()
	if err != nil {
		return Organization{}, err
	}

	var organizations []Organization
	if err := json.Unmarshal(data, &organizations); err != nil {
		return Organization{}, err
	}

	if len(organizations) == 0 {
		return Organization{}, fmt.Errorf("organization not found with stripe customer ID: %s", stripeCustomerID)
	}

	return organizations[0], nil
}

func CreateOrganization(clerkID string, stripeCustomerID string) error {
	organizationData := map[string]interface{}{
		"clerk_organization_id": clerkID,
		"stripe_customer_id":    stripeCustomerID,
	}

	_, _, err := client.From("organizations").Insert(organizationData, false, "", "", "").Execute()
	if err != nil {
		return err
	}

	return nil
}

func UpdateOrganizationStripeCustomerID(clerkID string, stripeCustomerID string) error {
	_, _, err := client.From("organizations").
		Update(map[string]interface{}{
			"stripe_customer_id": stripeCustomerID,
		}, "", "").
		Eq("clerk_organization_id", clerkID).
		Execute()

	if err != nil {
		return err
	}

	return nil
}

func UpdateOrganizationClerkID(stripeCustomerID string, clerkID string) error {
	_, _, err := client.From("organizations").
		Update(map[string]interface{}{
			"clerk_organization_id": clerkID,
		}, "", "").
		Eq("stripe_customer_id", stripeCustomerID).
		Execute()

	if err != nil {
		return err
	}

	return nil
}

func DeleteOrganizationByClerkID(clerkID string) error {
	_, _, err := client.From("organizations").
		Delete("", "").
		Eq("clerk_organization_id", clerkID).
		Execute()

	if err != nil {
		return err
	}

	return nil
}

func DeleteOrganizationByStripeCustomerID(stripeCustomerID string) error {
	_, _, err := client.From("organizations").
		Delete("", "").
		Eq("stripe_customer_id", stripeCustomerID).
		Execute()

	if err != nil {
		return err
	}

	return nil
}
