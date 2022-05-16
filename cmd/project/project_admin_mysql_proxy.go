package project

import (
	"github.com/spf13/cobra"
)

var projectAdminMysqlCmd = &cobra.Command{
	Use:   "admin-api [method] [path]",
	Short: "pre authenticated curl interface to the Admin API",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	projectAdminMysqlCmd.AddCommand(projectAdminApiCmd)
}
