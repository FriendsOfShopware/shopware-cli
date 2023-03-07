package extension

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/FriendsOfShopware/shopware-cli/extension"
)

var extensionAssetBundleCmd = &cobra.Command{
	Use:   "build [path]",
	Short: "Builds assets for extensions",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		assetCfg := extension.AssetBuildConfig{EnableESBuildForAdmin: false, EnableESBuildForStorefront: false}
		validatedExtensions := make([]extension.Extension, 0)

		for _, arg := range args {
			path, err := filepath.Abs(arg)
			if err != nil {
				return fmt.Errorf("cannot open file: %w", err)
			}

			ext, err := extension.GetExtensionByFolder(path)
			if err != nil {
				return fmt.Errorf("cannot open extension: %w", err)
			}

			validatedExtensions = append(validatedExtensions, ext)
		}

		if len(args) == 1 {
			extCfg, err := extension.ReadExtensionConfig(validatedExtensions[0].GetPath())
			if err != nil {
				return fmt.Errorf("cannot read extension config: %w", err)
			}

			assetCfg.EnableESBuildForAdmin = extCfg.Build.Zip.Assets.EnableESBuildForAdmin
			assetCfg.EnableESBuildForStorefront = extCfg.Build.Zip.Assets.EnableESBuildForStorefront
		}

		err := extension.BuildAssetsForExtensions(os.Getenv("SHOPWARE_PROJECT_ROOT"), validatedExtensions, assetCfg)
		if err != nil {
			return fmt.Errorf("cannot build assets: %w", err)
		}

		log.Infof("Assets has been built")

		return nil
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionAssetBundleCmd)
}
