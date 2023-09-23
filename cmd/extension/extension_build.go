package extension

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/FriendsOfShopware/shopware-cli/extension"
	"github.com/FriendsOfShopware/shopware-cli/logging"
)

var extensionAssetBundleCmd = &cobra.Command{
	Use:   "build [path]",
	Short: "Builds assets for extensions",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		assetCfg := extension.AssetBuildConfig{
			EnableESBuildForAdmin:      false,
			EnableESBuildForStorefront: false,
			ShopwareRoot:               os.Getenv("SHOPWARE_PROJECT_ROOT"),
		}
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
			extCfg := validatedExtensions[0].GetExtensionConfig()

			assetCfg.EnableESBuildForAdmin = extCfg.Build.Zip.Assets.EnableESBuildForAdmin
			assetCfg.EnableESBuildForStorefront = extCfg.Build.Zip.Assets.EnableESBuildForStorefront
		}

		constraint, err := validatedExtensions[0].GetShopwareVersionConstraint()
		if err != nil {
			return fmt.Errorf("cannot get shopware version constraint: %w", err)
		}

		assetCfg.ShopwareVersion = constraint

		err = extension.BuildAssetsForExtensions(cmd.Context(), extension.ConvertExtensionsToSources(cmd.Context(), validatedExtensions), assetCfg)
		if err != nil {
			return fmt.Errorf("cannot build assets: %w", err)
		}

		logging.FromContext(cmd.Context()).Infof("Assets has been built")

		return nil
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionAssetBundleCmd)
}
