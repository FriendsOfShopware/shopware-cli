package cmd

import (
	"fmt"
	"os"
	accountApi "shopware-cli/account-api"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into your Shopware Account",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		email := viper.GetString(ConfigAccountUser)
		password := viper.GetString(ConfigAccountPassword)
		newCredentials := false

		if len(email) == 0 || len(password) == 0 {
			email, password = askUserForEmailAndPassword()
			newCredentials = true

			viper.Set(ConfigAccountUser, email)
			viper.Set(ConfigAccountPassword, password)
		} else {
			log.Infof("Using existing credentials. Use account:logout to logout")
		}

		client, err := accountApi.NewApi(accountApi.LoginRequest{Email: email, Password: password})

		if err != nil {
			return errors.Wrap(err, "login failed with error")
		}

		if viper.GetInt(ConfigAccountCompany) > 0 {
			err = changeAPIMembership(client, viper.GetInt(ConfigAccountCompany))

			if err != nil {
				return errors.Wrap(err, "cannot change company member ship")
			}
		}

		if newCredentials {
			err := saveConfig()

			if err != nil {
				return errors.Wrap(err, "cannot save config")
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
