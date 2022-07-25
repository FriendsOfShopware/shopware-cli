package account

import (
	account_api "github.com/FriendsOfShopware/shopware-cli/account-api"
	"github.com/FriendsOfShopware/shopware-cli/config"
	"github.com/spf13/cobra"
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

func Register(rootCmd *cobra.Command, onInit func(commandName string) (*ServiceContainer, error)) {
	accountRootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		ser, err := onInit(cmd.Name())
		services = ser
		return err
	}
	rootCmd.AddCommand(accountRootCmd)
}
