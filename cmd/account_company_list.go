package cmd

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var accountCompanyListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all available company for your Account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := getAccountAPIByConfig()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Customer ID"})

		for _, membership := range client.GetMemberships() {
			table.Append([]string{strconv.FormatInt(int64(membership.Company.Id), 10), membership.Company.Name, strconv.FormatInt(int64(membership.Company.CustomerNumber), 10)})
		}

		table.Render()
	},
}

func init() {
	accountCompanyRootCmd.AddCommand(accountCompanyListCmd)
}
