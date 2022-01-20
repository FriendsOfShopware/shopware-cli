package cmd

import (
	"context"
	"fmt"
	accountApi "shopware-cli/account-api"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfgFile string
var version = "dev"
var ctx context.Context

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
	ctx = context.Background()
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

	if err := initApplicationConfig(); err != nil {
		log.Fatalln(err)
	}
}

func getAccountAPIByConfig() (*accountApi.Client, error) {
	if appConfig.Account.Email == "" {
		return nil, fmt.Errorf("please login first using shopware-cli account login")
	}

	client, err := accountApi.NewApi(accountApi.LoginRequest{Email: appConfig.Account.Email, Password: appConfig.Account.Password})

	if err != nil {
		return nil, err
	}

	if appConfig.Account.Company > 0 {
		err = changeAPIMembership(client, appConfig.Account.Company)

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
