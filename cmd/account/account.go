package account

import (
	"github.com/spf13/cobra"
	account_api "shopware-cli/account-api"
	"shopware-cli/config"
)

var accountRootCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage your Shopware Account",
}

type ServiceContainer struct {
	Conf          config.Config
	AccountClient *account_api.Client
}

var services *ServiceContainer

func Register(rootCmd *cobra.Command, onInit func() (*ServiceContainer, error)) {
	accountRootCmd.PreRunE = func(_ *cobra.Command, _ []string) error {
		ser, err := onInit()
		services = ser
		return err
	}
	rootCmd.AddCommand(accountRootCmd)
}
