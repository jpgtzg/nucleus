package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"nucleus/stripe"

	"github.com/joho/godotenv"
	stripeSDK "github.com/stripe/stripe-go/v82"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	stripeSDK.Key = os.Getenv("STRIPE_KEY")

	http.HandleFunc("/webhook", stripe.HandleWebhook)
	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	fmt.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
