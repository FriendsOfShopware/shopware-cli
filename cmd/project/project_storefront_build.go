package project

import (
	"path/filepath"

	"github.com/FriendsOfShopware/shopware-cli/internal/phpexec"
	"github.com/FriendsOfShopware/shopware-cli/shop"

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
			// We need an absolute path for webpack
			projectRoot, err = filepath.Abs(args[0])
			if err != nil {
				return err
			}
		} else if projectRoot, err = findClosestShopwareProject(); err != nil {
			return err
		}

		shopCfg, err := shop.ReadConfig(projectConfigPath, true)
		if err != nil {
			return err
		}

		logging.FromContext(cmd.Context()).Infof("Looking for extensions to build assets in project")

		sources := extension.FindAssetSourcesOfProject(cmd.Context(), projectRoot, shopCfg)

		forceInstall, _ := cmd.PersistentFlags().GetBool("force-install-dependencies")

		shopwareConstraint, err := extension.GetShopwareProjectConstraint(projectRoot)
		if err != nil {
			return err
		}

		assetCfg := extension.AssetBuildConfig{
			DisableAdminBuild: true,
			ShopwareRoot:      projectRoot,
			ShopwareVersion:   shopwareConstraint,
			NPMForceInstall:   forceInstall,
			ContributeProject: extension.IsContributeProject(projectRoot),
		}

		if err := extension.BuildAssetsForExtensions(cmd.Context(), sources, assetCfg); err != nil {
			return err
		}

		skipThemeCompile, _ := cmd.PersistentFlags().GetBool("skip-theme-compile")
		if skipThemeCompile {
			return nil
		}

		return runTransparentCommand(commandWithRoot(phpexec.ConsoleCommand(cmd.Context(), "theme:compile"), projectRoot))
	},
}

func init() {
	projectRootCmd.AddCommand(projectStorefrontBuildCmd)
	projectStorefrontBuildCmd.PersistentFlags().Bool("skip-theme-compile", false, "Skip theme compilation")
	projectStorefrontBuildCmd.PersistentFlags().Bool("force-install-dependencies", false, "Force install NPM dependencies")
}
