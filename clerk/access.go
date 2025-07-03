// TODO THIS FILE IS NOT IMPLEMENTED YET
package clerk

import (
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

// AccessControl provides methods to check user access to products
type AccessControl struct{}

// NewAccessControl creates a new access control instance
func NewAccessControl() *AccessControl {
	return &AccessControl{}
}

// HasAccess checks if a user has access to a specific product
// This is the main function you should use to check access in your application
func (ac *AccessControl) HasAccess(customerId string, productId string) bool {
	return HasActiveSubscription(customerId, productId)
}

// GetUserProducts returns all products the user has active access to
func (ac *AccessControl) GetUserProducts(customerId string) []string {
	activeSubscriptions := GetActiveSubscriptions(customerId)
	var products []string

	for _, sub := range activeSubscriptions {
		if productId, ok := sub["product_id"].(string); ok {
			products = append(products, productId)
		}
	}

	return products
}

// GetSubscriptionDetails returns detailed information about a user's subscription for a product
func (ac *AccessControl) GetSubscriptionDetails(customerId string, productId string) map[string]interface{} {
	activeSubscriptions := GetActiveSubscriptions(customerId)

	for _, sub := range activeSubscriptions {
		if sub["product_id"] == productId {
			return sub
		}
	}

	return nil
}

// ValidateAccess is a middleware-style function for HTTP handlers
func (ac *AccessControl) ValidateAccess(customerId string, requiredProduct string) (bool, string) {
	if !ac.HasAccess(customerId, requiredProduct) {
		return false, "Access denied: Active subscription required"
	}

	return true, "Access granted"
}
