package account

import (
	"errors"
	"fmt"
	"os"

	accountApi "github.com/FriendsOfShopware/shopware-cli/account-api"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into your Shopware Account",
	Long:  "",
	RunE: func(_ *cobra.Command, _ []string) error {
		email := services.Conf.GetAccountEmail()
		password := services.Conf.GetAccountPassword()
		newCredentials := false

		if len(email) == 0 || len(password) == 0 {
			email, password = askUserForEmailAndPassword()
			newCredentials = true

			if err := services.Conf.SetAccountEmail(email); err != nil {
				return err
			}
			if err := services.Conf.SetAccountPassword(password); err != nil {
				return err
			}
		} else {
			log.Infof("Using existing credentials. Use account:logout to logout")
		}

		client, err := accountApi.NewApi(accountApi.LoginRequest{Email: email, Password: password})
		if err != nil {
			return fmt.Errorf("login failed with error: %w", err)
		}

		if companyId := services.Conf.GetAccountCompanyId(); companyId > 0 {
			err = changeAPIMembership(client, companyId)

			if err != nil {
				return fmt.Errorf("cannot change company member ship: %w", err)
			}
		}

		if newCredentials {
			err := services.Conf.Save()
			if err != nil {
				return fmt.Errorf("cannot save config: %w", err)
			}
		}

		profile, err := client.GetMyProfile()
		if err != nil {
			return err
		}

		log.Infof(
			"Hey %s %s. You are now authenticated on company %s and can use all account commands",
			profile.PersonalData.FirstName,
			profile.PersonalData.LastName,
			client.GetActiveMembership().Company.Name,
		)

		return nil
	},
}

func init() {
	accountRootCmd.AddCommand(loginCmd)
}

func askUserForEmailAndPassword() (string, string) {
	emailPrompt := promptui.Prompt{
		Label:    "Email",
		Validate: emptyValidator,
	}

	email, err := emailPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	passwordPrompt := promptui.Prompt{
		Label:    "Password",
		Validate: emptyValidator,
		Mask:     '*',
	}

	password, err := passwordPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return email, password
}

func emptyValidator(s string) error {
	if len(s) == 0 {
		return errors.New("this cannot be empty")
	}

	return nil
}

func changeAPIMembership(client *accountApi.Client, companyID int) error {
	if companyID == 0 || client.GetActiveCompanyID() == companyID {
		log.Tracef("Client is on correct membership skip")
		return nil
	}

	for _, membership := range client.GetMemberships() {
		if membership.Company.Id == companyID {
			log.Tracef("Changing member ship from %s (%d) to %s (%d)", client.ActiveMembership.Company.Name, client.ActiveMembership.Company.Id, membership.Company.Name, membership.Company.Id)
			return client.ChangeActiveMembership(membership)
		}
	}

	return fmt.Errorf("could not find configured company with id %d", companyID)
}
