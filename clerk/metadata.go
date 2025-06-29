package clerk

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

var globalCtx context.Context

func init() {

	globalCtx = context.Background()
	clerkAPIKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkAPIKey == "" {
		log.Fatal("CLERK_SECRET_KEY environment variable is required")
	}

	clerk.SetKey(clerkAPIKey)
}

func UpdateUserMetadata(userId string, metadata map[string]string) error {
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
