package api

import (
	"encoding/json"
	"net/http"
	"nucleus/auth"
	"nucleus/clerk"
)

// GetUserSuscriptionsHandler is a handler that returns the user's subscriptions
func GetUserSuscriptionsHandler(w http.ResponseWriter, r *http.Request) {
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
