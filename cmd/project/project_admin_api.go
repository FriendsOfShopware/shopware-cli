package project

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"shopware-cli/shop"
)

var projectAdminApiCmd = &cobra.Command{
	Use:   "admin-api [method] [path]",
	Short: "pre authenticated curl interface to the Admin API",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		var cfg *shop.Config
		var err error

		if cfg, err = shop.ReadConfig(projectConfigPath); err != nil {
			return err
		}

		if cfg.AdminApi == nil {
			return fmt.Errorf("admin api is not activated in the config")
		}

		source, err := cfg.AdminApi.GetTokenSource(cobraCmd.Context(), cfg.URL)

		if err != nil {
			return err
		}

		token, err := source.Token()

		if err != nil {
			return err
		}

		tokenOnly, _ := cobraCmd.PersistentFlags().GetBool("output-token")

		if tokenOnly {
			fmt.Println(token.AccessToken)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("command needs 2 arguments")
		}

		options := []string{"-X", args[0], fmt.Sprintf("%s/api%s", cfg.URL, args[1]), "--header", fmt.Sprintf("Authorization: %s", token.AccessToken)}

		if len(args) > 2 {
			options = append(options, args[2:]...)
		}

		cmd := exec.Command("curl", options...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	},
}

func init() {
	projectAdminApiCmd.PersistentFlags().Bool("output-token", false, "Output only token")
	projectRootCmd.AddCommand(projectAdminApiCmd)
}
