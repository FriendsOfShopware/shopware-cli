package cmd

import (
	"context"

	"github.com/spf13/cobra"

	accountApi "github.com/FriendsOfShopware/shopware-cli/account-api"
	"github.com/FriendsOfShopware/shopware-cli/cmd/account"
	"github.com/FriendsOfShopware/shopware-cli/cmd/extension"
	"github.com/FriendsOfShopware/shopware-cli/cmd/project"
	"github.com/FriendsOfShopware/shopware-cli/config"
	"github.com/FriendsOfShopware/shopware-cli/logging"
)

var (
	cfgFile string
	version = "dev"
)

var rootCmd = &cobra.Command{
	Use:     "shopware-cli",
	Short:   "A cli for common Shopware tasks",
	Long:    `This application contains some utilities like extension management`,
	Version: version,
}

func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		logging.FromContext(ctx).Fatalln(err)
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
	account.Register(rootCmd, func(commandName string) (*account.ServiceContainer, error) {
		err := config.InitConfig(cfgFile)
		if err != nil {
			return nil, err
		}
		conf := config.Config{}
		if commandName == "login" || commandName == "logout" {
			return &account.ServiceContainer{
				Conf:          conf,
				AccountClient: nil,
			}, nil
		}
		client, err := accountApi.NewApi(rootCmd.Context(), conf)
		if err != nil {
			return nil, err
		}
		return &account.ServiceContainer{
			Conf:          conf,
			AccountClient: client,
		}, nil
	})
}
