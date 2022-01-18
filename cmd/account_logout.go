package cmd

import (
	"log"
	accountApi "shopware-cli/account-api"

	termColor "github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Shopware Account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := accountApi.InvalidateTokenCache()
		if err != nil {
			log.Fatalln(err)
		}

		viper.Set(ConfigAccountUser, "")
		viper.Set(ConfigAccountPassword, "")
		viper.Set(ConfigAccountCompany, "")

		err = viper.WriteConfig()

		if err != nil {
			log.Fatalln(err)
		}

		termColor.Green("You have been logged out")
	},
}

func init() {
	accountRootCmd.AddCommand(logoutCmd)
}
