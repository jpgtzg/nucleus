package clerk

import (
	"encoding/json"
	"time"
)

// ClerkWebhookEvent represents the structure of a Clerk webhook event
type ClerkWebhookEvent struct {
	Data      json.RawMessage `json:"data"`
	Object    string          `json:"object"`
	Type      string          `json:"type"`
	CreatedAt int64           `json:"created_at"`
}

// UserData represents user data in webhook events
type UserData struct {
	ID                string                 `json:"id"`
	EmailAddresses    []EmailAddress         `json:"email_addresses,omitempty"`
	PhoneNumbers      []PhoneNumber          `json:"phone_numbers,omitempty"`
	Web3Wallets       []Web3Wallet           `json:"web3_wallets,omitempty"`
	Username          *string                `json:"username,omitempty"`
	FirstName         *string                `json:"first_name,omitempty"`
	LastName          *string                `json:"last_name,omitempty"`
	ImageURL          string                 `json:"image_url,omitempty"`
	HasImage          bool                   `json:"has_image"`
	CreatedAt         int64                  `json:"created_at"`
	UpdatedAt         int64                  `json:"updated_at"`
	PublicMetadata    map[string]interface{} `json:"public_metadata,omitempty"`
	PrivateMetadata   map[string]interface{} `json:"private_metadata,omitempty"`
	UnsafeMetadata    map[string]interface{} `json:"unsafe_metadata,omitempty"`
	ExternalID        *string                `json:"external_id,omitempty"`
	ExternalAccounts  []ExternalAccount      `json:"external_accounts,omitempty"`
	LastSignInAt      *int64                 `json:"last_sign_in_at,omitempty"`
	Banned            bool                   `json:"banned"`
	Locked            bool                   `json:"locked"`
	LockoutUntil      *int64                 `json:"lockout_until,omitempty"`
	Verification      UserVerification       `json:"verification,omitempty"`
	PasswordEnabled   bool                   `json:"password_enabled"`
	TOTPEnabled       bool                   `json:"totp_enabled"`
	BackupCodeEnabled bool                   `json:"backup_code_enabled"`
	TwoFactorEnabled  bool                   `json:"two_factor_enabled"`
	Gender            *string                `json:"gender,omitempty"`
	Birthday          *string                `json:"birthday,omitempty"`
}

// EmailAddress represents an email address
type EmailAddress struct {
	ID           string `json:"id"`
	EmailAddress string `json:"email_address"`
	Verification struct {
		Status   string `json:"status"`
		Strategy string `json:"strategy"`
	} `json:"verification"`
	LinkedTo []struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"linked_to,omitempty"`
	Object string `json:"object"`
}

// PhoneNumber represents a phone number
type PhoneNumber struct {
	ID           string `json:"id"`
	PhoneNumber  string `json:"phone_number"`
	Verification struct {
		Status   string `json:"status"`
		Strategy string `json:"strategy"`
	} `json:"verification"`
	LinkedTo []struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"linked_to,omitempty"`
	Object string `json:"object"`
}

// Web3Wallet represents a web3 wallet
type Web3Wallet struct {
	ID           string `json:"id"`
	Web3Wallet   string `json:"web3_wallet"`
	Verification struct {
		Status   string `json:"status"`
		Strategy string `json:"strategy"`
	} `json:"verification"`
	Object string `json:"object"`
}

// ExternalAccount represents an external account
type ExternalAccount struct {
	ID              string                 `json:"id"`
	Object          string                 `json:"object"`
	Provider        string                 `json:"provider"`
	EmailAddress    string                 `json:"email_address"`
	ApprovedScopes  string                 `json:"approved_scopes"`
	CreatedAt       int64                  `json:"created_at"`
	UpdatedAt       int64                  `json:"updated_at"`
	PublicMetadata  map[string]interface{} `json:"public_metadata"`
	PrivateMetadata map[string]interface{} `json:"private_metadata"`
	Label           *string                `json:"label,omitempty"`
	Verification    struct {
		Status   string `json:"status"`
		Strategy string `json:"strategy"`
		Attempts *int   `json:"attempts,omitempty"`
		ExpireAt *int64 `json:"expire_at,omitempty"`
	} `json:"verification"`
}

// UserVerification represents user verification status
type UserVerification struct {
	Status   string `json:"status"`
	Strategy string `json:"strategy"`
	Attempts *int   `json:"attempts,omitempty"`
	ExpireAt *int64 `json:"expire_at,omitempty"`
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

// GetPrimaryEmail returns the primary email address from user data
func (u *UserData) GetPrimaryEmail() string {
	if len(u.EmailAddresses) > 0 {
		return u.EmailAddresses[0].EmailAddress
	}
	return ""
}

// GetFullName returns the full name from user data
func (u *UserData) GetFullName() string {
	firstName := ""
	lastName := ""

	if u.FirstName != nil {
		firstName = *u.FirstName
	}
	if u.LastName != nil {
		lastName = *u.LastName
	}

	if firstName != "" && lastName != "" {
		return firstName + " " + lastName
	} else if firstName != "" {
		return firstName
	} else if lastName != "" {
		return lastName
	}
	return ""
}

// GetCreatedAt returns the creation time as a time.Time
func (u *UserData) GetCreatedAt() time.Time {
	return time.Unix(u.CreatedAt, 0)
}

// GetUpdatedAt returns the update time as a time.Time
func (u *UserData) GetUpdatedAt() time.Time {
	return time.Unix(u.UpdatedAt, 0)
}
