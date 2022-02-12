package project

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"shopware-cli/shop"
)

var projectExtensionDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a extension",
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

		extensions, err := client.GetAvailableExtensions(cmd.Context())

		if err != nil {
			return err
		}

		failed := false

		for _, arg := range args {
			extension := extensions.GetByName(arg)

			if extension == nil {
				log.Errorf("Cannot find extension by name %s", arg)
				continue
			}

			if extension.Active {
				if err := client.DeactivateExtension(cmd.Context(), extension.Type, extension.Name); err != nil {
					failed = true

					log.Errorf("Deactivation of %s failed with error: %v", extension.Name, err)
				}

				continue
			}

			if extension.InstalledAt != nil {
				if err := client.UninstallExtension(cmd.Context(), extension.Type, extension.Name); err != nil {
					failed = true

					log.Errorf("Uninstall of %s failed with error: %v", extension.Name, err)
				}
			}

			if err := client.RemoveExtension(cmd.Context(), extension.Type, extension.Name); err != nil {
				failed = true

				log.Errorf("Remove of %s failed with error: %v", extension.Name, err)
			}

			log.Infof("Removed %s", extension.Name)
		}

		if failed {
			return fmt.Errorf("remove failed")
		}

		return nil
	},
}

func init() {
	projectExtensionCmd.AddCommand(projectExtensionDeleteCmd)
}
