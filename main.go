package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"nucleus/api"
	"nucleus/auth"
	"nucleus/clerk"
	"nucleus/stripe"

	"github.com/joho/godotenv"
	stripeSDK "github.com/stripe/stripe-go/v82"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	stripeSDK.Key = os.Getenv("STRIPE_KEY")

	http.HandleFunc("/stripe/webhook", stripe.HandleWebhook)
	http.HandleFunc("/clerk/webhook", clerk.WebhookHandler)
	http.Handle("/user/subscriptions", auth.VerifyingMiddleware(http.HandlerFunc(api.GetUserSuscriptionsHandler)))

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	fmt.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
