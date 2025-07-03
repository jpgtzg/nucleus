package stripe

import (
	"log"
	"nucleus/clerk"

	"github.com/stripe/stripe-go/v82"
)

// HandleSubscriptionCreated handles the subscription created event
// It adds the subscription information to the organization metadata
func HandleSubscriptionCreated(subscription stripe.Subscription) {
	customerId := subscription.Customer.ID
	clerk.AddSubscriptionToOrganizationMetadata(customerId, subscription)
	log.Printf("Subscription created for customer: %s, subscription: %s", customerId, subscription.ID)
}

// HandleSubscriptionUpdated handles the subscription updated event
// It updates the subscription information in the organization metadata
func HandleSubscriptionUpdated(subscription stripe.Subscription) {
	customerId := subscription.Customer.ID
	clerk.UpdateSubscriptionInOrganizationMetadata(customerId, subscription)
	log.Printf("Subscription updated for customer: %s, subscription: %s", customerId, subscription.ID)
}

// HandleSubscriptionDeleted handles the subscription deleted event
// It removes the subscription from the organization metadata
func HandleSubscriptionDeleted(subscription stripe.Subscription) {
	customerId := subscription.Customer.ID
	clerk.RemoveSubscriptionFromOrganizationMetadata(customerId, subscription.ID)
	log.Printf("Subscription deleted for customer: %s, subscription: %s", customerId, subscription.ID)
}
