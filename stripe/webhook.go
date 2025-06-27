package stripe

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
// It logs the event ID and type
// It handles the payment intent succeeded event
// It logs the event type if it's not handled
func processWebhookEvent(event stripe.Event) {
	log.Printf("Processing event ID: %s", event.ID)
	log.Printf("Event type: %s, Event created: %s", event.Type, time.Unix(event.Created, 0).Format("2006-01-02 15:04:05"))

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			return
		}
		log.Printf("Successful payment for %d.", paymentIntent.Amount)
		// TODO: Handle payment intent success
	default:
		log.Printf("Unhandled event type: %s", event.Type)
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
