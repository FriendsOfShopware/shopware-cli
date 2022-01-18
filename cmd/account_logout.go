package cmd

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	accountApi "shopware-cli/account-api"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		viper.Set(ConfigAccountUser, "")
		viper.Set(ConfigAccountPassword, "")
		viper.Set(ConfigAccountCompany, "")

		err = viper.WriteConfig()

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
