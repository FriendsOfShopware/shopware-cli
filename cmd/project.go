package cmd

import (
	"github.com/spf13/cobra"
)

var projectRootCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage your Shopware Project",
}

func init() {
	rootCmd.AddCommand(projectRootCmd)
}
