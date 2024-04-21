package project

import (
	"fmt"
	"github.com/FriendsOfShopware/shopware-cli/extension"
	"github.com/FriendsOfShopware/shopware-cli/logging"
	"github.com/FriendsOfShopware/shopware-cli/shop"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var projectDebug = &cobra.Command{
	Use:   "debug",
	Short: "Shows detected Shopware version and detected extensions for further debugging",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		args[0], err = filepath.Abs(args[0])
		if err != nil {
			return err
		}

		shopCfg, err := shop.ReadConfig(projectConfigPath, true)
		if err != nil {
			return err
		}

		shopwareConstraint, err := extension.GetShopwareProjectConstraint(args[0])
		if err != nil {
			return err
		}

		if shopCfg.IsFallback() {
			fmt.Printf("Could not find a %s, using fallback config\n", projectConfigPath)
		} else {
			fmt.Printf("Found config: Yes\n")
		}
		fmt.Printf("Detected following Shopware version: %s\n", shopwareConstraint.String())

		sources := extension.FindAssetSourcesOfProject(logging.DisableLogger(cmd.Context()), args[0], shopCfg)

		fmt.Println("Following extensions/bundles has been detected")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetColWidth(100)
		table.SetHeader([]string{"Name", "Path"})

		for _, source := range sources {
			table.Append([]string{source.Name, source.Path})
		}

		table.Render()

		return nil
	},
}

func init() {
	projectRootCmd.AddCommand(projectDebug)
}
