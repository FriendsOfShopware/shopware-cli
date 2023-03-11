package project

import (
	"github.com/spf13/cobra"
)

var projectStorefrontBuildCmd = &cobra.Command{
	Use:   "storefront-build",
	Short: "Builds the Storefront",
	RunE: func(cmd *cobra.Command, _ []string) error {
		forceNpmInstall, _ := cmd.PersistentFlags().GetBool("npm-install")

		var projectRoot string
		var err error

		if projectRoot, err = findClosestShopwareProject(); err != nil {
			return err
		}

		return buildStorefront(cmd.Context(), projectRoot, forceNpmInstall)
	},
}

func init() {
	projectRootCmd.AddCommand(projectStorefrontBuildCmd)
	projectStorefrontBuildCmd.PersistentFlags().Bool("npm-install", false, "Run npm install")
}
