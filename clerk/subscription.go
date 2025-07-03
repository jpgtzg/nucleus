package clerk

import (
	"log"
	"time"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
)

// AddSubscriptionToOrganizationMetadata adds subscription information to user metadata
func AddSubscriptionToOrganizationMetadata(customerId string, subscription stripe.Subscription) {
	stripeCustomer, err := customer.Get(customerId, &stripe.CustomerParams{})
	if err != nil {
		log.Printf("Error getting stripe customer: %v", err)
		return
	}

	organizationId := stripeCustomer.Metadata["clerk_organization_id"]
	if organizationId == "" {
		log.Printf("Stripe customer %s has no clerk organization id", customerId)
		return
	}

	metadata, err := GetOrganizationMetadata(organizationId)
	if err != nil {
		log.Printf("Error getting user metadata: %v", err)
		return
	}

	// Get current period end from subscription items
	var currentPeriodEnd int64
	if len(subscription.Items.Data) > 0 {
		currentPeriodEnd = subscription.Items.Data[0].CurrentPeriodEnd
	}

	// Create subscription info
	subscriptionInfo := map[string]interface{}{
		"id":                 subscription.ID,
		"status":             subscription.Status,
		"current_period_end": currentPeriodEnd,
		"product_id":         subscription.Items.Data[0].Price.Product.ID,
		"price_id":           subscription.Items.Data[0].Price.ID,
	}

	// Initialize stripe data if it doesn't exist
	if stripeData, ok := metadata["stripe"].(map[string]interface{}); ok {
		if subscriptions, ok := stripeData["subscriptions"].([]interface{}); ok {
			// Check if subscription already exists
			for _, sub := range subscriptions {
				if subMap, ok := sub.(map[string]interface{}); ok {
					if subMap["id"] == subscription.ID {
						// Update existing subscription
						subMap["status"] = subscription.Status
						subMap["current_period_end"] = currentPeriodEnd
						subMap["product_id"] = subscription.Items.Data[0].Price.Product.ID
						subMap["price_id"] = subscription.Items.Data[0].Price.ID
						UpdateOrganizationMetadata(organizationId, metadata)
						return
					}
				}
			}
			// Add new subscription
			stripeData["subscriptions"] = append(subscriptions, subscriptionInfo)
		} else {
			stripeData["subscriptions"] = []interface{}{subscriptionInfo}
		}
	} else {
		metadata["stripe"] = map[string]interface{}{
			"subscriptions": []interface{}{subscriptionInfo},
		}
	}

	UpdateOrganizationMetadata(organizationId, metadata)
}

// UpdateSubscriptionInOrganizationMetadata updates existing subscription information
func UpdateSubscriptionInOrganizationMetadata(customerId string, subscription stripe.Subscription) {
	stripeCustomer, err := customer.Get(customerId, &stripe.CustomerParams{})
	if err != nil {
		log.Printf("Error getting stripe customer: %v", err)
		return
	}

	organizationId := stripeCustomer.Metadata["clerk_organization_id"]
	if organizationId == "" {
		log.Printf("Stripe customer %s has no clerk organization id", customerId)
		return
	}

	metadata, err := GetOrganizationMetadata(organizationId)
	if err != nil {
		log.Printf("Error getting user metadata: %v", err)
		return
	}

	// Get current period end from subscription items
	var currentPeriodEnd int64
	if len(subscription.Items.Data) > 0 {
		currentPeriodEnd = subscription.Items.Data[0].CurrentPeriodEnd
	}

	if stripeData, ok := metadata["stripe"].(map[string]interface{}); ok {
		if subscriptions, ok := stripeData["subscriptions"].([]interface{}); ok {
			for i, sub := range subscriptions {
				if subMap, ok := sub.(map[string]interface{}); ok {
					if subMap["id"] == subscription.ID {
						// Update subscription info
						subscriptions[i] = map[string]interface{}{
							"id":                 subscription.ID,
							"status":             subscription.Status,
							"current_period_end": currentPeriodEnd,
							"product_id":         subscription.Items.Data[0].Price.Product.ID,
							"price_id":           subscription.Items.Data[0].Price.ID,
						}
						UpdateOrganizationMetadata(organizationId, metadata)
						return
					}
				}
			}
		}
	}
}

// RemoveSubscriptionFromOrganizationMetadata removes a subscription from user metadata
func RemoveSubscriptionFromOrganizationMetadata(customerId string, subscriptionId string) {
	stripeCustomer, err := customer.Get(customerId, &stripe.CustomerParams{})
	if err != nil {
		log.Printf("Error getting stripe customer: %v", err)
		return
	}

	organizationId := stripeCustomer.Metadata["clerk_organization_id"]
	if organizationId == "" {
		log.Printf("Stripe customer %s has no clerk organization id", customerId)
		return
	}

	metadata, err := GetOrganizationMetadata(organizationId)
	if err != nil {
		log.Printf("Error getting user metadata: %v", err)
		return
	}

	if stripeData, ok := metadata["stripe"].(map[string]interface{}); ok {
		if subscriptions, ok := stripeData["subscriptions"].([]interface{}); ok {
			var updatedSubscriptions []interface{}
			for _, sub := range subscriptions {
				if subMap, ok := sub.(map[string]interface{}); ok {
					if subMap["id"] != subscriptionId {
						updatedSubscriptions = append(updatedSubscriptions, sub)
					}
				}
			}
			stripeData["subscriptions"] = updatedSubscriptions
			UpdateOrganizationMetadata(organizationId, metadata)
		}
	}
}

// HasActiveSubscription checks if a organization has an active subscription for a specific product
// Returns true if the organization has an active subscription that hasn't expired
func HasActiveSubscription(customerId string, productId string) bool {
	organizationId, err := GetUserOrganizationId(customerId)
	if err != nil {
		log.Printf("Error getting user organization id: %v", err)
		return false
	}

	metadata, err := GetOrganizationMetadata(organizationId)
	if err != nil {
		log.Printf("Error getting user metadata: %v", err)
		return false
	}

	if stripeData, ok := metadata["stripe"].(map[string]interface{}); ok {
		if subscriptions, ok := stripeData["subscriptions"].([]interface{}); ok {
			currentTime := time.Now().Unix()

			for _, sub := range subscriptions {
				if subMap, ok := sub.(map[string]interface{}); ok {
					// Check if this subscription is for the requested product
					if subMap["product_id"] == productId {
						// Check if subscription is active
						if status, ok := subMap["status"].(string); ok && status == "active" {
							// Check if subscription hasn't expired
							if periodEnd, ok := subMap["current_period_end"].(float64); ok {
								if int64(periodEnd) > currentTime {
									return true
								}
							}
						}
					}
				}
			}
		}
	}

	return false
}

// GetActiveSubscriptions returns all active subscriptions for a organization
func GetActiveSubscriptions(customerId string) []map[string]interface{} {
	organizationId, err := GetUserOrganizationId(customerId)
	if err != nil {
		log.Printf("Error getting user organization id: %v", err)
		return nil
	}

	metadata, err := GetOrganizationMetadata(organizationId)
	if err != nil {
		log.Printf("Error getting user metadata: %v", err)
		return nil
	}

	var activeSubscriptions []map[string]interface{}
	currentTime := time.Now().Unix()

	if stripeData, ok := metadata["stripe"].(map[string]interface{}); ok {
		if subscriptions, ok := stripeData["subscriptions"].([]interface{}); ok {
			for _, sub := range subscriptions {
				if subMap, ok := sub.(map[string]interface{}); ok {
					// Check if subscription is active
					if status, ok := subMap["status"].(string); ok && status == "active" {
						// Check if subscription hasn't expired
						if periodEnd, ok := subMap["current_period_end"].(float64); ok {
							if int64(periodEnd) > currentTime {
								activeSubscriptions = append(activeSubscriptions, subMap)
							}
						}
					}
				}
			}
		}
	}

	return activeSubscriptions
}
