package account

import (
	"github.com/spf13/cobra"
)

var accountCompanyMerchantCmd = &cobra.Command{
	Use:   "merchant",
	Short: "Manage merchants",
}

func init() {
	accountRootCmd.AddCommand(accountCompanyMerchantCmd)
}
