package cmd

import (
	"os"

	termColor "github.com/fatih/color"
	"github.com/spf13/cobra"
)

var accountProducerInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "List information about your producer account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := getAccountAPIByConfigOrFail()

		p, err := client.Producer()

		if err != nil {
			termColor.Red(err.Error())
			os.Exit(1)
		}

		profile, err := p.Profile()
		if err != nil {
			termColor.Red(err.Error())
			os.Exit(1)
		}

		termColor.Blue("Name: %s", profile.Name)
		termColor.Blue("Prefix: %s", profile.Prefix)
		termColor.Blue("Website: %s", profile.Website)
	},
}

func init() {
	accountCompanyProducerCmd.AddCommand(accountProducerInfoCmd)
}
