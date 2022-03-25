package cmd

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var accountCompanyMerchantShopListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all shops",
	Aliases: []string{"ls"},
	RunE: func(_ *cobra.Command, _ []string) error {
		client := getAccountAPIByConfigOrFail()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Domain", "Usage"})

		shops, err := client.Merchant().Shops()

		if err != nil {
			return err
		}

		for _, shop := range shops {
			table.Append([]string{
				strconv.FormatInt(int64(shop.Id), 10),
				shop.Domain,
				shop.Environment.Name,
			})
		}

		table.Render()

		return nil
	},
}

func init() {
	accountCompanyMerchantShopCmd.AddCommand(accountCompanyMerchantShopListCmd)
}
