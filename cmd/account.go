package cmd

import (
	"github.com/spf13/cobra"
)

var accountRootCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage your Shopware Account",
}

func init() {
	rootCmd.AddCommand(accountRootCmd)
}
