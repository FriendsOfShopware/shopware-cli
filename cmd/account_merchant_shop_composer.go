package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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

		return nil
	},
}

func init() {
	accountCompanyMerchantShopCmd.AddCommand(accountCompanyMerchantShopComposerCmd)
}
