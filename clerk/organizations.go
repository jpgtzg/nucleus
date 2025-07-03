package clerk

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/organization"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

func init() {
	clerkAPIKey := os.Getenv("CLERK_SECRET_KEY")
	clerk.SetKey(clerkAPIKey)
}

func GetUserOrganizations(userId string) (*clerk.OrganizationMembershipList, error) {
	orgMemberships, err := user.ListOrganizationMemberships(context.Background(), userId, &user.ListOrganizationMembershipsParams{})

	if err != nil {
		log.Printf("Error getting organization memberships: %v", err)
	}

	return orgMemberships, nil
}

func GetUserOrganizationId(userId string) (string, error) {
	orgMemberships, err := GetUserOrganizations(userId)
	if err != nil {
		return "", err
	}
	return orgMemberships.OrganizationMemberships[0].Organization.ID, nil
}

func GetOrganizationPublicMetadata(organizationId string) (map[string]interface{}, error) {
	organization, err := organization.Get(context.Background(), organizationId)
	if err != nil {
		return nil, err
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(organization.PublicMetadata, &metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func UpdateOrganizationPublicMetadata(organizationId string, metadata map[string]interface{}) error {
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	rawMessage := json.RawMessage(jsonData)
	_, err = organization.Update(context.Background(), organizationId, &organization.UpdateParams{
		PublicMetadata: &rawMessage,
	})

	return err
}
