package cmd

import (
	termColor "github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var accountCompanyProducerExtensionDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a extension",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getAccountApiByConfig()

		extensionId, err := strconv.Atoi(args[0])

		if err != nil {
			termColor.Red(err.Error())
			os.Exit(1)
		}

		p, err := client.Producer()

		if err != nil {
			termColor.Red(err.Error())
			os.Exit(1)
		}

		err = p.DeleteExtension(extensionId)

		if err != nil {
			termColor.Red(err.Error())
			os.Exit(1)
		}

		termColor.Green("Extension has been successfully deleted")
	},
}

func init() {
	accountCompanyProducerExtensionCmd.AddCommand(accountCompanyProducerExtensionDeleteCmd)
}
