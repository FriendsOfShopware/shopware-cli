package account

import (
	"github.com/spf13/cobra"
)

var accountCompanyMerchantShopCmd = &cobra.Command{
	Use:   "shop",
	Short: "Manage the shops",
}

func init() {
	accountCompanyMerchantCmd.AddCommand(accountCompanyMerchantShopCmd)
}
