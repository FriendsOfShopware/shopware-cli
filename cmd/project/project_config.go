package project

import (
	"github.com/spf13/cobra"
)

var projectConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the project config",
}

func init() {
	projectRootCmd.AddCommand(projectConfigCmd)
}
