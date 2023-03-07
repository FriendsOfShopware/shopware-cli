package project

import (
	"errors"
	"fmt"
	"os"

	"github.com/FriendsOfShopware/shopware-cli/shop"

	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var projectConfigInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a new project config in current dir",
	RunE: func(_ *cobra.Command, _ []string) error {
		config := &shop.Config{}
		var content []byte
		var err error

		urlPrompt := promptui.Prompt{
			Label:    "Shop-URL example: http://localhost",
			Validate: emptyValidator,
		}

		config.URL, err = urlPrompt.Run()

		if err != nil {
			return err
		}

		if err = askApi(config); err != nil {
			return err
		}

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}

		if content, err = yaml.Marshal(config); err != nil {
			return err
		}

		if err := os.WriteFile(".shopware-project.yml", content, os.ModePerm); err != nil {
			return err
		}

		log.Info("Created .shopware-project.yml")

		return nil
	},
}

func askApi(config *shop.Config) error {
	adminApi := promptui.Prompt{
		Label:     "Configure admin-api access",
		IsConfirm: true,
	}

	var result string
	_, err := adminApi.Run()
	if err != nil {
		return nil //nolint:nilerr
	}

	apiAuthType := promptui.Select{
		Label: "Auth type",
		Items: []string{"user-password", "integration"},
	}

	if _, result, err = apiAuthType.Run(); err != nil {
		return err
	}

	apiConfig := shop.ConfigAdminApi{}
	config.AdminApi = &apiConfig

	if result == "integration" {
		clientIdPrompt := promptui.Prompt{
			Label:    "Client-ID",
			Validate: emptyValidator,
		}

		clientSecretPrompt := promptui.Prompt{
			Label:    "Client-Secret",
			Validate: emptyValidator,
		}

		if id, err := clientIdPrompt.Run(); err != nil {
			return err
		} else {
			apiConfig.ClientId = id
		}

		if secret, err := clientSecretPrompt.Run(); err != nil {
			return err
		} else {
			apiConfig.ClientSecret = secret
		}

		return nil
	}

	adminUserPrompt := promptui.Prompt{
		Label:    "Admin User",
		Validate: emptyValidator,
	}

	adminPasswordPrompt := promptui.Prompt{
		Label:    "Admin Password",
		Validate: emptyValidator,
	}

	if user, err := adminUserPrompt.Run(); err != nil {
		return err
	} else {
		apiConfig.Username = user
	}

	if password, err := adminPasswordPrompt.Run(); err != nil {
		return err
	} else {
		apiConfig.Password = password
	}

	return nil
}

func init() {
	projectConfigCmd.AddCommand(projectConfigInitCmd)
}

func emptyValidator(s string) error {
	if len(s) == 0 {
		return errors.New("this cannot be empty")
	}

	return nil
}
