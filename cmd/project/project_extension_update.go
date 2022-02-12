package project

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"shopware-cli/shop"
)

var projectExtensionUpdateCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update a extension",
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

		disableStoreUpdates, _ := cmd.PersistentFlags().GetBool("disable-store-update")

		if err := client.RefreshExtensions(cmd.Context()); err != nil {
			return err
		}

		extensions, err := client.GetAvailableExtensions(cmd.Context())

		if err != nil {
			return err
		}

		failed := false

		if len(args) == 1 && args[0] == "all" {
			args = make([]string, 0)

			for _, extension := range extensions {
				args = append(args, extension.Name)
			}
		}

		for _, arg := range args {
			extension := extensions.GetByName(arg)

			if extension == nil {
				failed = true
				log.Errorf("Cannot find extension by name %s", arg)
				continue
			}

			if extension.LatestVersion == "" || extension.Version == extension.LatestVersion {
				log.Infof("Extension %s is up to date", arg)
				continue
			}

			if !extension.Active {
				log.Infof("Extension %s is not active skipping", arg)
				continue
			}

			if extension.UpdateSource == "store" && !disableStoreUpdates {
				if err := client.DownloadExtension(cmd.Context(), arg); err != nil {
					log.Errorf("Download of %s update failed with error: %v", extension.Name, err)
					failed = true
					continue
				}
			}

			if err := client.UpdateExtension(cmd.Context(), extension.Type, extension.Name); err != nil {
				failed = true

				log.Errorf("Update of %s failed with error: %v", extension.Name, err)
			}

			log.Infof("Updated %s", extension.Name)
		}

		if failed {
			return fmt.Errorf("update failed")
		}

		return nil
	},
}

func init() {
	projectExtensionCmd.AddCommand(projectExtensionUpdateCmd)
	projectExtensionUpdateCmd.PersistentFlags().Bool("disable-store-update", false, "Downloads updates from store.shopware.com")
}
