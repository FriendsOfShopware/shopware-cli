package cmd

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strconv"

	"github.com/spf13/cobra"
)

var accountCompanyProducerExtensionDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a extension",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getAccountAPIByConfigOrFail()

		extensionId, err := strconv.Atoi(args[0])

		if err != nil {
			return errors.Wrap(err, "cannot convert id to int")
		}

		p, err := client.Producer()

		if err != nil {
			return errors.Wrap(err, "cannot get producer endpoint")
		}

		err = p.DeleteExtension(extensionId)

		if err != nil {
			return errors.Wrap(err, "cannot delete extension")
		}

		log.Infof("Extension has been successfully deleted")

		return nil
	},
}

func init() {
	accountCompanyProducerExtensionCmd.AddCommand(accountCompanyProducerExtensionDeleteCmd)
}
