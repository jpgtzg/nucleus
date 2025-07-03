package clerk

import (
	"encoding/json"
)

// ClerkWebhookEvent represents the structure of a Clerk webhook event
type ClerkWebhookEvent struct {
	Data      json.RawMessage `json:"data"`
	Object    string          `json:"object"`
	Type      string          `json:"type"`
	CreatedAt int64           `json:"created_at"`
}
