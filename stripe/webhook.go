package stripe

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

// HandleWebhook handles the webhook request
// It verifies that the request is coming from a webhook IP
// It reads the request body (if it's not too large) and constructs the event (if it's valid stripe event)
// It processes the event asynchronously
// It returns immediately with a 200 status code to acknowledge receipt
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Verify that the request is coming from a webhook IP
	clientIP := getClientIP(r)
	if !IsWebhookIP(clientIP) {
		log.Printf("Webhook request received from non-webhook IP: %s", clientIP)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// This part sets a maximum byte limit to the request body
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Verifies that the payload is coming from Stripe via the ConstructEvent function
	// Passes the payload to construct the Event (Go Stripe handler)
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	signatureHeader := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Webhook signature verification failed. %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	go processWebhookEvent(event)
}

// processWebhookEvent processes the webhook event asynchronously
// It handles subscription and invoice events
// It logs the event type if it's not handled
func processWebhookEvent(event stripe.Event) {
	log.Printf("[STRIPE] Processing webhook event: %s", event.Type)

	switch event.Type {
	case "customer.subscription.created":
		var subscription stripe.Subscription
		jsonData, err := json.Marshal(event.Data.Object)
		if err != nil {
			log.Printf("Error marshaling webhook data: %v", err)
			return
		}

		err = json.Unmarshal(jsonData, &subscription)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			return
		}
		HandleSubscriptionCreated(subscription)

	case "customer.subscription.updated":
		var subscription stripe.Subscription
		jsonData, err := json.Marshal(event.Data.Object)
		if err != nil {
			log.Printf("Error marshaling webhook data: %v", err)
			return
		}

		err = json.Unmarshal(jsonData, &subscription)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			return
		}
		HandleSubscriptionUpdated(subscription)

	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		jsonData, err := json.Marshal(event.Data.Object)
		if err != nil {
			log.Printf("Error marshaling webhook data: %v", err)
			return
		}

		err = json.Unmarshal(jsonData, &subscription)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			return
		}
		HandleSubscriptionDeleted(subscription)

	default:
		log.Printf("Unhandled event type: %s", event.Type)

		jsonData, err := json.Marshal(event.Data.Object)
		if err != nil {
			log.Printf("Error marshaling webhook data: %v", err)
			return
		}

		os.WriteFile(event.ID+"_"+string(event.Type)+".json", jsonData, 0644)

		return
	}
}

// getClientIP extracts the real client IP address from the request
// It checks various headers that might contain the real IP when behind a proxy
func getClientIP(r *http.Request) string {
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		if commaIndex := strings.Index(forwardedFor, ","); commaIndex != -1 {
			return strings.TrimSpace(forwardedFor[:commaIndex])
		}
		return strings.TrimSpace(forwardedFor)
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}

	if clientIP := r.Header.Get("X-Client-IP"); clientIP != "" {
		return strings.TrimSpace(clientIP)
	}

	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return strings.TrimSpace(cfIP)
	}
	return r.RemoteAddr
}
