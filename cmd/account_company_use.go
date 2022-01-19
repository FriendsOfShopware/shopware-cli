package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"

	"github.com/spf13/cobra"
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
				appConfig.Account.Company = companyID

				err := saveApplicationConfig()
				if err != nil {
					return err
				}

				log.Infof("Successfully changed your company to %s (%d)", membership.Company.Name, membership.Company.CustomerNumber)
				return nil
			}
		}

		return fmt.Errorf("company with ID \"%d\" not found", companyID)
	},
}

func init() {
	accountCompanyRootCmd.AddCommand(accountCompanyUseCmd)
}
