package cmd

import (
	termColor "github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var logoutCmd = &cobra.Command{
	Use:   "account:logout",
	Short: "Logout from Shopware Account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set(ConfigAccountUser, "")
		viper.Set(ConfigAccountPassword, "")
		viper.Set(ConfigAccountCompany, "")

		err := viper.WriteConfig()

		if err != nil {
			log.Fatalln(err)
		}

		termColor.Green("You have been logged out")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
