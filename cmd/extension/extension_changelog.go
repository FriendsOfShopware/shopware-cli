package extension

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/FriendsOfShopware/shopware-cli/extension"
)

var extensionChangelogCmd = &cobra.Command{
	Use:   "get-changelog [path]",
	Short: "Get the changelog",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("cannot find path: %w", err)
		}

		stat, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("cannot find path: %w", err)
		}

		var ext extension.Extension

		if stat.IsDir() {
			ext, err = extension.GetExtensionByFolder(path)
		} else {
			ext, err = extension.GetExtensionByZip(path)
		}

		if err != nil {
			return fmt.Errorf("changelog: cannot open extension %w", err)
		}

		changelog, err := ext.GetChangelog()
		if err != nil {
			return fmt.Errorf("cannot generate changelog: %w", err)
		}

		isGermanChangelog, _ := cmd.PersistentFlags().GetBool("german")

		if isGermanChangelog {
			fmt.Println(changelog.German)
		} else {
			fmt.Println(changelog.English)
		}

		return nil
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionChangelogCmd)
	extensionChangelogCmd.PersistentFlags().Bool("german", false, "Get the german changelog")
}
