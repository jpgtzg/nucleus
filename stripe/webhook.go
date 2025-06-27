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
// It reads the request body and constructs the event
// It processes the event asynchronously
// It returns a 200 status code to acknowledge receipt
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Verify that the request is coming from a webhook IP
	clientIP := getClientIP(r)
	if !IsWebhookIP(clientIP) {
		log.Printf("Webhook request received from non-webhook IP: %s", clientIP)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Passes the payload to construct the Event (Go Stripe handler), also verifies that the payload is coming from Stripe
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	signatureHeader := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Webhook signature verification failed. %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Return HTTP 200 immediately to acknowledge receipt
	w.WriteHeader(http.StatusOK)

	// Process the webhook event asynchronously
	go processWebhookEvent(event)
}

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
	// Check for X-Forwarded-For header (most common proxy header)
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs separated by commas
		// The first one is usually the original client IP
		if commaIndex := strings.Index(forwardedFor, ","); commaIndex != -1 {
			return strings.TrimSpace(forwardedFor[:commaIndex])
		}
		return strings.TrimSpace(forwardedFor)
	}

	// Check for X-Real-IP header (used by some reverse proxies)
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}

	// Check for X-Client-IP header
	if clientIP := r.Header.Get("X-Client-IP"); clientIP != "" {
		return strings.TrimSpace(clientIP)
	}

	// Check for CF-Connecting-IP header (Cloudflare)
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return strings.TrimSpace(cfIP)
	}

	// Fall back to RemoteAddr if no proxy headers are present
	return r.RemoteAddr
}
