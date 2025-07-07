package supabase

type Organization struct {
	ID               int    `json:"id"`
	ClerkID          string `json:"clerk_organization_id"`
	StripeCustomerID string `json:"stripe_customer_id"`
}
