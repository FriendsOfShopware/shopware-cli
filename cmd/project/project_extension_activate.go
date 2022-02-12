package project

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"shopware-cli/shop"
)

var projectExtensionActivateCmd = &cobra.Command{
	Use:   "activate [name]",
	Short: "Activate a extension",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg *shop.Config
		var err error

		if cfg, err = shop.ReadConfig(projectConfigPath); err != nil {
			return err
		}

		client, err := shop.NewShopClient(cmd.Context(), cfg, nil)
		if err != nil {
			return err
		}

		extensions, err := client.GetInstalledExtensions(cmd.Context())

		if err != nil {
			return err
		}

		failed := false

		for _, arg := range args {
			extension := extensions.GetByName(arg)

			if extension == nil {
				failed = true
				log.Errorf("Cannot find extension by name %s", arg)
				continue
			}

			if extension.Active {
				log.Infof("Extension %s is already active", arg)
				continue
			}

			if extension.InstalledAt == nil {
				if err := client.InstallExtension(cmd.Context(), extension.Type, extension.Name); err != nil {
					failed = true

					log.Errorf("Installation of %s failed with error: %v", extension.Name, err)
				}
			}

			if err := client.ActivateExtension(cmd.Context(), extension.Type, extension.Name); err != nil {
				failed = true

				log.Errorf("Activate of %s failed with error: %v", extension.Name, err)
			}

			log.Infof("Activated %s", extension.Name)
		}

		if failed {
			return fmt.Errorf("activation failed")
		}

		return nil
	},
}

func init() {
	projectExtensionCmd.AddCommand(projectExtensionActivateCmd)
}
