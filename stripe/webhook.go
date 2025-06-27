package stripe

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
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
