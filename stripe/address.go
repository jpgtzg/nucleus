package stripe

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// fetchWebhookIPs fetches the list of webhook IPs from the Stripe API
func fetchWebhookIPs() []string {
	response, err := http.Get("https://stripe.com/files/ips/ips_webhooks.json")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data struct {
		Webhooks []string `json:"WEBHOOKS"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data.Webhooks
}

// isWebhookIP checks if the given IP is in the list of webhook IPs
func IsWebhookIP(ip string) bool {
	ips := fetchWebhookIPs()
	for _, webhookIP := range ips {
		if webhookIP == ip {
			return true
		}
	}
	return false
}
