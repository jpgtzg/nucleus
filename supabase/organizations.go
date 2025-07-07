package supabase

import (
	"encoding/json"
	"log"
	"os"

	supabaseTypes "nucleus/types/supabase"

	"github.com/supabase-community/supabase-go"
)

type Organization = supabaseTypes.Organization

var client *supabase.Client

func init() {
	var err error
	client, err = supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllOrganizations() ([]Organization, error) {
	data, _, err := client.From("organizations").Select("*", "exact", false).Execute()
	if err != nil {
		return nil, err
	}

	var organizations []Organization
	if err := json.Unmarshal(data, &organizations); err != nil {
		return nil, err
	}

	return organizations, nil
}
