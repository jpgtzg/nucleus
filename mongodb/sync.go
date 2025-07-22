package mongodb

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	mongodbTypes "nucleus/types/mongodb"
)

var Client *mongo.Client

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URI")).SetServerAPIOptions(serverAPI)

	var err error
	Client, err = mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	if err := Client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("[MONGO] Pinged deployment. Successfully connected to MongoDB!")
}

func CreateOrganizationSync(clerkID string, stripeCustomerID string) error {

	coll := Client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION_SYNC"))

	_, err := coll.InsertOne(context.TODO(), bson.M{
		"clerk_organization_id": clerkID,
		"stripe_customer_id":    stripeCustomerID,
	})

	log.Printf("[MONGO] Created organization sync for clerkID: %s, stripeCustomerID: %s", clerkID, stripeCustomerID)
	return err
}

func GetOrganizationByClerkID(clerkID string) (mongodbTypes.Organization, error) {

	coll := Client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION_SYNC"))

	var result mongodbTypes.Organization
	err := coll.FindOne(context.Background(), bson.M{"clerk_organization_id": clerkID}).Decode(&result)
	if err != nil {
		return mongodbTypes.Organization{}, err
	}

	return result, nil
}

func GetOrganizationByStripeCustomerID(stripeCustomerID string) (mongodbTypes.Organization, error) {
	coll := Client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION_SYNC"))

	var result mongodbTypes.Organization
	err := coll.FindOne(context.Background(), bson.M{"stripe_customer_id": stripeCustomerID}).Decode(&result)
	if err != nil {
		return mongodbTypes.Organization{}, err
	}

	return result, nil
}

func DeleteOrganizationByClerkID(clerkID string) error {
	coll := Client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION_SYNC"))

	_, err := coll.DeleteOne(context.Background(), bson.M{"clerk_organization_id": clerkID})
	if err != nil {
		return err
	}

	return nil
}

func DeleteOrganizationByStripeCustomerID(stripeCustomerID string) error {
	coll := Client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION_SYNC"))

	_, err := coll.DeleteOne(context.Background(), bson.M{"stripe_customer_id": stripeCustomerID})
	if err != nil {
		return err
	}

	return nil
}

func UpdateOrganizationStripeCustomerID(clerkID string, stripeCustomerID string) error {
	coll := Client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION_SYNC"))

	_, err := coll.UpdateOne(context.Background(), bson.M{"clerk_organization_id": clerkID}, bson.M{"$set": bson.M{"stripe_customer_id": stripeCustomerID}})
	if err != nil {
		return err
	}

	return nil
}

func UpdateOrganizationClerkID(stripeCustomerID string, clerkID string) error {
	coll := Client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION_SYNC"))

	_, err := coll.UpdateOne(context.Background(), bson.M{"stripe_customer_id": stripeCustomerID}, bson.M{"$set": bson.M{"clerk_organization_id": clerkID}})
	if err != nil {
		return err
	}

	return nil
}
