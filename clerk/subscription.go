package clerk

import (
	"log"
	"nucleus/mongodb"

	"time"

	"github.com/stripe/stripe-go/v82"
)

// AddSubscriptionToOrganizationMetadata adds subscription information to user metadata
func AddSubscriptionToOrganizationMetadata(customerId string, subscription *stripe.Subscription) {
	organization, err := mongodb.GetOrganizationByStripeCustomerID(customerId)
	if err != nil {
		log.Printf("Error getting organization: %v", err)
		return
	}

	metadata, err := GetOrganizationPublicMetadata(organization.ClerkID)
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
						UpdateOrganizationPublicMetadata(organization.ClerkID, metadata)
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

	UpdateOrganizationPublicMetadata(organization.ClerkID, metadata)
}

// UpdateSubscriptionInOrganizationMetadata updates existing subscription information
func UpdateSubscriptionInOrganizationMetadata(customerId string, subscription *stripe.Subscription) {
	organization, err := mongodb.GetOrganizationByStripeCustomerID(customerId)
	if err != nil {
		log.Printf("Error getting organization: %v", err)
		return
	}

	metadata, err := GetOrganizationPublicMetadata(organization.ClerkID)
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
						UpdateOrganizationPublicMetadata(organization.ClerkID, metadata)
						return
					}
				}
			}
		}
	}
}

// RemoveSubscriptionFromOrganizationMetadata removes a subscription from user metadata
func RemoveSubscriptionFromOrganizationMetadata(customerId string, subscriptionId string) {
	organization, err := mongodb.GetOrganizationByStripeCustomerID(customerId)
	if err != nil {
		log.Printf("Error getting organization: %v", err)
		return
	}

	metadata, err := GetOrganizationPublicMetadata(organization.ClerkID)
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
			UpdateOrganizationPublicMetadata(organization.ClerkID, metadata)
		}
	}
}

// GetActiveSubscriptionsByCustomerID returns all active subscriptions for a organization
func GetActiveSubscriptionsByCustomerID(customerId string) []map[string]interface{} {
	organization, err := mongodb.GetOrganizationByStripeCustomerID(customerId)
	if err != nil {
		log.Printf("Error getting organization: %v", err)
		return nil
	}

	metadata, err := GetOrganizationPublicMetadata(organization.ClerkID)
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

// GetActiveSubscriptions returns all active subscriptions for a organization
func GetActiveSubscriptionsByOrganizationID(organizationID string) []map[string]interface{} {

	metadata, err := GetOrganizationPublicMetadata(organizationID)
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
