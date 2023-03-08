package account

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/FriendsOfShopware/shopware-cli/logging"
)

var accountCompanyProducerExtensionDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a extension",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		extensionId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("cannot convert id to int: %w", err)
		}

		p, err := services.AccountClient.Producer(cmd.Context())

		if err != nil {
			return fmt.Errorf("cannot get producer endpoint: %w", err)
		}

		err = p.DeleteExtension(cmd.Context(), extensionId)

		if err != nil {
			return fmt.Errorf("cannot delete extension: %w", err)
		}

		logging.FromContext(cmd.Context()).Infof("Extension has been successfully deleted")

		return nil
	},
}

func init() {
	accountCompanyProducerExtensionCmd.AddCommand(accountCompanyProducerExtensionDeleteCmd)
}
