package project

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"shopware-cli/shop"
)

var projectClearCacheCmd = &cobra.Command{
	Use:   "clear-cache",
	Short: "Clears the Shop cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg *shop.Config
		var err error

		if cfg, err = shop.ReadConfig(projectConfigPath); err != nil {
			return err
		}

		if cfg.AdminApi == nil {
			log.Infof("Clearing cache localy")

			projectRoot, err := findClosestShopwareProject()

			if err != nil {
				return err
			}

			return os.RemoveAll(fmt.Sprintf("%s/var/cache", projectRoot))
		}

		log.Infof("Clearing cache using admin-api")

		client, err := shop.NewShopClient(cmd.Context(), cfg, nil)
		if err != nil {
			return err
		}

		return client.ClearCache(cmd.Context())
	},
}

func init() {
	projectRootCmd.AddCommand(projectClearCacheCmd)
}
