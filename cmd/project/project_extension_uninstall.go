package project

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"shopware-cli/shop"
)

var projectExtensionUninstallCmd = &cobra.Command{
	Use:   "uninstall [name]",
	Short: "Uninstall a extension",
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

			if extension.InstalledAt == nil {
				log.Infof("Extension %s is already uninstalled", arg)
				continue
			}

			if extension.Active {
				if err := client.DeactivateExtension(cmd.Context(), extension.Type, extension.Name); err != nil {
					failed = true

					log.Errorf("Deactivation of %s failed with error: %v", extension.Name, err)
				} else {
					log.Infof("Deactivated %s", extension.Name)
				}
			}

			if err := client.UninstallExtension(cmd.Context(), extension.Type, extension.Name); err != nil {
				failed = true

				log.Errorf("Installation of %s failed with error: %v", extension.Name, err)
			}

			log.Infof("Uninstalled %s", extension.Name)
		}

		if failed {
			return fmt.Errorf("uninstall failed")
		}

		return nil
	},
}

func init() {
	projectExtensionCmd.AddCommand(projectExtensionUninstallCmd)
}
