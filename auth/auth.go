package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nucleus/clerk"
	"strings"
	"time"

	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	clerkjwt "github.com/clerk/clerk-sdk-go/v2/jwt"
)

// OrganizationIDKey is the context key for storing organization ID
type OrganizationIDKey struct{}

// GetOrganizationID retrieves the organization ID from the request context
func GetOrganizationID(r *http.Request) (string, bool) {
	organizationID, ok := r.Context().Value(OrganizationIDKey{}).(string)
	return organizationID, ok
}

// VerifyingMiddleware is the general middleware that verifies the passed JWT Token from clerk and extracts the user ID and organization ID to pass it to the next handler
func VerifyingMiddleware(next http.Handler) http.Handler {
	return clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[API] Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		startTime := time.Now()

		userID, err := extractUserIDFromAuthHeader(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		organizationID, err := clerk.GetUserOrganizationId(userID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), OrganizationIDKey{}, organizationID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
		log.Printf("[API] Response: %s %s -> STATUS: %d completed in %v", r.Method, r.URL.Path, http.StatusOK, time.Since(startTime))
	}))
}

// extractUserIDFromAuthHeader extracts the user ID from the Authorization header
func extractUserIDFromAuthHeader(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	// Check if it's a Bearer token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Verify the JWT token and extract claims
	claims, err := clerkjwt.Verify(context.Background(), &clerkjwt.VerifyParams{
		Token: token,
	})
	if err != nil {
		return "", fmt.Errorf("failed to verify token: %v", err)
	}

	// Extract user ID from the subject claim
	userID := claims.RegisteredClaims.Subject
	if userID == "" {
		return "", fmt.Errorf("no user ID found in token")
	}

	return userID, nil
}
