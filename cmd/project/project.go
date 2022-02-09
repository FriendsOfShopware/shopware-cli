package project

import (
	"github.com/spf13/cobra"
)

var projectConfigPath string

var projectRootCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage your Shopware Project",
}

func Register(rootCmd *cobra.Command) {
	rootCmd.AddCommand(projectRootCmd)
	projectRootCmd.PersistentFlags().StringVar(&projectConfigPath, "project-config", ".shopware-project.yml", "Path to .shopware-project.yml")
}
