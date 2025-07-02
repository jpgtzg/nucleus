// Package clerk is a wrapper for the Clerk API
// It provides functions to get and update user metadata

package clerk

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/joho/godotenv"
)

var globalCtx context.Context

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	globalCtx = context.Background()
	clerkAPIKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkAPIKey == "" {
		log.Fatal("CLERK_SECRET_KEY environment variable is required")
	}

	clerk.SetKey(clerkAPIKey)
}

func GetUserMetadata(userId string) (map[string]interface{}, error) {
	user, err := user.Get(globalCtx, userId)
	if err != nil {
		return nil, err
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(user.PublicMetadata, &metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func UpdateUserMetadata(userId string, metadata map[string]interface{}) error {
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	rawMessage := json.RawMessage(jsonData)
	_, err = user.Update(globalCtx, userId, &user.UpdateParams{
		PublicMetadata: &rawMessage,
	})
	return err
}
