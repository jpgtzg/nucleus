package stripe

import (
	"log"
	"nucleus/clerk"
	"time"

	"github.com/stripe/stripe-go/v82"
)

// AddSubscriptionToUserMetadata adds subscription information to user metadata
func AddSubscriptionToUserMetadata(customerId string, subscription stripe.Subscription) {
	metadata, err := clerk.GetUserMetadata(customerId)
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
		"id":                   subscription.ID,
		"status":               subscription.Status,
		"current_period_end":   currentPeriodEnd,
		"cancel_at_period_end": subscription.CancelAtPeriodEnd,
		"product_id":           subscription.Items.Data[0].Price.Product.ID,
		"price_id":             subscription.Items.Data[0].Price.ID,
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
						subMap["cancel_at_period_end"] = subscription.CancelAtPeriodEnd
						subMap["product_id"] = subscription.Items.Data[0].Price.Product.ID
						subMap["price_id"] = subscription.Items.Data[0].Price.ID
						clerk.UpdateUserMetadata(customerId, metadata)
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

	clerk.UpdateUserMetadata(customerId, metadata)
}

// UpdateSubscriptionInUserMetadata updates existing subscription information
func UpdateSubscriptionInUserMetadata(customerId string, subscription stripe.Subscription) {
	metadata, err := clerk.GetUserMetadata(customerId)
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
							"id":                   subscription.ID,
							"status":               subscription.Status,
							"current_period_end":   currentPeriodEnd,
							"cancel_at_period_end": subscription.CancelAtPeriodEnd,
							"product_id":           subscription.Items.Data[0].Price.Product.ID,
							"price_id":             subscription.Items.Data[0].Price.ID,
						}
						clerk.UpdateUserMetadata(customerId, metadata)
						return
					}
				}
			}
		}
	}
}

// RemoveSubscriptionFromUserMetadata removes a subscription from user metadata
func RemoveSubscriptionFromUserMetadata(customerId string, subscriptionId string) {
	metadata, err := clerk.GetUserMetadata(customerId)
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
			clerk.UpdateUserMetadata(customerId, metadata)
		}
	}
}

// HasActiveSubscription checks if a user has an active subscription for a specific product
// Returns true if the user has an active subscription that hasn't expired
func HasActiveSubscription(customerId string, productId string) bool {
	metadata, err := clerk.GetUserMetadata(customerId)
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

// GetActiveSubscriptions returns all active subscriptions for a user
func GetActiveSubscriptions(customerId string) []map[string]interface{} {
	metadata, err := clerk.GetUserMetadata(customerId)
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
