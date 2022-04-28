package cmd

import (
	"context"
	accountApi "shopware-cli/account-api"
	"shopware-cli/cmd/account"
	"shopware-cli/cmd/extension"
	"shopware-cli/cmd/project"
	"shopware-cli/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfgFile string
var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "shopware-cli",
	Short:   "A cli for common Shopware tasks",
	Long:    `This application contains some utilities like extension management`,
	Version: version,
}

func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.SilenceErrors = true

	cobra.OnInitialize(func() {
		_ = config.InitConfig(cfgFile)
	})

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.shopware-cli.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "show debug output")

	project.Register(rootCmd)
	extension.Register(rootCmd)
	account.Register(rootCmd, func() (*account.ServiceContainer, error) {
		err := config.InitConfig(cfgFile)
		if err != nil {
			return nil, err
		}
		conf := config.Config{}
		client, err := accountApi.NewApi(conf)
		if err != nil {
			return nil, err
		}
		return &account.ServiceContainer{
			Conf:          conf,
			AccountClient: client,
		}, nil
	})
}
