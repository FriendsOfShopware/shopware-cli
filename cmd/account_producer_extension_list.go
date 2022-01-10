package cmd

import (
	termColor "github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

var accountCompanyProducerExtensionListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all your extensions",
	Run: func(cmd *cobra.Command, args []string) {
		client := getAccountApiByConfig()

		p, err := client.Producer()

		if err != nil {
			termColor.Red(err.Error())
			os.Exit(1)
		}

		extensions, err := p.Extensions()

		if err != nil {
			log.Fatalln(err)
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
	},
}

func init() {
	accountCompanyProducerExtensionCmd.AddCommand(accountCompanyProducerExtensionListCmd)
}
