package project

import (
	"os/exec"

	"github.com/FriendsOfShopware/shopware-cli/extension"
	"github.com/FriendsOfShopware/shopware-cli/logging"
	"github.com/spf13/cobra"
)

var projectStorefrontBuildCmd = &cobra.Command{
	Use:   "storefront-build [path]",
	Short: "Builds the Storefront",
	RunE: func(cmd *cobra.Command, args []string) error {
		var projectRoot string
		var err error

		if len(args) == 1 {
			projectRoot = args[0]
		} else if projectRoot, err = findClosestShopwareProject(); err != nil {
			return err
		}

		logging.FromContext(cmd.Context()).Infof("Looking for extensions to build assets in project")

		sources := extension.FindAssetSourcesOfProject(cmd.Context(), args[0])
		constraint, err := extension.GetShopwareProjectConstraint(args[0])
		if err != nil {
			return err
		}

		assetCfg := extension.AssetBuildConfig{
			DisableAdminBuild: true,
			ShopwareRoot:      projectRoot,
			ShopwareVersion:   constraint,
		}

		if err := extension.BuildAssetsForExtensions(cmd.Context(), sources, assetCfg); err != nil {
			return err
		}

		return runTransparentCommand(commandWithRoot(exec.CommandContext(cmd.Context(), "php", "bin/console", "theme:compile"), projectRoot))
	},
}

func init() {
	projectRootCmd.AddCommand(projectStorefrontBuildCmd)
}
