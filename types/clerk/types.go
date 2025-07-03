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

// SessionData represents session data in webhook events
type SessionData struct {
	ID              string                 `json:"id"`
	Object          string                 `json:"object"`
	UserID          string                 `json:"user_id"`
	Status          string                 `json:"status"`
	LastActiveAt    int64                  `json:"last_active_at"`
	ExpireAt        int64                  `json:"expire_at"`
	AbandonAt       int64                  `json:"abandon_at"`
	CreatedAt       int64                  `json:"created_at"`
	UpdatedAt       int64                  `json:"updated_at"`
	PublicMetadata  map[string]interface{} `json:"public_metadata"`
	PrivateMetadata map[string]interface{} `json:"private_metadata"`
	Actor           *Actor                 `json:"actor,omitempty"`
}

// Actor represents an actor in session events
type Actor struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Action    string `json:"action"`
	CreatedAt int64  `json:"created_at"`
}

// EmailData represents email data in webhook events
type EmailData struct {
	ID             string                 `json:"id"`
	Object         string                 `json:"object"`
	FromEmailName  string                 `json:"from_email_name"`
	Subject        string                 `json:"subject"`
	Body           string                 `json:"body"`
	BodyPlain      *string                `json:"body_plain,omitempty"`
	EmailAddress   string                 `json:"email_address"`
	Status         string                 `json:"status"`
	CreatedAt      int64                  `json:"created_at"`
	UpdatedAt      int64                  `json:"updated_at"`
	DeliveredAt    *int64                 `json:"delivered_at,omitempty"`
	OpenedAt       *int64                 `json:"opened_at,omitempty"`
	ClickedAt      *int64                 `json:"clicked_at,omitempty"`
	BouncedAt      *int64                 `json:"bounced_at,omitempty"`
	ComplainedAt   *int64                 `json:"complained_at,omitempty"`
	UnsubscribedAt *int64                 `json:"unsubscribed_at,omitempty"`
	Data           map[string]interface{} `json:"data"`
}
