package account

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	accountApi "shopware-cli/account-api"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Shopware Account",
	Long:  ``,
	RunE: func(_ *cobra.Command, _ []string) error {
		err := accountApi.InvalidateTokenCache()
		if err != nil {
			return errors.Wrap(err, "cannot invalidate token cache")
		}

		_ = services.Conf.SetAccountCompanyId(0)
		_ = services.Conf.SetAccountEmail("")
		_ = services.Conf.SetAccountPassword("")

		if err := services.Conf.Save(); err != nil {
			return errors.Wrap(err, "cannot write config")
		}

		log.Infof("You have been logged out")

		return nil
	},
}

func init() {
	accountRootCmd.AddCommand(logoutCmd)
}
