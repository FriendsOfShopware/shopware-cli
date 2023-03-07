package project

import (
	"fmt"

	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/FriendsOfShopware/shopware-cli/shop"
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

		client, err := shop.NewShopClient(cmd.Context(), cfg)
		if err != nil {
			return err
		}

		extensions, _, err := client.ExtensionManager.ListAvailableExtensions(adminSdk.NewApiContext(cmd.Context()))

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
				if _, err := client.ExtensionManager.DeactivateExtension(adminSdk.NewApiContext(cmd.Context()), extension.Type, extension.Name); err != nil {
					failed = true

					log.Errorf("Deactivation of %s failed with error: %v", extension.Name, err)
					continue
				}

				log.Infof("Deactivated %s", extension.Name)
			}

			if extension.InstalledAt != nil {
				if _, err := client.ExtensionManager.UninstallExtension(adminSdk.NewApiContext(cmd.Context()), extension.Type, extension.Name); err != nil {
					failed = true

					log.Errorf("Uninstall of %s failed with error: %v", extension.Name, err)
					continue
				}

				log.Infof("Uninstalled %s", extension.Name)
			}

			if _, err := client.ExtensionManager.RemoveExtension(adminSdk.NewApiContext(cmd.Context()), extension.Type, extension.Name); err != nil {
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
