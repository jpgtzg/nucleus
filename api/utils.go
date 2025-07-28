package api

import (
	"net/http"
	"strings"
)

// VerifyMethod checks if the HTTP request method is in the list of allowed methods.
// This function is used to ensure endpoints only accept the correct HTTP methods.
//
// Parameters:
//   - r: HTTP request to verify
//   - allowedMethods: Array of allowed HTTP methods (e.g., ["GET", "POST"])
//
// Returns:
//   - bool: True if the request method is allowed, false otherwise
//
// Example:
//
//	VerifyMethod(r, []string{"POST"}) // Only allows POST
//	VerifyMethod(r, []string{"GET", "POST"}) // Allows both GET and POST
func VerifyMethod(r *http.Request, allowedMethods []string) bool {
	for _, method := range allowedMethods {
		if r.Method == strings.ToUpper(method) {
			return true
		}
	}
	return false
}
