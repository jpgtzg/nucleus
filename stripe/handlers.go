package stripe

import (
	"log"

	"github.com/stripe/stripe-go/v82"
)

// HandleSubscriptionCreated handles the subscription created event
// It adds the subscription information to the user metadata
func HandleSubscriptionCreated(subscription stripe.Subscription) {
	customerId := subscription.Customer.ID
	AddSubscriptionToUserMetadata(customerId, subscription)
	log.Printf("Subscription created for customer: %s, subscription: %s", customerId, subscription.ID)
}

// HandleSubscriptionUpdated handles the subscription updated event
// It updates the subscription information in the user metadata
func HandleSubscriptionUpdated(subscription stripe.Subscription) {
	customerId := subscription.Customer.ID
	UpdateSubscriptionInUserMetadata(customerId, subscription)
	log.Printf("Subscription updated for customer: %s, subscription: %s", customerId, subscription.ID)
}

// HandleSubscriptionDeleted handles the subscription deleted event
// It removes the subscription from the user metadata
func HandleSubscriptionDeleted(subscription stripe.Subscription) {
	customerId := subscription.Customer.ID
	RemoveSubscriptionFromUserMetadata(customerId, subscription.ID)
	log.Printf("Subscription deleted for customer: %s, subscription: %s", customerId, subscription.ID)
}
