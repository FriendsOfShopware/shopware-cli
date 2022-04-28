package account

import (
	"fmt"
	accountApi "shopware-cli/account-api"
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var accountCompanyUseCmd = &cobra.Command{
	Use:   "use [companyId]",
	Short: "Use another company for your Account",
	Args:  cobra.MinimumNArgs(1),
	Long:  ``,
	RunE: func(_ *cobra.Command, args []string) error {
		companyID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		for _, membership := range services.AccountClient.GetMemberships() {
			if membership.Company.Id == companyID {
				if err := services.Conf.SetAccountCompanyId(companyID); err != nil {
					return err
				}

				if err := services.Conf.Save(); err != nil {
					return err
				}

				err = accountApi.InvalidateTokenCache()
				if err != nil {
					return errors.Wrap(err, "cannot invalidate token cache")
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
