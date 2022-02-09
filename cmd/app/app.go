package app

import (
	"github.com/spf13/cobra"
)

var appDir string

var appRootCommand = &cobra.Command{
	Use:   "app",
	Short: "Helpers for app development",
}

func Register(rootCmd *cobra.Command) {
	pushCommand.PersistentFlags().StringVar(&appDir, "dir", "./", "directory of the app")
	rootCmd.AddCommand(appRootCommand)
}
