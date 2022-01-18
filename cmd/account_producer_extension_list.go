package cmd

import (
	"github.com/pkg/errors"
	"os"
	account_api "shopware-cli/account-api"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var accountCompanyProducerExtensionListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all your extensions",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getAccountAPIByConfigOrFail()

		p, err := client.Producer()

		if err != nil {
			return errors.Wrap(err, "cannot get producer endpoint")
		}

		criteria := account_api.ListExtensionCriteria{
			Limit: 100,
		}

		if len(listExtensionSearch) > 0 {
			criteria.Search = listExtensionSearch
			criteria.OrderBy = "name"
			criteria.OrderSequence = "asc"
		}

		extensions, err := p.Extensions(&criteria)

		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Type", "Compatible with latest version", "Status"})

		for _, extension := range extensions {
			if extension.Status.Name == "deleted" {
				continue
			}

			compatible := "No"

			if extension.IsCompatibleWithLatestShopwareVersion {
				compatible = "Yes"
			}

			table.Append([]string{
				strconv.FormatInt(int64(extension.Id), 10),
				extension.Name,
				extension.Generation.Description,
				compatible,
				extension.Status.Name,
			})
		}

		table.Render()

		return nil
	},
}

var listExtensionSearch string

func init() {
	accountCompanyProducerExtensionCmd.AddCommand(accountCompanyProducerExtensionListCmd)
	accountCompanyProducerExtensionListCmd.Flags().StringVar(&listExtensionSearch, "search", "", "Filter for name")
}
