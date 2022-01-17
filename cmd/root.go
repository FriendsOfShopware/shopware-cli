package cmd

import (
	"fmt"
	"os"
	accountApi "shopware-cli/account-api"

	termColor "github.com/fatih/color"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ConfigAccountUser     = "account.email"
	ConfigAccountPassword = "account.password"
	ConfigAccountCompany  = "account.company"
)

var cfgFile string
var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "shopware-cli",
	Short:   "A cli for common Shopware tasks",
	Long:    `This application contains some utilities like extension management`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
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
	cobra.CheckErr(viper.ReadInConfig())
}

func getAccountAPIByConfig() *accountApi.Client {
	email := viper.GetString(ConfigAccountUser)
	password := viper.GetString(ConfigAccountPassword)

	client, err := accountApi.NewApi(accountApi.LoginRequest{Email: email, Password: password})

	if err != nil {
		termColor.Red("Login failed with error: %s", err.Error())
		os.Exit(1)
	}

	companyID := viper.GetInt(ConfigAccountCompany)

	err = changeAPIMembership(client, companyID)

	if err != nil {
		termColor.Red(err.Error())
		os.Exit(1)
	}

	return client
}

func changeAPIMembership(client *accountApi.Client, companyID int) error {
	if companyID == 0 || client.GetActiveCompanyID() == companyID {
		return nil
	}

	for _, membership := range client.GetMemberships() {
		if membership.Company.Id == companyID {
			return client.ChangeActiveMembership(membership)
		}
	}

	return fmt.Errorf("could not find configured company with id %d", companyID)
}

func saveConfig() error {
	err := viper.SafeWriteConfig()

	if err != nil {
		err = viper.WriteConfig()
	}

	return err
}
