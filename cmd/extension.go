package cmd

import (
	"github.com/spf13/cobra"
)

var extensionRootCmd = &cobra.Command{
	Use:   "extension",
	Short: "Shopware Extension utilities",
}

func init() {
	rootCmd.AddCommand(extensionRootCmd)
}
