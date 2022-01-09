package cmd

import (
	termColor "github.com/fatih/color"
	"os"
	accountApi "shopware-cli/account-api"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "shopware-cli",
	Short: "A cli for common Shopware tasks",
	Long:  `This application contains some utilities like extension management`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.shopware-cli.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".shopware-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".shopware-cli")
	}

	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
	}
}

func getAccountApiByConfig() *accountApi.Client {
	email := viper.GetString("account_email")
	password := viper.GetString("account_password")

	client, err := accountApi.NewApi(accountApi.LoginRequest{Email: email, Password: password})

	if err != nil {
		termColor.Red("Login failed with error: %s", err.Error())
		os.Exit(1)
	}

	return client
}
