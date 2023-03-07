package account

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var accountCompanyProducerExtensionDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a extension",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		extensionId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("cannot convert id to int: %w", err)
		}

		p, err := services.AccountClient.Producer()
		if err != nil {
			return fmt.Errorf("cannot get producer endpoint: %w", err)
		}

		err = p.DeleteExtension(extensionId)

		if err != nil {
			return fmt.Errorf("cannot delete extension: %w", err)
		}

		log.Infof("Extension has been successfully deleted")

		return nil
	},
}

func init() {
	accountCompanyProducerExtensionCmd.AddCommand(accountCompanyProducerExtensionDeleteCmd)
}
