package cmd

import (
	termColor "github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	accountApi "shopware-cli/account-api"
	"strings"
)

var accountCompanyProducerExtensionCreateCmd = &cobra.Command{
	Use:   "create [name] [classic|platform|themes|apps]",
	Short: "Creates a new extension",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getAccountApiByConfig()

		p, err := client.Producer()

		if err != nil {
			termColor.Red(err.Error())
			os.Exit(1)
		}

		profile, err := p.Profile()
		if err != nil {
			termColor.Red(err.Error())
			os.Exit(1)
		}

		if args[1] != accountApi.GenerationApps && args[1] != accountApi.GenerationPlatform && args[1] != accountApi.GenerationThemes && args[1] != accountApi.GenerationClassic {
			termColor.Red("Generation must be one of these options: %s %s %s %s", accountApi.GenerationPlatform, accountApi.GenerationThemes, accountApi.GenerationClassic, accountApi.GenerationApps)
			os.Exit(1)
		}

		if !strings.HasPrefix(args[0], profile.Prefix) {
			termColor.Red("Extension name must start with the prefix %s", profile.Prefix)
			os.Exit(1)
		}

		extension, err := p.CreateExtension(accountApi.CreateExtensionRequest{
			Name: args[0],
			Generation: struct {
				Name string `json:"name"`
			}{Name: args[1]},
			ProducerId: p.GetId(),
		})

		if err != nil {
			termColor.Red(err.Error())
			os.Exit(1)
		}

		termColor.Green("Extension with name %s has been successfully created", extension.Name)
	},
}

func init() {
	accountCompanyProducerExtensionCmd.AddCommand(accountCompanyProducerExtensionCreateCmd)
}
