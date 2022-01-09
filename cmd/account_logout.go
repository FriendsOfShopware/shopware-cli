package cmd

import (
	termColor "github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "account:logout",
	Short: "Logout from Shopware Account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("account_email", "")
		viper.Set("account_password", "")

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
