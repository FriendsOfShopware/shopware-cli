package cmd

import (
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var accountCompanyListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Lists all available company for your Account",
	Aliases: []string{"ls"},
	Long:    ``,
	Run: func(_ *cobra.Command, _ []string) {
		client := getAccountAPIByConfigOrFail()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Customer ID", "Roles"})

		for _, membership := range client.GetMemberships() {
			table.Append([]string{
				strconv.FormatInt(int64(membership.Company.Id), 10),
				membership.Company.Name,
				strconv.FormatInt(int64(membership.Company.CustomerNumber), 10),
				strings.Join(membership.GetRoles(), ", "),
			})
		}

		table.Render()
	},
}

func init() {
	accountCompanyRootCmd.AddCommand(accountCompanyListCmd)
}
