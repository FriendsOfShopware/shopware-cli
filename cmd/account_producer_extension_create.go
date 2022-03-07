package cmd

import (
	"fmt"
	accountApi "shopware-cli/account-api"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
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
	RunE: func(_ *cobra.Command, args []string) error {
		client := getAccountAPIByConfigOrFail()

		p, err := client.Producer()

		if err != nil {
			return errors.Wrap(err, "cannot get producer endpoint")
		}

		profile, err := p.Profile()
		if err != nil {
			return errors.Wrap(err, "cannot get producer profile")
		}

		if args[1] != accountApi.GenerationApps && args[1] != accountApi.GenerationPlatform && args[1] != accountApi.GenerationThemes && args[1] != accountApi.GenerationClassic {
			return fmt.Errorf("generation must be one of these options: %s %s %s %s", accountApi.GenerationPlatform, accountApi.GenerationThemes, accountApi.GenerationClassic, accountApi.GenerationApps)
		}

		if !strings.HasPrefix(args[0], profile.Prefix) {
			return fmt.Errorf("extension name must start with the prefix %s", profile.Prefix)
		}

		extension, err := p.CreateExtension(accountApi.CreateExtensionRequest{
			Name: args[0],
			Generation: struct {
				Name string `json:"name"`
			}{Name: args[1]},
			ProducerID: p.GetId(),
		})

		if err != nil {
			return errors.Wrap(err, "cannot create extension")
		}

		log.Infof("Extension with name %s has been successfully created", extension.Name)

		return nil
	},
}

func init() {
	accountCompanyProducerExtensionCmd.AddCommand(accountCompanyProducerExtensionCreateCmd)
}
