package cmd

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
	RunE: func(cmd *cobra.Command, args []string) error {
		err := accountApi.InvalidateTokenCache()
		if err != nil {
			return errors.Wrap(err, "cannot invalidate token cache")
		}

		appConfig.Account.Company = 0
		appConfig.Account.Email = ""
		appConfig.Account.Password = ""
		err = saveApplicationConfig()

		if err != nil {
			return errors.Wrap(err, "cannot write config")
		}

		log.Infof("You have been logged out")

		return nil
	},
}

func init() {
	accountRootCmd.AddCommand(logoutCmd)
}
