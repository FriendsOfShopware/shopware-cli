package cmd

import (
	"fmt"
	"os"
	accountApi "shopware-cli/account-api"

	log "github.com/sirupsen/logrus"
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
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.SilenceErrors = true

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.shopware-cli.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "show debug output")
}

func initConfig() {
	if verbose, _ := rootCmd.PersistentFlags().GetBool("verbose"); verbose {
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			log.Println(err)
		}

		// Search config in home directory with name ".shopware-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".shopware-cli")
	}

	viper.AutomaticEnv()
	_ = viper.ReadInConfig()
}

func getAccountAPIByConfig() (*accountApi.Client, error) {
	email := viper.GetString(ConfigAccountUser)
	password := viper.GetString(ConfigAccountPassword)

	client, err := accountApi.NewApi(accountApi.LoginRequest{Email: email, Password: password})

	if err != nil {
		return nil, err
	}

	companyID := viper.GetInt(ConfigAccountCompany)

	if companyID > 0 {
		err = changeAPIMembership(client, companyID)

		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func getAccountAPIByConfigOrFail() *accountApi.Client {
	client, err := getAccountAPIByConfig()

	if err != nil {
		log.Fatalln(err)
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
