package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var projectStorefrontBuildCmd = &cobra.Command{
	Use:   "storefront-build",
	Short: "Builds the Storefront",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		var projectRoot string
		var err error

		if projectRoot, err = findClosestShopwareProject(); err != nil {
			return err
		}

		adminRoot := getPlatformPath("Storefront", "Resources/app/storefront")

		if err := runSimpleCommand(projectRoot, "php", "bin/console", "bundle:dump"); err != nil {
			return err
		}

		// Optional command, allowed to failure
		_ = runSimpleCommand(projectRoot, "php", "bin/console", "feature:dump")

		// Optional npm install

		_, err = os.Stat(getPlatformPath("Storefront", "Resources/app/storefront/node_modules"))

		forceNpmInstall, _ := cobraCmd.PersistentFlags().GetBool("npm-install")

		if forceNpmInstall || os.IsNotExist(err) {
			if installErr := runSimpleCommand(projectRoot, "npm", "install", "--prefix", adminRoot); err != nil {
				return installErr
			}
		}

		envs := []string{
			fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
			fmt.Sprintf("PROJECT_ROOT=%s", projectRoot),
		}

		npmRun := exec.Command("npm", "--prefix", adminRoot, "run", "production")
		npmRun.Env = envs
		npmRun.Stdin = os.Stdin
		npmRun.Stdout = os.Stdout
		npmRun.Stderr = os.Stderr

		if err := npmRun.Run(); err != nil {
			return err
		}

		if err := runSimpleCommand(projectRoot, "php", "bin/console", "assets:install"); err != nil {
			return err
		}

		return runSimpleCommand(projectRoot, "php", "bin/console", "theme:compile")
	},
}

func init() {
	projectRootCmd.AddCommand(projectStorefrontBuildCmd)
	projectStorefrontBuildCmd.PersistentFlags().Bool("npm-install", false, "Run npm install")
}
