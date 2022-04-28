package account

import (
	"github.com/spf13/cobra"
)

var accountCompanyProducerExtensionInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Manage store page",
}

func init() {
	accountCompanyProducerExtensionCmd.AddCommand(accountCompanyProducerExtensionInfoCmd)
}
