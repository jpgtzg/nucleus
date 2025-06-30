package stripe

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"nucleus/clerk"
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
// It handles the invoice paid event
// It logs the event type if it's not handled
func processWebhookEvent(event stripe.Event) {

	switch event.Type {
	case "invoice.paid":
		var invoice stripe.Invoice
		jsonData, err := json.Marshal(event.Data.Object)
		if err != nil {
			log.Printf("Error marshaling webhook data: %v", err)
			return
		}

		err = json.Unmarshal(jsonData, &invoice)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			return
		}
		handleInvoicePaid(invoice)
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

// handleInvoicePaid handles the invoice paid event
// It adds the product ID to the user metadata (it requires the stripe ID to be the same as the clerk user ID - front end task)
func handleInvoicePaid(invoice stripe.Invoice) {
	productId := invoice.Lines.Data[0].Pricing.PriceDetails.Product
	customerId := invoice.Customer.ID
	addProductIdToUserMetadata(customerId, productId)
	log.Println(" Stripe User ID:", invoice.Customer.ID)
}

// addProductIdToUserMetadata adds the product ID to the user metadata
// It gets the user metadata, adds the product ID to the products_id array, and updates the user metadata
func addProductIdToUserMetadata(customerId string, productId string) {

	metadata, err := clerk.GetUserMetadata(customerId)
	if err != nil {
		log.Printf("Error getting user metadata: %v", err)
		return
	}

	// appends the current metadata with the new product ID, if it doesn't exist, it creates a new products_id array
	if stripeData, ok := metadata["stripe"].(map[string]interface{}); ok {
		if productsID, ok := stripeData["products_id"].([]interface{}); ok {
			stripeData["products_id"] = append(productsID, productId)
		} else {
			stripeData["products_id"] = []interface{}{productId}
		}
	} else {
		metadata["stripe"] = map[string]interface{}{
			"products_id": []interface{}{productId},
		}
	}

	clerk.UpdateUserMetadata(customerId, metadata)
}
