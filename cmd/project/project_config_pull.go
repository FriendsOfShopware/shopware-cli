package project

import (
	"io/ioutil"
	"os"
	"shopware-cli/shop"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var projectConfigPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Synchronizes your shop config to local",
	RunE: func(cmd *cobra.Command, _ []string) error {
		var cfg *shop.Config
		var err error

		if cfg, err = shop.ReadConfig(projectConfigPath); err != nil {
			return err
		}

		client, err := shop.NewShopClient(cmd.Context(), cfg, nil)
		if err != nil {
			return err
		}

		if cfg.Sync == nil {
			cfg.Sync = &shop.ConfigSync{}
		}

		for _, applyer := range NewSyncApplyers() {
			if err := applyer.Pull(cmd.Context(), client, cfg); err != nil {
				return err
			}
		}

		content, err := yaml.Marshal(cfg)

		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(projectConfigPath, content, os.ModePerm); err != nil {
			return err
		}

		log.Infof("%s has been updated", projectConfigPath)

		return nil
	},
}

func init() {
	projectConfigCmd.AddCommand(projectConfigPullCmd)
}
