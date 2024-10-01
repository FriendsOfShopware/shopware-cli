package main

import (
	"encoding/json"
	"github.com/FriendsOfShopware/shopware-cli/shop"
	"github.com/invopop/jsonschema"
	"os"
)

func generateProjectSchema() error {
	r := new(jsonschema.Reflector)
	r.FieldNameTag = "yaml"
	r.RequiredFromJSONSchemaTags = true

	if err := r.AddGoComments("github.com/FriendsOfShopware/shopware-cli", "./shop"); err != nil {
		return err
	}

	schema := r.Reflect(&shop.Config{})

	bytes, err := json.MarshalIndent(schema, "", "  ")

	if err != nil {
		return err
	}

	if err := os.WriteFile("shop/shopware-project-schema.json", bytes, 0644); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := generateProjectSchema(); err != nil {
		panic(err)
	}
}
