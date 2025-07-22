package mongodb

import "go.mongodb.org/mongo-driver/v2/bson"

type Organization struct {
	ID               bson.ObjectID `json:"id" bson:"_id"`
	ClerkID          string        `json:"clerk_organization_id" bson:"clerk_organization_id"`
	StripeCustomerID string        `json:"stripe_customer_id" bson:"stripe_customer_id"`
}
