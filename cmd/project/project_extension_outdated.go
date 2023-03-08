package project

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/FriendsOfShopware/shopware-cli/logging"
	"github.com/FriendsOfShopware/shopware-cli/shop"
	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"

	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/FriendsOfShopware/shopware-cli/shop"
)

var projectExtensionOutdatedCmd = &cobra.Command{
	Use:   "outdated",
	Short: "List all outdated extensions",
	RunE: func(cmd *cobra.Command, _ []string) error {
		var cfg *shop.Config
		var err error

		outputAsJson, _ := cmd.PersistentFlags().GetBool("json")

		if cfg, err = shop.ReadConfig(projectConfigPath); err != nil {
			return err
		}

		client, err := shop.NewShopClient(cmd.Context(), cfg)
		if err != nil {
			return err
		}

		if _, err := client.ExtensionManager.Refresh(adminSdk.NewApiContext(cmd.Context())); err != nil {
			return err
		}

		extensions, _, err := client.ExtensionManager.ListAvailableExtensions(adminSdk.NewApiContext(cmd.Context()))
		extensions = extensions.FilterByUpdateable()

		if err != nil {
			return err
		}

		if outputAsJson {
			content, err := json.Marshal(extensions)

			if err != nil {
				return err
			}

			fmt.Println(string(content))

			return nil
		}

		if len(extensions) == 0 {
			logging.FromContext(cmd.Context()).Infof("All extensions are up-to-date")
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetColWidth(100)
		table.SetHeader([]string{"Name", "Current Version", "Latest Version", "Update Source"})

		for _, extension := range extensions {
			table.Append([]string{extension.Name, extension.Version, extension.LatestVersion, extension.UpdateSource})
		}

		table.Render()

		os.Exit(1)

		return nil
	},
}

func init() {
	projectExtensionCmd.AddCommand(projectExtensionOutdatedCmd)
	projectExtensionOutdatedCmd.PersistentFlags().Bool("json", false, "Output as json")
}
