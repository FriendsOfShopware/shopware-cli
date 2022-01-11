package cmd

import (
	termColor "github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strconv"
)

var accountCompanyUseCmd = &cobra.Command{
	Use:   "use [companyId]",
	Short: "Use another company for your Account",
	Args:  cobra.MinimumNArgs(1),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		companyId, err := strconv.Atoi(args[0])

		if err != nil {
			log.Fatalln(err)
		}

		client := getAccountApiByConfig()

		for _, membership := range *client.GetMemberships() {
			if membership.Company.Id == companyId {
				viper.Set(ConfigAccountCompany, companyId)
				err := saveConfig()

				if err != nil {
					log.Fatalln(err)
				}

				termColor.Green("Successfully changed your company to %s (%d)", membership.Company.Name, membership.Company.CustomerNumber)
				return
			}
		}

		termColor.Red("Could noy find company by id %s", companyId)
	},
}

func init() {
	accountCompanyRootCmd.AddCommand(accountCompanyUseCmd)
}
