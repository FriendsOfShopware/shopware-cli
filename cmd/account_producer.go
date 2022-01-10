package cmd

import (
	"github.com/spf13/cobra"
)

var accountCompanyProducerCmd = &cobra.Command{
	Use:   "company",
	Short: "Manage your Shopware manufacturer",
}

func init() {
	accountRootCmd.AddCommand(accountCompanyProducerCmd)
}
