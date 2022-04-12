package account

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var accountProducerInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "List information about your producer account",
	Long:  ``,
	RunE: func(_ *cobra.Command, _ []string) error {
		p, err := services.AccountClient.Producer()

		if err != nil {
			return errors.Wrap(err, "cannot get producer endpoint")
		}

		profile, err := p.Profile()
		if err != nil {
			return errors.Wrap(err, "cannot get producer profile")
		}

		log.Infof("Name: %s", profile.Name)
		log.Infof("Prefix: %s", profile.Prefix)
		log.Infof("Website: %s", profile.Website)

		return nil
	},
}

func init() {
	accountCompanyProducerCmd.AddCommand(accountProducerInfoCmd)
}
