package cmd

import (
	"fmt"
	termColor "github.com/fatih/color"
	"os"
	accountApi "shopware-cli/account-api"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ConfigAccountUser     = "account.email"
	ConfigAccountPassword = "account.password"
	ConfigAccountCompany  = "account.company"
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
	email := viper.GetString(ConfigAccountUser)
	password := viper.GetString(ConfigAccountPassword)

	client, err := accountApi.NewApi(accountApi.LoginRequest{Email: email, Password: password})

	if err != nil {
		termColor.Red("Login failed with error: %s", err.Error())
		os.Exit(1)
	}

	companyId := viper.GetInt(ConfigAccountCompany)

	err = changeApiMembership(client, companyId)

	if err != nil {
		termColor.Red(err.Error())
		os.Exit(1)
	}

	return client
}

func changeApiMembership(client *accountApi.Client, companyId int) error {
	if companyId == 0 || client.GetActiveCompanyId() == companyId {
		return nil
	}

	for _, membership := range *client.GetMemberships() {
		if membership.Company.Id == companyId {
			return client.ChangeActiveMembership(&membership)
		}
	}

	return fmt.Errorf("could not find configured company with id %d", companyId)
}

func saveConfig() error {
	err := viper.SafeWriteConfig()

	if err != nil {
		err = viper.WriteConfig()
	}

	return err
}
