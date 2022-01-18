package cmd

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

var accountCompanyMerchantShopListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all shops",
	Run: func(cmd *cobra.Command, args []string) {
		client := getAccountAPIByConfig()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Domain", "Usage"})

		shops, err := client.Merchant().Shops()

		if err != nil {
			log.Fatalln(err)
		}

		for _, shop := range shops {
			table.Append([]string{
				strconv.FormatInt(int64(shop.Id), 10),
				shop.Domain,
				shop.Environment.Name,
			})
		}

		table.Render()
	},
}

func init() {
	accountCompanyMerchantShopCmd.AddCommand(accountCompanyMerchantShopListCmd)
}
