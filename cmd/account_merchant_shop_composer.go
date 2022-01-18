package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var accountCompanyMerchantShopComposerCmd = &cobra.Command{
	Use:   "configure-composer [domain]",
	Short: "Configure local composer.json to use packages.shopware.com",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		completions := make([]string, 0)

		client, err := getAccountAPIByConfig()

		if err != nil {
			return completions, cobra.ShellCompDirectiveNoFileComp
		}

		shops, err := client.Merchant().Shops()

		if err != nil {
			return completions, cobra.ShellCompDirectiveNoFileComp
		}

		for _, shop := range shops {
			completions = append(completions, shop.Domain)
		}

		return completions, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getAccountAPIByConfigOrFail()

		shops, err := client.Merchant().Shops()

		if err != nil {
			return errors.Wrap(err, "cannot get shops")
		}

		shop := shops.GetByDomain(args[0])

		if shop == nil {
			return fmt.Errorf("cannot find shop by domain %s", args[0])
		}

		token, err := client.Merchant().GetComposerToken(shop.Id)

		if err != nil {
			return err
		}

		if token == "" {
			generatedToken, err := client.Merchant().GenerateComposerToken(shop.Id)

			if err != nil {
				return err
			}

			if err := client.Merchant().SaveComposerToken(shop.Id, generatedToken); err != nil {
				return err
			}

			token = generatedToken
		}

		log.Infof("The composer token is %s", token)

		if _, err := os.Stat("composer.json"); err == nil {
			log.Info("Found composer.json, adding it now as repository")

			var content []byte

			if content, err = ioutil.ReadFile("composer.json"); err != nil {
				return err
			}

			var composer map[string]interface{}

			if err := json.Unmarshal(content, &composer); err != nil {
				return err
			}

			if _, ok := composer["repositories"]; !ok {
				composer["repositories"] = make(map[string]interface{})
			}

			repositories, _ := composer["repositories"].(map[string]interface{})

			repositories["shopware-packages"] = struct {
				Type string `json:"type"`
				Url  string `json:"url"`
			}{
				Type: "composer",
				Url:  "https://packages.shopware.com",
			}

			if content, err = json.MarshalIndent(composer, "", "    "); err != nil {
				return err
			}

			if err = ioutil.WriteFile("composer.json", content, os.ModePerm); err != nil {
				return err
			}

			var authJson map[string]interface{}

			if content, err = ioutil.ReadFile("auth.json"); err == nil {
				if err := json.Unmarshal(content, &authJson); err != nil {
					return err
				}
			} else {
				authJson = make(map[string]interface{})
			}

			if _, ok := authJson["bearer"]; !ok {
				authJson["bearer"] = make(map[string]interface{})
			}

			bearer, _ := authJson["bearer"].(map[string]interface{})

			bearer["packages.shopware.com"] = token

			if content, err = json.MarshalIndent(authJson, "", "    "); err != nil {
				return err
			}

			if err = ioutil.WriteFile("auth.json", content, os.ModePerm); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	accountCompanyMerchantShopCmd.AddCommand(accountCompanyMerchantShopComposerCmd)
}
