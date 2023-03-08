package account

import (
	"fmt"
	"strings"

	accountApi "github.com/FriendsOfShopware/shopware-cli/account-api"
	"github.com/FriendsOfShopware/shopware-cli/logging"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"

	accountApi "github.com/FriendsOfShopware/shopware-cli/account-api"
)

var accountCompanyProducerExtensionCreateCmd = &cobra.Command{
	Use:   "create [name] [classic|platform|themes|apps]",
	Short: "Creates a new extension",
	Args:  cobra.ExactArgs(2),
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			return []string{accountApi.GenerationApps, accountApi.GenerationClassic, accountApi.GenerationThemes, accountApi.GenerationPlatform}, cobra.ShellCompDirectiveNoFileComp
		}

		return []string{}, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := services.AccountClient.Producer(cmd.Context())

		if err != nil {
			return fmt.Errorf("cannot get producer endpoint: %w", err)
		}

		profile, err := p.Profile(cmd.Context())
		if err != nil {
			return fmt.Errorf("cannot get producer profile: %w", err)
		}

		if args[1] != accountApi.GenerationApps && args[1] != accountApi.GenerationPlatform && args[1] != accountApi.GenerationThemes && args[1] != accountApi.GenerationClassic {
			return fmt.Errorf("generation must be one of these options: %s %s %s %s", accountApi.GenerationPlatform, accountApi.GenerationThemes, accountApi.GenerationClassic, accountApi.GenerationApps)
		}

		if !strings.HasPrefix(args[0], profile.Prefix) {
			return fmt.Errorf("extension name must start with the prefix %s", profile.Prefix)
		}

		extension, err := p.CreateExtension(cmd.Context(), accountApi.CreateExtensionRequest{
			Name: args[0],
			Generation: struct {
				Name string `json:"name"`
			}{Name: args[1]},
			ProducerID: p.GetId(),
		})
		if err != nil {
			return fmt.Errorf("cannot create extension: %w", err)
		}

		logging.FromContext(cmd.Context()).Infof("Extension with name %s has been successfully created", extension.Name)

		return nil
	},
}

func init() {
	accountCompanyProducerExtensionCmd.AddCommand(accountCompanyProducerExtensionCreateCmd)
}
