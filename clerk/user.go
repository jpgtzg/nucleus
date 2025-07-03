package clerk

import (
	"context"
	"encoding/json"

	"github.com/clerk/clerk-sdk-go/v2/user"
)

func GetUserMetadata(userId string) (map[string]interface{}, error) {
	user, err := user.Get(context.Background(), userId)
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
	_, err = user.Update(context.Background(), userId, &user.UpdateParams{
		PublicMetadata: &rawMessage,
	})
	return err
}
