package api

import (
	"encoding/json"
	"net/http"
	"nucleus/auth"
	"nucleus/clerk"
)

// GetUserSuscriptionsHandler is a handler that returns the user's subscriptions
func GetUserSuscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	subscriptions := clerk.GetActiveSubscriptions(userId)
	json.NewEncoder(w).Encode(subscriptions)
}
