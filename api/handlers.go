package api

import (
	"encoding/json"
	"net/http"
	"nucleus/auth"
	"nucleus/clerk"
	"nucleus/mongodb"
)

// GetUserSuscriptionsHandler is a handler that returns the user's subscriptions
func GetUserSuscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the organization ID from the user's organization memberships
	organizationID, ok := auth.GetOrganizationID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the active subscriptions for the organization
	subscriptions := clerk.GetActiveSubscriptionsByOrganizationID(organizationID)

	// Return the subscriptions
	json.NewEncoder(w).Encode(subscriptions)
}

func GetUserStripeCustomerIDHandler(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the organization ID from the user's organization memberships
	organizationID, ok := auth.GetOrganizationID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the organization object from the database using the clerkID
	organization, err := mongodb.GetOrganizationByClerkID(organizationID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return the stripe customer ID
	json.NewEncoder(w).Encode(organization.StripeCustomerID)
}
