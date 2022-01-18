package cmd

import (
	"fmt"
	"strconv"

	termColor "github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var accountCompanyUseCmd = &cobra.Command{
	Use:   "use [companyId]",
	Short: "Use another company for your Account",
	Args:  cobra.MinimumNArgs(1),
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		companyID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		client := getAccountAPIByConfigOrFail()

		for _, membership := range client.GetMemberships() {
			if membership.Company.Id == companyID {
				viper.Set(ConfigAccountCompany, companyID)

				err := saveConfig()
				if err != nil {
					return err
				}

				termColor.Green("Successfully changed your company to %s (%d)", membership.Company.Name, membership.Company.CustomerNumber)
				return nil
			}
		}

		return fmt.Errorf("company with ID \"%d\" not found", companyID)
	},
}

func init() {
	accountCompanyRootCmd.AddCommand(accountCompanyUseCmd)
}
